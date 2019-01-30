package cache

import (
	"encoding/json"
	"reflect"

	"errors"

	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/utils"
)

type TCacheFunc struct {
	store CacheStore
	cf    func(...interface{}) string
}

func (c *TCacheFunc) Cache(expire int64, f interface{}, args ...interface{}) (interface{}, error) {
	return CacheFunc(c.store, expire, c.cf, f, args...)
}

//support args string int  []string []int []int64 []float64   map[string]xx map[int]xx  plain struct{}
func CacheFunc(store CacheStore, expire int64, cf func(...interface{}) string,
	f interface{}, args ...interface{}) (interface{}, error) {

	ft := reflect.TypeOf(f)
	if ft.NumOut() == 0 {
		log4go.Error("CacheFunc f must have one return value")
		return nil, nil
	}

	key := "CF:" + utils.MD5(cf(args...))

	cacheVal, err := store.Get(key)
	if err != nil {
		fv := reflect.ValueOf(f)

		var argValues []reflect.Value
		for _, arg := range args {
			argValues = append(argValues, reflect.ValueOf(arg))
		}

		retValues := fv.Call(argValues)
		var result interface{}
		var err error
		if len(retValues) > 0 {
			result = retValues[0].Interface()
			if ss, err := json.Marshal(result); err == nil {
				err = store.Set(key, string(ss), expire)
			}
		}

		if len(retValues) > 1 {
			if retValues[1].IsValid() && !retValues[1].IsNil() {
				err = retValues[1].Interface().(error)
			}
		}
		return result, err
	} else {
		log4go.Debug("CacheFunc hit cache : %v", string(cacheVal))
		//fmt.Printf("CacheFunc hit cache : %v\n", ss)

		origRetType := ft.Out(0)
		retType := origRetType

		if retType.Kind() == reflect.Ptr {
			retType = retType.Elem()
		}

		retValue := reflect.New(retType)
		retValueInterface := retValue.Interface()

		var result interface{}

		err = json.Unmarshal(cacheVal, retValueInterface)
		if err == nil {
			if origRetType.Kind() == reflect.Ptr {
				result = retValueInterface
			} else {
				result = reflect.Indirect(retValue).Interface()
			}
		}

		return result, err
	}

	return nil, errors.New("CacheFunc can not process")
}
