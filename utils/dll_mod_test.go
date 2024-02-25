package utils

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"
)

func TestDllCall(t *testing.T) {
	// mod := NewDllMod("machineinfo.dll")

	// result := int32(0)

	// retCode, err := mod.Call("GetDiskType", "C:", &result)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if retCode != 0 {
	// 	t.FailNow()
	// }

	// if result != 4 {
	// 	t.FailNow()
	// }
}

func TestDllConvertString(t *testing.T) {
	mod := NewDllMod("test.dll")

	testStr := "abcde很棒"
	var arg uintptr
	var err error
	arg, err = mod.convertArg(testStr)
	if err != nil {
		t.FailNow()
	}

	var slice []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Data = arg
	header.Len = len(testStr)
	header.Cap = header.Len

	if string(slice) != testStr {
		t.FailNow()
	}
}

func TestDllConvertInt(t *testing.T) {
	mod := NewDllMod("test.dll")

	var arg uintptr
	var err error
	arg, err = mod.convertArg(12345)
	if err != nil {
		t.FailNow()
	}

	if arg != 12345 {
		t.FailNow()
	}

	intptr := int(1080)
	arg, err = mod.convertArg(&intptr)
	if err != nil {
		t.FailNow()
	}

	if *(*int)(unsafe.Pointer(arg)) != intptr {
		t.FailNow()
	}

	uintptr1 := uintptr(11080)
	arg, err = mod.convertArg(&uintptr1)
	if err != nil {
		t.FailNow()
	}

	if *(*uintptr)(unsafe.Pointer(arg)) != uintptr1 {
		t.FailNow()
	}
}

func TestDllConvertBool(t *testing.T) {
	mod := NewDllMod("test.dll")

	var arg uintptr
	var err error
	arg, err = mod.convertArg(true)
	if err != nil {
		t.FailNow()
	}

	if arg != 1 {
		t.FailNow()
	}
}

func TestDllConvertSlice(t *testing.T) {
	mod := NewDllMod("test.dll")

	origSlice := []byte("testslicecvt")

	var arg uintptr
	var err error
	arg, err = mod.convertArg(origSlice)
	if err != nil {
		t.FailNow()
	}

	var slice []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Data = arg
	header.Len = len(origSlice)
	header.Cap = header.Len

	if bytes.Compare(origSlice, slice) != 0 {
		t.FailNow()
	}
}

type testDllModStruct struct {
	x1 int32
	x2 int64
	x4 uintptr
}

func TestDllConvertStructPtr(t *testing.T) {
	mod := NewDllMod("test.dll")

	s := testDllModStruct{100, 200, 300}

	var arg uintptr
	var err error
	arg, err = mod.convertArg(&s)
	if err != nil {
		t.FailNow()
	}

	s2 := *(*testDllModStruct)(unsafe.Pointer(arg))
	if s2.x1 != s.x1 || s2.x2 != s.x2 || s2.x4 != s.x4 {
		t.FailNow()
	}
}

func TestGetCStrFromUintptr(t *testing.T) {
	mod := NewDllMod("test.dll")

	testStr := "abcde很棒"
	var arg uintptr
	var err error
	arg, err = mod.convertArg(testStr)
	if err != nil {
		t.FailNow()
	}

	origStr := mod.GetCStrFromUintptr(arg)

	if testStr != origStr {
		t.FailNow()
	}
}

func TestDllConvertFunc(t *testing.T) {
	//cannot convert back
	// mod := NewDllMod("test.dll")

	// var testCallback = func(s uintptr) uintptr {
	// 	fmt.Println("test callback")
	// 	return s + 900000
	// }

	// var arg uintptr
	// var err error
	// arg, err = mod.convertArg(testCallback)
	// if err != nil {
	// 	t.FailNow()
	// }

	// callback := *(*(func(s uintptr) uintptr))(unsafe.Pointer(arg))

	// t.Log(callback(12345))
}
