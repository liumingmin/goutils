package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/conv"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

var (
	sg                  = singleflight.Group{}
	DurPrevCacheThrough = 10 * time.Second //防止缓存穿透
)

func RdsDeleteCache(ctx context.Context, rds redis.UniversalClient, keyFmt string, args ...interface{}) (err error) {
	key := fmt.Sprintf(keyFmt, args...)

	log.Debug(ctx, "RdsDeleteCache cache key : %v", key)

	return rds.Del(ctx, key).Err()
}

func RdsCacheFunc0[Tr any](ctx context.Context, rds redis.UniversalClient, expire time.Duration, fn func(context.Context) (Tr, error),
	key string) (tr Tr, err error) {
	retValue, err := rds.Get(ctx, key).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheFunc0 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				rds.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		rds.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}

func RdsCacheFunc1[T1 any, Tr any](ctx context.Context, rds redis.UniversalClient, expire time.Duration, fn func(context.Context, T1) (Tr, error),
	keyFmt string, t1 T1) (tr Tr, err error) {
	key := fmt.Sprintf(keyFmt, t1)

	retValue, err := rds.Get(ctx, key).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheFunc1 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx, t1)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				rds.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		rds.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}

func RdsCacheFunc2[T1 any, T2 any, Tr any](ctx context.Context, rds redis.UniversalClient, expire time.Duration, fn func(context.Context, T1, T2) (Tr, error),
	keyFmt string, t1 T1, t2 T2) (tr Tr, err error) {
	key := fmt.Sprintf(keyFmt, t1, t2)

	retValue, err := rds.Get(ctx, key).Result()
	if err == nil {
		log.Debug(ctx, "RdsCacheFunc2 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx, t1, t2)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				rds.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		rds.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}
