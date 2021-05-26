package cache_func

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/liumingmin/goutils/log"
)

func TestRdscCacheFunc(t *testing.T) {
	ctx := context.Background()

	const cacheKey = "UT:%v:%v"
	const RDSC_DB = "rdscdb"

	result, err := RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, RDSC_DB, 60, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, RDSC_DB, cacheKey, "p1", "p2")
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
