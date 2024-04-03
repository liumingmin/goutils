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

	result, err := RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
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

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	err = RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
	if err != nil {
		t.Error(err)
	}

	//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
	//log.Info(ctx, "%v %v %v", result, err, printKind(result))
	//
	//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
	//log.Info(ctx, "%v %v %v", result, err, printKind(result))

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

func drainToArray(v interface{}) interface{} {
	vv := reflect.ValueOf(v)
	if vv.IsValid() && !vv.IsNil() && !vv.IsZero() {
		return reflect.ValueOf(v).Index(0).Interface()
	}
	return nil
}

func drainToMap(v interface{}) interface{} {
	vv := reflect.ValueOf(v)
	if vv.IsValid() && !vv.IsNil() && !vv.IsZero() {
		return vv.MapIndex(reflect.ValueOf("abc")).Interface()
	}
	return nil
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

func rawGetFunc0(p1, p2 string) (string, error) {
	return fmt.Sprintf("TEST:%v:%v", "p1", "p2"), mockErr()
}

func rawGetFunc1(p1, p2 string) (float64, error) {
	return 21.85, mockErr()
}

func rawGetFunc2(p1, p2 string) (int64, error) {
	return 21, mockErr()
}

func rawGetFunc3(p1, p2 string) (bool, error) {
	return true, mockErr()
}

func rawGetFunc4(p1, p2 string) (cacheDataStruct, error) {
	return cacheDataStruct{
		PersonId:   p1,
		Subject:    p2,
		NotifyType: 2,
		Amount:     19.55,
		Extra:      map[string]string{"123": "444"},
	}, mockErr()
}

func rawGetFunc5(p1, p2 string) (*cacheDataStruct, error) {
	result, _ := rawGetFunc4(p1, p2)
	return &result, mockErr()
}

func rawGetFunc6(p1, p2 string) ([]*cacheDataStruct, error) {
	result, _ := rawGetFunc4(p1, p2)
	return []*cacheDataStruct{&result}, mockErr()
}

func rawGetFunc7(p1, p2 string) (map[string]*cacheDataStruct, error) {
	result, _ := rawGetFunc4(p1, p2)
	return map[string]*cacheDataStruct{"abc": &result}, mockErr()
}

func rawGetFunc8(p1, p2 string) ([]string, error) {
	return []string{p1, p2}, mockErr()
}

func rawGetFunc9(p1, p2 string) (map[string]string, error) {
	return map[string]string{p1: p2}, mockErr()
}

//func rawGetFunc10(p1, p2 string) (*ws.P_MESSAGE, error) {
//	return &ws.P_MESSAGE{
//		ProtocolId: 100,
//		Data:       []byte("pb" + p1 + p2),
//	}, nil
//}

func TestRdsCacheMultiFunc(t *testing.T) {
	if !isRdsRun() {
		return
	}

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
}

type Thing struct {
	Id   string
	Name string
}

func getThingsByIds(ctx context.Context, ids []string) (map[string]*Thing, error) {
	return map[string]*Thing{
		"1": {Id: "1"},
		"2": {Id: "2"},
		"3": {Id: "3"},
		"4": {Id: "4"},
		"5": {Id: "5"},
	}, nil
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
