package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

//缓存多个参数应该封装为结构传递，如果是基本类型可以type定义后实现接口,类似StringCacheKeyer
type CacheKeyer interface {
	RCMCacheKey() string
}

type StringCacheKeyer string

func (t StringCacheKeyer) RCMCacheKey() string {
	return string(t)
}

//无法防击穿，使用场景需要注意
func RdsCacheMultiFunc(ctx context.Context, rds redis.UniversalClient, rdsExpire int, fMulti interface{}, keyFmt string,
	args []CacheKeyer) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("RdsCacheMultiFunc err: %v", e)
	})

	ft := reflect.TypeOf(fMulti)
	if ft.NumOut() == 0 {
		log.Error(ctx, "RdsCacheMultiFunc f must have one return value")
		return nil, errors.New("f must have one return value")
	}

	//get and check return value
	retSliceType := ft.Out(0)
	if retSliceType.Kind() == reflect.Ptr { //todo return value addr()
		retSliceType = retSliceType.Elem()
	}

	if retSliceType.Kind() != reflect.Slice {
		log.Error(ctx, "RdsCacheMultiFunc f must have return slice value")
		return nil, errors.New("f must have return slice value")
	}

	//slice item type, may be ptr or struct...
	retItemType := retSliceType.Elem()

	//prepare cache key
	keys := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		keys[i] = fmt.Sprintf(keyFmt, args[i].RCMCacheKey())
	}

	noCachedArgs := make([]CacheKeyer, 0, len(keys))
	resultValues := reflect.MakeSlice(retSliceType, 0, len(keys))

	retValues, err := rds.MGet(ctx, keys...).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheMultiFunc hit caches : %v", len(retValues))

		for i, retValue := range retValues {
			if utils.IsNil(retValue) {
				noCachedArgs = append(noCachedArgs, args[i])
			} else {
				retValueStr, _ := retValue.(string)
				resultValues = reflect.Append(resultValues, reflect.ValueOf(convertStringTo(ctx, retValueStr, retItemType)))
			}
		}
	}

	if len(noCachedArgs) == 0 {
		return resultValues.Interface(), nil
	}

	callRetValues, err := rdsCacheMultiCallFunc(ctx, rds, rdsExpire, fMulti, keyFmt, noCachedArgs, retItemType)
	if err == nil && callRetValues != nil {
		callRetValueSlice := reflect.ValueOf(callRetValues)
		for i := 0; i < callRetValueSlice.Len(); i++ {
			resultValues = reflect.Append(resultValues, callRetValueSlice.Index(i))
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
	args []CacheKeyer, retItemType reflect.Type) (interface{}, error) {
	argValues := make([]reflect.Value, 0)

	//prepare arg1
	ft := reflect.TypeOf(fMulti)
	if ft.NumIn() > 0 && ft.In(0).Implements(ctxIface) {
		argValues = append(argValues, reflect.ValueOf(ctx))
	}

	//check in value
	var argType reflect.Type
	if ft.NumIn() > 0 && ft.In(0).Implements(ctxIface) {
		if ft.NumIn() > 1 {
			argType = ft.In(1)
		}
	} else {
		if ft.NumIn() > 0 {
			argType = ft.In(0)
		}
	}
	if argType.Kind() != reflect.Slice {
		return nil, errors.New("arg type is not slice")
	}

	//convert to real type? todo?
	argItemType := argType.Elem()
	argSliceValue := reflect.MakeSlice(argType, len(args), len(args))
	for i := 0; i < len(args); i++ {
		argSliceValue = reflect.Append(argSliceValue, reflect.ValueOf(args[i]).Convert(argItemType))
	}

	//prepare arg2
	argValues = append(argValues, argSliceValue)

	fv := reflect.ValueOf(fMulti)
	retValues := fv.Call(argValues)

	var retErr error
	if len(retValues) > 1 && !utils.SafeIsNil(&retValues[1]) {
		retErr, _ = retValues[1].Interface().(error)
	}

	var result interface{}
	if len(retValues) > 0 && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		retValueSlice := retValues[0]
		for i := 0; i < retValueSlice.Len(); i++ {
			retValueItem := retValueSlice.Index(i).Interface()
			cacheKeyer, ok := retValueItem.(CacheKeyer)
			if !ok {
				continue
			}

			rds.Set(ctx, cacheKeyer.RCMCacheKey(), convertRetValueToString(ctx, retValueItem, retItemType),
				time.Duration(rdsExpire+(i%10))*time.Second)
		}
		result = retValueSlice.Interface()
	} else {
		for i := 0; i < len(args); i++ {
			key := fmt.Sprintf(keyFmt, args[i].RCMCacheKey())
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
