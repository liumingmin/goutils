package cache_func

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/liumingmin/goutils/log4go"
	"github.com/robfig/go-cache"
)

var (
	gDefMemCache *cache.Cache
)

const (
	MEM_CACHE_FUNC_KEY      = "MCF:%s"
	MEM_CACHE_CONF_FUNC_KEY = "def_cache_func"
)

func DefMemCacheFunc(expire time.Duration, f interface{}, prefix string, args ...interface{}) (interface{}, error) {
	return MemCacheFunc(gDefMemCache, expire, f, prefix, args...)
}

func DefMemCacheDelete(prefix string, args ...interface{}) bool {
	return MemCacheDelete(gDefMemCache, prefix, args...)
}

func MemCacheFunc(cc *cache.Cache, expire time.Duration, f interface{}, prefix string, args ...interface{}) (interface{}, error) {
	//defer logutil.Recover(context.Background(), func(e interface{}) string {
	//	err := e.(error)
	//	return fmt.Sprintf("CacheFunc err: %v", err)
	//})

	ft := reflect.TypeOf(f)
	if ft.NumOut() == 0 {
		log4go.Error("CacheFunc f must have one return value")
		return nil, nil
	}

	key := prefix + ":" + getMemCacheKey(args...)
	log4go.Debug("MemCacheFunc cache key : %v", key)

	retValue, ok := cc.Get(key)
	if ok {
		log4go.Debug("MemCacheFunc hit cache : %v", retValue)

		return retValue, nil
	} else {
		fv := reflect.ValueOf(f)

		var argValues []reflect.Value
		for _, arg := range args {
			argValues = append(argValues, reflect.ValueOf(arg))
		}

		retValues := fv.Call(argValues)
		var result interface{}
		var err error
		if len(retValues) > 0 && retValues[0].IsValid() {
			result = retValues[0].Interface()
			cc.Set(key, result, expire)
		} else {
			cc.Set(key, nil, expire) //防止缓存穿透
		}

		if len(retValues) > 1 && retValues[1].IsValid() {
			if !retValues[1].IsNil() {
				err = retValues[1].Interface().(error)
			}
		}
		return result, err
	}
}

func MemCacheDelete(cc *cache.Cache, prefix string, args ...interface{}) bool {
	key := prefix + ":" + getMemCacheKey(args...)
	return cc.Delete(key)
}

func DefMemCacheGetValue(prefix string, args ...interface{}) (interface{}, bool) {
	return MemCacheGetValue(gDefMemCache, prefix, args...)
}

func MemCacheGetValue(cc *cache.Cache, prefix string, args ...interface{}) (interface{}, bool) {
	key := prefix + ":" + getMemCacheKey(args...)
	log4go.Debug("GetMemValue cache key : %v", key)

	return cc.Get(key)
}

func getMemCacheKey(args ...interface{}) string {
	if len(args) == 0 {
		log4go.Error("getMemCacheKey args len is 0")
		return MEM_CACHE_FUNC_KEY
	}

	if len(args) == 1 {
		t := reflect.TypeOf(args[0])
		if t.Kind() == reflect.String || t.Kind() == reflect.Int ||
			t.Kind() == reflect.Int32 || t.Kind() == reflect.Int64 {
			return fmt.Sprintf(MEM_CACHE_FUNC_KEY, fmt.Sprint(args[0]))
		}
	}

	return fmt.Sprintf(MEM_CACHE_FUNC_KEY, MD5Args(args...))
}

func MD5(origStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(origStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//支持原生值类型，值结构体，原生类型数组
func MD5Args(args ...interface{}) string {
	var b bytes.Buffer
	for _, arg := range args {
		b.WriteString("^")

		switch arg.(type) {
		case []string:
			s := arg.([]string)
			sort.Strings(s)
			b.WriteString(fmt.Sprintf("%v", s))
		case []int:
			ii := arg.([]int)
			sort.Ints(ii)
			b.WriteString(fmt.Sprintf("%v", ii))
		case []int64:
			ii := arg.([]int64)
			var is []int
			for _, i := range ii {
				is = append(is, int(i))
			}
			sort.Ints(is)
			b.WriteString(fmt.Sprintf("%v", is))
		case []float64:
			f := arg.([]float64)
			sort.Float64s(f)
			b.WriteString(fmt.Sprintf("%v", f))
		default:
			b.WriteString(fmt.Sprintf("%#v", arg))
		}
	}

	return MD5(b.String())
}

func init() {
	gDefMemCache = Get(MEM_CACHE_CONF_FUNC_KEY)
}
