package cache

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	redisDao "github.com/liumingmin/goutils/db/redis"
	"github.com/liumingmin/goutils/log"
)

func TestRdscCacheFunc(t *testing.T) {
	ctx := context.Background()
	cacher := mockGetCacher()

	err := DeleteCache(ctx, cacher, "UTKey")
	if err != nil {
		t.Error(err)
	}

	value1, err := CacheFunc0(ctx, cacher, 60*time.Second, rawGetFunc1, "UTKey")
	if err != nil {
		t.Error(err)
	}

	value2, err := CacheFunc0(ctx, cacher, 60*time.Second, rawGetFunc1, "UTKey")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

	err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v", "p1"))
	if err != nil {
		t.Error(err)
	}

	value1, err = CacheFunc1(ctx, cacher, 60*time.Second, rawGetFunc2, fmt.Sprintf("UT:%v", "p1"), "p1")
	if err != nil {
		t.Error(err)
	}

	value2, err = CacheFunc1(ctx, cacher, 60*time.Second, rawGetFunc2, fmt.Sprintf("UT:%v", "p1"), "p1")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

	err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v", "p1", "p2"))
	if err != nil {
		t.Error(err)
	}

	value1, err = CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc0, fmt.Sprintf("UT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	value2, err = CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc0, fmt.Sprintf("UT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

	param3 := &testCacheParam{Param1: "p3"}
	err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1))
	if err != nil {
		t.Error(err)
	}

	value1, err = CacheFunc3(ctx, cacher, 60*time.Second, rawGetFunc3, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1), "p1", "p2", param3)
	if err != nil {
		t.Error(err)
	}

	value2, err = CacheFunc3(ctx, cacher, 60*time.Second, rawGetFunc3, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1), "p1", "p2", param3)
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}
}

func TestRdsDeleteCacheTestMore(t *testing.T) {
	ctx := context.Background()
	cacher := mockGetCacher()

	var err error

	err = DeleteCache(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"))
	if err != nil {
		t.Error(err)
	}

	result1, err := CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc4, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result1, err, mockPrintKind(result1))

	result2, err := CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc4, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result2, err, mockPrintKind(result2))

	if !reflect.DeepEqual(result1, result2) {
		t.Error(result1, result2)
	}

	err = DeleteCache(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"))
	if err != nil {
		t.Error(err)
	}

	result3, err := CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc5, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result3, err, mockPrintKind(result3))

	result4, err := CacheFunc2(ctx, cacher, 60*time.Second, rawGetFunc5, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result4, err, mockPrintKind(result4))

	if !reflect.DeepEqual(result3, result4) {
		t.Error(result3, result4)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func mockGetCacher() ICacher {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", time.Second*2)
	if err != nil {
		fmt.Println("Please install redis on local and start at port: 6379, then run test.")
		return &mockMemCacher{m: make(map[string]string)}
	}
	conn.Close()

	redisDao.InitRedises()
	const RDSC_DB = "rdscdb"
	rds := &RdsCacher{redisDao.Get(RDSC_DB)}
	return rds
}

func mockPrintKind(result interface{}) reflect.Kind {
	if result == nil {
		return reflect.Invalid
	}
	return reflect.TypeOf(result).Kind()
}

func mockErr() error {
	return nil
}

type testCacheParam struct {
	Param1 string
}

type mockMemCacher struct {
	m map[string]string
}

func (r *mockMemCacher) Del(ctx context.Context, key string) error {
	delete(r.m, key)
	return nil
}

func (r *mockMemCacher) Get(ctx context.Context, key string) (string, error) {
	if str, ok := r.m[key]; ok {
		return str, nil
	}

	return "", errors.New("not found")
}

func (r *mockMemCacher) Set(ctx context.Context, key string, str string, expire time.Duration) error {
	r.m[key] = str
	return nil
}

type mockCacheDataStruct struct {
	PersonId   string  `json:"personId" binding:"required"`
	Subject    string  `json:"subject" `
	NotifyType int     `json:"notifyType" `
	Amount     float64 `json:"amount" `
}

func rawGetFunc0(ctx context.Context, p1, p2 string) (string, error) {
	return fmt.Sprintf("TEST:%v:%v", p1, p2), mockErr()
}

func rawGetFunc1(ctx context.Context) (string, error) {
	return "TEST", mockErr()
}

func rawGetFunc2(ctx context.Context, p1 string) (string, error) {
	return fmt.Sprintf("TEST:%v", p1), mockErr()
}

func rawGetFunc3(ctx context.Context, p1, p2 string, p3 *testCacheParam) (string, error) {
	return fmt.Sprintf("TEST:%v:%v:%v", p1, p2, p3.Param1), mockErr()
}

func rawGetFunc4(ctx context.Context, p1, p2 string) (mockCacheDataStruct, error) {
	return mockCacheDataStruct{
		PersonId:   p1,
		Subject:    p2,
		NotifyType: 2,
		Amount:     19.55,
	}, mockErr()
}

func rawGetFunc5(ctx context.Context, p1, p2 string) (*mockCacheDataStruct, error) {
	result, _ := rawGetFunc4(ctx, p1, p2)
	return &result, mockErr()
}
