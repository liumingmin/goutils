package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"

	"github.com/demdxx/gocast"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/proto"
)

var (
	sg       = singleflight.Group{}
	ctxIface reflect.Type
	pbIface  reflect.Type
)

func RdsDeleteCache(ctx context.Context, rds redis.UniversalClient, keyFmt string, args ...interface{}) (err error) {
	key := fmt.Sprintf(keyFmt, args...)

	log.Debug(ctx, "RdsDeleteCache cache key : %v", key)

	return rds.Del(ctx, key).Err()
}

func RdsCacheFunc(ctx context.Context, rds redis.UniversalClient, rdsExpire int, f interface{}, keyFmt string,
	args ...interface{}) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("RdsCacheFunc err: %v", e)
	})

	ft := reflect.TypeOf(f)
	if ft.NumOut() == 0 {
		log.Error(ctx, "RdsCacheFunc f must have one return value")
		return nil, errors.New("f must have one return value")
	}

	key := fmt.Sprintf(keyFmt, args...)
	log.Debug(ctx, "RdsCacheFunc cache key : %v", key)

	retValue, err := rds.Get(ctx, key).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheFunc hit cache : %v", retValue)

		return convertStringTo(ctx, retValue, ft.Out(0)), nil
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		return rdsCacheCallFunc(ctx, rds, rdsExpire, f, keyFmt, args...)
	})
	return data, err
}

//通用场景过于复杂，限定参数,返回类型为 map[string]*struct
//fMulti示例: func getThingsByIds(ctx context.Context, ids []string) (map[string]*Thing, error)
func RdsCacheMultiFunc(ctx context.Context, rds redis.UniversalClient, rdsExpire int, fMulti interface{}, keyFmt string,
	args []string) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("RdsCacheMultiFunc err: %v", e)
	})

	//get and check return value
	ft := reflect.TypeOf(fMulti)
	retMapType := ft.Out(0)
	if retMapType.Kind() != reflect.Map {
		log.Error(ctx, "RdsCacheMultiFunc f must have return map value")
		return nil, errors.New("f must have return map value")
	}

	//map item type, may be ptr or struct...
	retItemType := retMapType.Elem()

	//prepare cache key
	keys := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		keys[i] = fmt.Sprintf(keyFmt, args[i])
	}

	noCachedArgs := make([]string, 0, len(args))
	resultValues := reflect.MakeMapWithSize(retMapType, len(args))

	retValues, err := rds.MGet(ctx, keys...).Result()
	if err == nil {
		for i, retValue := range retValues {
			if utils.IsNil(retValue) {
				noCachedArgs = append(noCachedArgs, args[i])
			} else {
				retValueStr, _ := retValue.(string)
				resultValues.SetMapIndex(reflect.ValueOf(args[i]), reflect.ValueOf(convertStringTo(ctx, retValueStr, retItemType)))
			}
		}
	}

	if len(noCachedArgs) == 0 {
		log.Debug(ctx, "RdsCacheMultiFunc all hit caches: %v", len(retValues))
		return resultValues.Interface(), nil
	}

	//防击穿
	sort.Strings(noCachedArgs)
	sgKey := keyFmt + utils.MD5(strings.Join(noCachedArgs, ","))
	callRetValue, err, _ := sg.Do(sgKey, func() (interface{}, error) {
		return rdsCacheMultiCallFunc(ctx, rds, rdsExpire, fMulti, keyFmt, noCachedArgs, retItemType)
	})
	if err == nil && !utils.IsNil(callRetValue) {
		retMapValue := reflect.ValueOf(callRetValue)

		mapIter := retMapValue.MapRange()
		for mapIter.Next() {
			key := mapIter.Key()
			value := mapIter.Value()

			resultValues.SetMapIndex(key, value)
		}
	}

	return resultValues.Interface(), nil
}

func rdsCacheCallFunc(ctx context.Context, rds redis.UniversalClient, rdsExpire int, f interface{}, keyFmt string,
	args ...interface{}) (interface{}, error) {
	argValues := make([]reflect.Value, 0)

	ft := reflect.TypeOf(f)
	if ft.NumIn() > 0 && ft.In(0).Implements(ctxIface) {
		argValues = append(argValues, reflect.ValueOf(ctx))
	}

	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}

	fv := reflect.ValueOf(f)
	retValues := fv.Call(argValues)

	var retErr error
	if len(retValues) > 1 && !utils.SafeIsNil(&retValues[1]) {
		retErr, _ = retValues[1].Interface().(error)
	}

	key := fmt.Sprintf(keyFmt, args...)

	var result interface{}
	if len(retValues) > 0 && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		result = retValues[0].Interface()
		rds.Set(ctx, key, convertRetValueToString(ctx, result, ft.Out(0)), time.Duration(rdsExpire)*time.Second)
	} else {
		rds.Set(ctx, key, "", time.Duration(utils.Min(rdsExpire, 10))*time.Second) //防止缓存穿透
		log.Debug(ctx, "RdsCacheFunc avoid cache through: %v", key)
	}
	return result, retErr
}

func rdsCacheMultiCallFunc(ctx context.Context, rds redis.UniversalClient, rdsExpire int, fMulti interface{}, keyFmt string,
	args []string, retItemType reflect.Type) (interface{}, error) {
	argValues := make([]reflect.Value, 0, 2)
	argValues = append(argValues, reflect.ValueOf(ctx))
	argValues = append(argValues, reflect.ValueOf(args))

	fv := reflect.ValueOf(fMulti)
	retValues := fv.Call(argValues)

	var retErr error
	if len(retValues) > 1 && !utils.SafeIsNil(&retValues[1]) {
		retErr, _ = retValues[1].Interface().(error)
	}

	var result interface{}
	if len(retValues) > 0 && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		mapIter := retValues[0].MapRange()
		i := 0
		for mapIter.Next() {
			key := mapIter.Key()
			value := mapIter.Value()

			rds.Set(ctx, key.String(), convertRetValueToString(ctx, value.Interface(), retItemType),
				time.Duration(rdsExpire+(i%10))*time.Second) //防缓存雪崩
			i++
		}

		result = retValues[0].Interface()
	} else {
		//prepare cache key
		for i := 0; i < len(args); i++ {
			key := fmt.Sprintf(keyFmt, args[i])
			rds.Set(ctx, key, "", time.Duration(utils.Min(rdsExpire, 10))*time.Second) //防止缓存穿透
		}
		log.Debug(ctx, "rdsCacheMultiCallFunc avoid cache through: %v", len(args))
	}
	return result, retErr
}

func convertRetValueToString(ctx context.Context, retValue interface{}, t reflect.Type) string {
	if t.Implements(pbIface) {
		if pbMsg, ok := retValue.(proto.Message); ok {
			data, err := proto.Marshal(pbMsg)
			if err != nil {
				log.Error(ctx, "proto.Marshal err: %v", err)
			}
			return string(data)
		}
	}

	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Struct, reflect.Map:
		bs, err := json.Marshal(retValue)
		if err != nil {
			log.Error(ctx, "json.Marshal err: %v", err)
		}
		return string(bs)
	default:
		return gocast.ToString(retValue)
	}
	return gocast.ToString(retValue)
}

func convertStringTo(ctx context.Context, cacheValue string, t reflect.Type) interface{} {
	if strings.TrimSpace(cacheValue) == "" {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
			return nil
		}
		return reflect.Zero(t)
	}

	if t.Implements(pbIface) {
		tt := t.Elem()
		if retValue, ok := reflect.New(tt).Interface().(proto.Message); ok {
			err := proto.Unmarshal([]byte(cacheValue), retValue)
			if err != nil {
				log.Error(ctx, "proto.Unmarshal err: %v", err)
			}
			return retValue
		}
		return nil
	}

	switch t.Kind() {
	case reflect.String:
		return cacheValue
	case reflect.Ptr:
		tt := t.Elem()
		retValue := reflect.New(tt).Interface()
		err := json.Unmarshal([]byte(cacheValue), retValue)
		if err != nil {
			log.Error(ctx, "json.Unmarshal err: %v", err)
		}
		return retValue
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Interface, reflect.Struct:
		retValue := reflect.New(t)
		retValueInterface := retValue.Interface()
		err := json.Unmarshal([]byte(cacheValue), retValueInterface)
		if err != nil {
			log.Error(ctx, "json.Unmarshal err: %v", err)
		}
		return retValue.Elem().Interface()
	default:
		result, err := gocast.ToT(cacheValue, t, "")
		if err != nil {
			log.Error(ctx, "gocast.ToT err: %v", err)
		}
		return result
	}
	return cacheValue
}

func init() {
	var ctx context.Context
	ctxIface = reflect.TypeOf(&ctx).Elem()

	var pbMsg proto.Message
	pbIface = reflect.TypeOf(&pbMsg).Elem()
}
