package conv

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/demdxx/gocast"
	"google.golang.org/protobuf/proto"
)

var (
	_pbIface reflect.Type
)

func init() {
	var pbMsg proto.Message
	_pbIface = reflect.TypeOf(&pbMsg).Elem()
}

func ValueToString(value interface{}) (string, error) {
	typ := reflect.TypeOf(value)

	if typ.Implements(_pbIface) {
		data, err := proto.Marshal(value.(proto.Message))
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	switch typ.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Struct, reflect.Map:
		bs, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}
	return gocast.ToString(value), nil
}

func StringToValue[T any](str string) (t T, err error) {
	if strings.TrimSpace(str) == "" {
		return
	}

	ptyp := reflect.TypeOf(&t)
	if ptyp.Implements(_pbIface) {
		var pb interface{} = &t
		err = proto.Unmarshal([]byte(str), pb.(proto.Message))
		return
	}

	typ := reflect.TypeOf(t)
	switch typ.Kind() {
	case reflect.String:
		var s interface{} = t
		return s.(T), nil
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface, reflect.Struct:
		err = json.Unmarshal([]byte(str), &t)
		return
	}

	var other interface{}
	other, err = gocast.ToT(str, typ, "")
	if err != nil {
		return
	}
	return other.(T), nil
}
