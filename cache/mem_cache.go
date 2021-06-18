package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/robfig/go-cache"
)

func MemCacheFunc(ctx context.Context, cc *cache.Cache, expire time.Duration, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("MemCacheFunc err: %v", e)
	})

	ft := reflect.TypeOf(f)
	if ft.NumOut() == 0 {
		log.Error(ctx, "MemCacheFunc f must have one return value")
		return nil, nil
	}

	key := fmt.Sprintf(keyFmt, args...)
	log.Debug(ctx, "MemCacheFunc cache key : %v", key)

	retValue, ok := cc.Get(key)
	if ok {
		log.Debug(ctx, "MemCacheFunc hit cache : %v", retValue)

		return retValue, nil
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		return memCacheCallFunc(ctx, cc, expire, f, keyFmt, args...)
	})
	return data, err
}

func memCacheCallFunc(ctx context.Context, cc *cache.Cache, expire time.Duration, f interface{}, keyFmt string, args ...interface{}) (interface{}, error) {
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

	key := fmt.Sprintf(keyFmt, args...)

	var result interface{}
	if len(retValues) > 0 && retValues[0].IsValid() && !utils.SafeIsNil(&retValues[0]) && retErr == nil {
		result = retValues[0].Interface()
		cc.Set(key, result, expire)
	} else {
		cc.Set(key, nil, time.Duration(utils.Min64(int64(expire), 20*int64(time.Second)))) //防止缓存穿透
		log.Debug(ctx, "MemCacheFunc avoid cache through: %v", key)
	}
	return result, retErr
}

func MemCacheDelete(ctx context.Context, cc *cache.Cache, keyFmt string, args ...interface{}) bool {
	key := fmt.Sprintf(keyFmt, args...)
	return cc.Delete(key)
}
