package utils

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestAnyIndirect(t *testing.T) {
	val := reflect.ValueOf(10)
	if AnyIndirect(val) != val {
		t.Error(val)
	}

	x := 10
	val2 := reflect.ValueOf(&x)
	if AnyIndirect(val2) == val2 {
		t.Error(val2)
	}

	if AnyIndirect(val2).Int() != int64(x) {
		t.Error(val2)
	}
}

func TestIsNil(t *testing.T) {
	var m map[string]string
	if !IsNil(m) {
		t.Error(m)
	}

	var c chan string
	if !IsNil(c) {
		t.Error(c)
	}

	var fun func()
	if !IsNil(fun) {
		t.Error("func not nil")
	}

	var s []string
	if !IsNil(s) {
		t.Error(s)
	}

	var sp *string
	if !IsNil(sp) {
		t.Error(sp)
	}

	var up unsafe.Pointer
	if !IsNil(up) {
		t.Error(up)
	}

	testIsNil[map[string]string](t)
	testIsNil[chan string](t)
	testIsNil[func()](t)
	testIsNil[[]string](t)
	testIsNil[*string](t)
	testIsNil[unsafe.Pointer](t)
}

func testIsNil[T any](t *testing.T) {
	value := testWrapperNil[T]()

	if value == nil {
		t.Error(value)
	}

	if !IsNil(value) {
		t.Error(value)
	}
}

func testWrapperNil[T any]() interface{} {
	var data T
	return data
}
