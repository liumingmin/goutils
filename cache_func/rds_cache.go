package cache_func

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/demdxx/gocast"
	"github.com/liumingmin/goutils/goredis"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"golang.org/x/sync/singleflight"
)

var (
	sg = singleflight.Group{}
)

func RdsDeleteCache(ctx context.Context, dbName string, keyFmt string, args ...interface{}) (err error) {
	rds := goredis.Get(dbName)

	key := fmt.Sprintf(keyFmt, args...)

	log.Debug(ctx, "RdsDeleteCache cache key : %v", key)

	return rds.Del(ctx, key).Err()
}

func RdsCacheFunc(ctx context.Context, dbName string, rdsExpire int, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err := e.(error)
		return fmt.Sprintf("RdscCacheFuncCtx err: %v", err)
	})

	rds := goredis.Get(dbName)

	ft := reflect.TypeOf(f)
	if ft.NumOut() == 0 {
		log.Error(ctx, "RdsCacheFunc f must have one return value")
		return nil, nil
	}

	key := fmt.Sprintf(keyFmt, args...)
	log.Debug(ctx, "RdsCacheFunc cache key : %v", key)

	retValue, err := rds.Get(ctx, key).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheFunc hit cache : %v", retValue)

		return convertStringTo(retValue, ft.Out(0)), nil
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		return rdsCacheCallFunc(ctx, dbName, rdsExpire, f, keyFmt, args...)
	})
	return data, err
}

func rdsCacheCallFunc(ctx context.Context, dbName string, rdsExpire int, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
	argValues := make([]reflect.Value, 0)

	ft := reflect.TypeOf(f)

	var iface context.Context
	ctxIface := reflect.TypeOf(&iface).Elem()
	if ft.NumIn() > 0 && ft.In(0).Implements(ctxIface) {
		argValues = append(argValues, reflect.ValueOf(ctx))
	}

	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}

	fv := reflect.ValueOf(f)
	retValues := fv.Call(argValues)

	var retErr error
	if len(retValues) > 1 && retValues[1].IsValid() && !utils.SafeIsNil(&retValues[1]) {
		retErr, _ = retValues[1].Interface().(error)
	}

	rds := goredis.Get(dbName)
	key := fmt.Sprintf(keyFmt, args...)

	var result interface{}
	if len(retValues) > 0 && retValues[0].IsValid() && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		result = retValues[0].Interface()
		rds.Set(ctx, key, convertRetValueToString(result, ft.Out(0)), time.Duration(rdsExpire)*time.Second)
	} else {
		rds.Set(ctx, key, "", time.Duration(utils.Min(rdsExpire, 20))*time.Second) //防止缓存穿透
		log.Debug(ctx, "RdsCacheFunc avoid cache through: %v", key)
	}
	return result, retErr
}

func convertRetValueToString(retValue interface{}, t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Struct, reflect.Map:
		bs, _ := json.Marshal(retValue)
		return string(bs)
	default:
		return gocast.ToString(retValue)
	}
	return gocast.ToString(retValue)
}

func convertStringTo(cacheValue string, t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.String:
		return cacheValue
	case reflect.Ptr:
		if strings.TrimSpace(cacheValue) == "" {
			return nil
		}

		tt := t.Elem()
		retValue := reflect.New(tt).Interface()
		json.Unmarshal([]byte(cacheValue), retValue)
		return retValue
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		if strings.TrimSpace(cacheValue) == "" {
			return nil
		}

		retValue := reflect.New(t)
		retValueInterface := retValue.Interface()
		json.Unmarshal([]byte(cacheValue), retValueInterface)
		return retValue.Elem().Interface()
	case reflect.Struct:
		retValue := reflect.New(t)
		retValueInterface := retValue.Interface()
		json.Unmarshal([]byte(cacheValue), retValueInterface)
		return retValue.Elem().Interface()
	default:
		result, _ := gocast.ToT(cacheValue, t, "")
		return result
	}
	return cacheValue
}
