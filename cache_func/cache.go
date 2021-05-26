package cache_func

import (
	"fmt"
	"os"
	"time"

	"github.com/liumingmin/goutils/conf"

	"github.com/robfig/go-cache"
)

var pools = cachePools{make(map[string]*cache.Cache)}

type CachePools interface {
	Get(key string) *cache.Cache
}

type cachePools struct {
	pools map[string]*cache.Cache
}

func (p *cachePools) Get(key string) *cache.Cache {
	return p.pools[key]
}

func (p *cachePools) add(key string, cache *cache.Cache) {
	p.pools[key] = cache
}

func Get(key string) *cache.Cache {
	return pools.Get(key)
}

func init() {
	cachesInfo := conf.Ext("caches", []map[string]string{})
	if cachesInfo == nil {
		fmt.Fprintf(os.Stderr, "no caches key in conf file...")
		return
	}
	confs, _ := cachesInfo.([]interface{})
	for _, v := range confs {
		vv, _ := v.(map[interface{}]interface{})
		cacheKey := getStr(vv, "key")
		cacheTime := getStr(vv, "time")
		cacheGc := getStr(vv, "gc")
		pools.add(cacheKey, newCache(cacheTime, cacheGc))
	}
}

func newCache(timeStr string, gcStr string) *cache.Cache {
	var cacheConf *cache.Cache
	var timeD, gcD time.Duration
	if timeStr == "" {
		timeStr = "5m"
	}
	if gcStr == "" {
		gcStr = "30s"
	}
	timeD, _ = time.ParseDuration(timeStr)
	gcD, _ = time.ParseDuration(gcStr)
	cacheConf = cache.New(timeD, gcD)
	return cacheConf
}

func getStr(m map[interface{}]interface{}, key string) string {
	v := m[key]
	if v != nil {
		return fmt.Sprint(v)
	} else {
		return ""
	}
}
