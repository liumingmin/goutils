package cache

import (
	"context"
	"testing"
	"time"

	"github.com/robfig/go-cache"

	"github.com/liumingmin/goutils/log"
)

func TestMemCacheFunc(t *testing.T) {
	ctx := context.Background()

	const cacheKey = "UT:%v:%v"

	var lCache = cache.New(5*time.Minute, 5*time.Minute)
	result, err := MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")
}
