# cache 缓存模块
## mem_cache_test.go 内存缓存
### TestMemCacheFunc
```go

ctx := context.Background()

const cacheKey = "UT:%v:%v"

var lCache = cache.New(5*time.Minute, 5*time.Minute)
result, err := MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
log.Info(ctx, "%v %v %v", result, err, printKind(result))

_memCacheFuncTestMore(ctx, lCache, cacheKey)
```
## rds_cache_test.go Redis缓存
### TestRdscCacheFunc
```go

redisDao.InitRedises()
ctx := context.Background()

const cacheKey = "UT:%v:%v"
const RDSC_DB = "rdscdb"

rds := redisDao.Get(RDSC_DB)

result, err := RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
log.Info(ctx, "%v %v %v", result, err, printKind(result))

_rdsDeleteCacheTestMore(ctx, rds, cacheKey)
```
### TestRdsCacheMultiFunc
```go

redisDao.InitRedises()
ctx := context.Background()
const RDSC_DB = "rdscdb"

rds := redisDao.Get(RDSC_DB)
result, err := RdsCacheMultiFunc(ctx, rds, 30, getThingsByIds, "multikey:%s", []string{"1", "2", "5", "3", "4", "10"})
if err == nil && result != nil {
	mapValue, ok := result.(map[string]*Thing)
	if ok {
		for key, value := range mapValue {
			log.Info(ctx, "%v===%v", key, value)
		}
	}
}
```
