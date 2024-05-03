package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/conv"

	"golang.org/x/sync/singleflight"
)

var (
	sg                  = singleflight.Group{}
	DurPrevCacheThrough = 10 * time.Second //防止缓存穿透
)

type ICacher interface {
	Del(context.Context, string) error
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error
}

type RdsCacher struct {
	Rds redis.UniversalClient
}

func (r *RdsCacher) Del(ctx context.Context, key string) error {
	return r.Rds.Del(ctx, key).Err()
}

func (r *RdsCacher) Get(ctx context.Context, key string) (string, error) {
	return r.Rds.Get(ctx, key).Result()
}

func (r *RdsCacher) Set(ctx context.Context, key string, str string, expire time.Duration) error {
	return r.Rds.Set(ctx, key, str, expire).Err()
}

func DeleteCache(ctx context.Context, cacher ICacher, key string) (err error) {
	return cacher.Del(ctx, key) //rds
}

func CacheFunc0[Tr any](ctx context.Context, cacher ICacher, expire time.Duration, fn func(context.Context) (Tr, error),
	key string) (tr Tr, err error) {
	retValue, err := cacher.Get(ctx, key)
	if err == nil {
		log.Debug(ctx, "CacheFunc0 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				cacher.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		cacher.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}

func CacheFunc1[T1, Tr any](ctx context.Context, cacher ICacher, expire time.Duration, fn func(context.Context, T1) (Tr, error),
	key string, t1 T1) (tr Tr, err error) {

	retValue, err := cacher.Get(ctx, key)
	if err == nil {
		log.Debug(ctx, "CacheFunc1 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx, t1)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				cacher.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		cacher.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}

func CacheFunc2[T1, T2, Tr any](ctx context.Context, cacher ICacher, expire time.Duration, fn func(context.Context, T1, T2) (Tr, error),
	key string, t1 T1, t2 T2) (tr Tr, err error) {

	retValue, err := cacher.Get(ctx, key)
	if err == nil {
		log.Debug(ctx, "CacheFunc2 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx, t1, t2)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				cacher.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		cacher.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}

func CacheFunc3[T1, T2, T3, Tr any](ctx context.Context, cacher ICacher, expire time.Duration, fn func(context.Context, T1, T2, T3) (Tr, error),
	key string, t1 T1, t2 T2, t3 T3) (tr Tr, err error) {

	retValue, err := cacher.Get(ctx, key)
	if err == nil {
		log.Debug(ctx, "CacheFunc3 hit cache: key: %v, value: %v", key, retValue)

		return conv.StringToValue[Tr](retValue)
	}

	data, err, _ := sg.Do(key, func() (interface{}, error) {
		tr, err := fn(ctx, t1, t2, t3)
		if err == nil {
			str, err := conv.ValueToString(tr)
			if err == nil {
				cacher.Set(ctx, key, str, expire)
				return tr, nil
			}
		}

		cacher.Set(ctx, key, "", utils.Min(expire, DurPrevCacheThrough))
		return tr, err
	})
	return data.(Tr), err
}
