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

	const cacheKey = "UT:%v:%v"
	const RDSC_DB = "rdscdb"

	rds := redisDao.Get(RDSC_DB)

	result, err := RdsCacheFunc2(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	log.Info(ctx, "%v %v %v", result, err, printKind(result))
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

	var result interface{}
	var err error

	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc2(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
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
	Extra      interface{}
}

func rawGetFunc0(ctx context.Context, p1, p2 string) (string, error) {
	return fmt.Sprintf("TEST:%v:%v", "p1", "p2"), mockErr()
}

func rawGetFunc4(ctx context.Context, p1, p2 string) (cacheDataStruct, error) {
	return cacheDataStruct{
		PersonId:   p1,
		Subject:    p2,
		NotifyType: 2,
		Amount:     19.55,
		Extra:      map[string]string{"123": "444"},
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
