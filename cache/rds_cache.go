package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/demdxx/gocast"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/redis"
	"github.com/liumingmin/goutils/utils"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/proto"
)

var (
	sg       = singleflight.Group{}
	ctxIface reflect.Type
	pbIface  reflect.Type
)

func RdsDeleteCache(ctx context.Context, dbName string, keyFmt string, args ...interface{}) (err error) {
	rds := redis.Get(dbName)

	key := fmt.Sprintf(keyFmt, args...)

	log.Debug(ctx, "RdsDeleteCache cache key : %v", key)

	return rds.Del(ctx, key).Err()
}

func RdsCacheFunc(ctx context.Context, dbName string, rdsExpire int, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("RdscCacheFuncCtx err: %v", e)
	})

	rds := redis.Get(dbName)

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

		return convertStringTo(ctx, retValue, ft.Out(0)), nil
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		return rdsCacheCallFunc(ctx, dbName, rdsExpire, f, keyFmt, args...)
	})
	return data, err
}

func rdsCacheCallFunc(ctx context.Context, dbName string, rdsExpire int, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
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
	if len(retValues) > 1 && retValues[1].IsValid() && !utils.SafeIsNil(&retValues[1]) {
		retErr, _ = retValues[1].Interface().(error)
	}

	rds := redis.Get(dbName)
	key := fmt.Sprintf(keyFmt, args...)

	var result interface{}
	if len(retValues) > 0 && retValues[0].IsValid() && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		result = retValues[0].Interface()
		rds.Set(ctx, key, convertRetValueToString(ctx, result, ft.Out(0)), time.Duration(rdsExpire)*time.Second)
	} else {
		rds.Set(ctx, key, "", time.Duration(utils.Min(rdsExpire, 10))*time.Second) //防止缓存穿透
		log.Debug(ctx, "RdsCacheFunc avoid cache through: %v", key)
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
