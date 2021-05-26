package cache

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/liumingmin/goutils/utils"

	"github.com/gin-gonic/gin"
)

func cachefunctest(a int, s []string, b map[string]string) (map[string]string, error) {
	//return map[string]string{"aaa": "bbb", "cccc": b[0]}
	time.Sleep(time.Millisecond * 100)
	return b, errors.New("3333")
}

func TestCacheFunc(t *testing.T) {
	memcache := SimpleMemCache{}

	result, err := CacheFunc(&memcache, 1, utils.ConsistArgs,
		cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})

	t.Log("1", result, err)

	result, err = CacheFunc(&memcache, 1, utils.ConsistArgs,
		cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})

	t.Log("2", result, err)

	result, err = CacheFunc(&memcache, 1, utils.ConsistArgs,
		cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})

	t.Log("2", result, err)

	result, err = CacheFunc(&memcache, 1, utils.ConsistArgs,
		cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})

	t.Log("2", result, err)
}

func BenchmarkCacheFunc0(b *testing.B) {
	//memcache := SimpleMemCache{}

	for i := 0; i < b.N; i++ {
		cachefunctest(222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})
		//b.Log("1", result, err)
	}
}

func BenchmarkCacheFunc(b *testing.B) {
	memcache := SimpleMemCache{}

	for i := 0; i < b.N; i++ {
		CacheFunc(&memcache, 1, utils.ConsistArgs,
			cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})
		//b.Log("1", result, err)
	}
}

func TestTCacheFunc(t *testing.T) {
	memcache := SimpleMemCache{}
	tcf := TCacheFunc{store: &memcache, cf: utils.ConsistArgs}
	result, err := tcf.Cache(1, cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})
	t.Log("1", result, err)

	result, err = tcf.Cache(1, cachefunctest, 222, []string{"aaa", "222"}, map[string]string{"aaa": "bbb", "111": "ddd"})
	t.Log("2", result, err)
}

func TestCachePage(t *testing.T) {
	router := gin.New()

	memcache := SimpleMemCache{}

	router.GET("/cache_ping", CachePage(&memcache, time.Second*2, SimpleRequestFunc, func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().UnixNano()))
	}))

	w1 := performRequest("GET", "/cache_ping", router)
	time.Sleep(time.Second * 3)
	w2 := performRequest("GET", "/cache_ping", router)

	fmt.Println(w1.Body.String(), w2.Body.String())

	//assert.Equal(t, 200, w1.Code)
	//assert.Equal(t, 200, w2.Code)
	//assert.Equal(t, w1.Body.String(), w2.Body.String())
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
