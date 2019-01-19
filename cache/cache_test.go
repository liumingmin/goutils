package cache

import (
	"testing"

	"time"

	"github.com/liumingmin/goutils/utils"
	"github.com/pkg/errors"
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
