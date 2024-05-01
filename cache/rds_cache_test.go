package cache

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	redisDao "github.com/liumingmin/goutils/db/redis"
	"github.com/liumingmin/goutils/log"
)

func TestRdscCacheFunc(t *testing.T) {
	if !isRdsRun() {
		return
	}

	redisDao.InitRedises()
	ctx := context.Background()

	const RDSC_DB = "rdscdb"

	rds := redisDao.Get(RDSC_DB)

	err := RdsDeleteCache(ctx, rds, "UTKey")
	if err != nil {
		t.Error(err)
	}

	value1, err := RdsCacheFunc0(ctx, rds, 60*time.Second, rawGetFunc1, "UTKey")
	if err != nil {
		t.Error(err)
	}

	value2, err := RdsCacheFunc0(ctx, rds, 60*time.Second, rawGetFunc1, "UTKey")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

	err = RdsDeleteCache(ctx, rds, "UT:%v", "p1")
	if err != nil {
		t.Error(err)
	}

	value1, err = RdsCacheFunc1(ctx, rds, 60*time.Second, rawGetFunc2, "UT:%v", "p1")
	if err != nil {
		t.Error(err)
	}

	value2, err = RdsCacheFunc1(ctx, rds, 60*time.Second, rawGetFunc2, "UT:%v", "p1")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

	err = RdsDeleteCache(ctx, rds, "UT:%v:%v", "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	value1, err = RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc0, "UT:%v:%v", "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	value2, err = RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc0, "UT:%v:%v", "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	if value1 != value2 {
		t.Error(value1, value2)
	}

}

func TestRdsDeleteCacheTestMore(t *testing.T) {
	if !isRdsRun() {
		return
	}

	redisDao.InitRedises()
	ctx := context.Background()

	const cacheKey = "UT:%v:%v"
	const RDSC_DB = "rdscdb"

	rds := redisDao.Get(RDSC_DB)

	var err error

	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result1, err := RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result1, err, printKind(result1))

	result2, err := RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result2, err, printKind(result2))

	if !reflect.DeepEqual(result1, result2) {
		t.Error(result1, result2)
	}

	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result3, err := RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result3, err, printKind(result3))

	result4, err := RdsCacheFunc2(ctx, rds, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result4, err, printKind(result4))

	if !reflect.DeepEqual(result3, result4) {
		t.Error(result3, result4)
	}
}

func printKind(result interface{}) reflect.Kind {
	if result == nil {
		return reflect.Invalid
	}
	return reflect.TypeOf(result).Kind()
}

func mockErr() error {
	return nil
}

type cacheDataStruct struct {
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

func rawGetFunc4(ctx context.Context, p1, p2 string) (cacheDataStruct, error) {
	return cacheDataStruct{
		PersonId:   p1,
		Subject:    p2,
		NotifyType: 2,
		Amount:     19.55,
	}, mockErr()
}

func rawGetFunc5(ctx context.Context, p1, p2 string) (*cacheDataStruct, error) {
	result, _ := rawGetFunc4(ctx, p1, p2)
	return &result, mockErr()
}

func isRdsRun() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", time.Second*2)
	if err != nil {
		fmt.Println("Please install redis on local and start at port: 6379, then run test.")
		return false
	}
	conn.Close()

	return true
}

func TestMain(m *testing.M) {
	if !isRdsRun() {
		return
	}
	m.Run()
}
