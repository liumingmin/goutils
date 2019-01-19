package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/liumingmin/goutils/log4go"
)

func MD5(origStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(origStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//support  string int  []string []int []int64 []float64   map[string]xx map[int]xx  plain struct{}
func ConsistArgs(args ...interface{}) string {
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
		case struct{}:
		default:
			t := reflect.TypeOf(arg)
			if t.Kind() == reflect.Map {
				b.WriteString(MapToOrderStr(arg))
				continue
			}

			b.WriteString(fmt.Sprintf("%#v", arg))
		}
	}

	return b.String()
}

func MapToOrderStr(arg interface{}) string {
	var b bytes.Buffer

	value := reflect.ValueOf(arg)
	keys := value.MapKeys()
	if len(keys) == 0 {
		return ""
	}

	switch keys[0].Kind() {
	case reflect.String:
		var ss []string
		for _, key := range keys {
			ss = append(ss, key.String())
			sort.Strings(ss)
		}

		for _, s := range ss {
			elem := value.MapIndex(reflect.ValueOf(s))
			b.WriteString(s)
			b.WriteString(":")
			b.WriteString(fmt.Sprintf("%#v", elem.Interface()))
			b.WriteString(",")
		}
		break
	case reflect.Int:
	case reflect.Uint:
	case reflect.Int32:
	case reflect.Uint32:
	case reflect.Int64:
	case reflect.Uint64:
		var ss []int
		for _, key := range keys {
			ss = append(ss, int(key.Int()))
			sort.Ints(ss)
		}

		for _, s := range ss {
			elem := value.MapIndex(reflect.ValueOf(s))
			b.WriteString(strconv.Itoa(s))
			b.WriteString(":")
			b.WriteString(fmt.Sprintf("%#v", elem.Interface()))
			b.WriteString(",")
		}
		break
	}

	return b.String()
}

//检查keyname的keyvalue是否符合预期值expectKeyValues，如果不存在keyvalue，使用defaultKeyValue判断
func CheckKeyValueExpected(keyValues map[string]string, keyName, defaultKeyValue string, expectKeyValues []string) bool {
	if keyValue, exist := keyValues[keyName]; exist {
		log4go.Debug("Found keyName: %v keyValue: %v, expectValue: %+v",
			keyName, keyValue, expectKeyValues)

		if found, _ := StringsInArray(expectKeyValues, keyValue); found {
			return true
		}
	} else {
		log4go.Debug("Not Found  keyName: %v, defaultValue: %v, expectValue: %+v",
			keyName, defaultKeyValue, expectKeyValues)

		if found, _ := StringsInArray(expectKeyValues, defaultKeyValue); found {
			return true
		}
	}

	return false
}
