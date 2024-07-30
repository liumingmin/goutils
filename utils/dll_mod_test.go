//go:build windows
// +build windows

package utils

import (
	"bytes"
	"reflect"
	"testing"
	"unsafe"
)

func TestDllCall(t *testing.T) {
	mod := NewDllMod("ntdll.dll")

	info := &struct {
		osVersionInfoSize uint32
		MajorVersion      uint32
		MinorVersion      uint32
		BuildNumber       uint32
		PlatformId        uint32
		CsdVersion        [128]uint16
		ServicePackMajor  uint16
		ServicePackMinor  uint16
		SuiteMask         uint16
		ProductType       byte
		_                 byte
	}{}

	info.osVersionInfoSize = uint32(unsafe.Sizeof(*info))
	retCode, err := mod.Call("RtlGetVersion", uintptr(unsafe.Pointer(info)))
	if err != nil {
		t.Error(err)
	}

	if retCode != 0 {
		t.Error(retCode)
	}

	if info.MajorVersion == 0 {
		t.Error(info.MajorVersion)
	}

	retCode, err = mod.Call("RtlGetVersion", uintptr(unsafe.Pointer(info)))
	if err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}

	if retCode != 0 {
		t.Error(retCode)
	}

	if info.MajorVersion == 0 {
		t.Error(info.MajorVersion)
	}
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

	testDllConvertNum(t, mod, int(-1080))
	testDllConvertNum(t, mod, uint(1080))
	testDllConvertNum(t, mod, int8(-128))
	testDllConvertNum(t, mod, uint8(255))
	testDllConvertNum(t, mod, int16(-30000))
	testDllConvertNum(t, mod, uint16(30000))
	testDllConvertNum(t, mod, int32(-3000000))
	testDllConvertNum(t, mod, uint32(3000000))
	testDllConvertNum(t, mod, int64(-3000000))
	testDllConvertNum(t, mod, uint64(3000000))
	testDllConvertNum(t, mod, uintptr(11080))

	testData := 123
	up := unsafe.Pointer(&testData)
	testDllConvertNum(t, mod, up)

	testDllConvertNumPtr(t, mod, int(-1080))
	testDllConvertNumPtr(t, mod, uint(1080))
	testDllConvertNumPtr(t, mod, int8(-128))
	testDllConvertNumPtr(t, mod, uint8(255))
	testDllConvertNumPtr(t, mod, int16(-30000))
	testDllConvertNumPtr(t, mod, uint16(30000))
	testDllConvertNumPtr(t, mod, int32(-3000000))
	testDllConvertNumPtr(t, mod, uint32(3000000))
	testDllConvertNumPtr(t, mod, int64(-3000000))
	testDllConvertNumPtr(t, mod, uint64(3000000))
	testDllConvertNumPtr(t, mod, uintptr(11080))

	testDllConvertNumPtr(t, mod, float32(100.12))
	testDllConvertNumPtr(t, mod, float64(100.12))
	testDllConvertNumPtr(t, mod, complex64(100.12))
	testDllConvertNumPtr(t, mod, complex128(100.12))
}

func testDllConvertNum[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | uintptr | unsafe.Pointer](t *testing.T, mod *DllMod, num T) {
	var arg uintptr
	var err error
	arg, err = mod.convertArg(num)
	if err != nil {
		t.Error(err)
	}

	if T(arg) != num {
		t.Errorf("%v not %v", arg, num)
	}
}

func testDllConvertNumPtr[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | uintptr | float32 | float64 | complex64 | complex128](t *testing.T, mod *DllMod, num T) {
	arg, err := mod.convertArg(&num)
	if err != nil {
		t.Error(err)
	}

	addrNum := *(*T)(unsafe.Pointer(arg))

	if addrNum != num {
		t.Errorf("%v is %v not %v", arg, addrNum, num)
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

	if !bytes.Equal(origSlice, slice) {
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

func TestDllConvertUnsupport(t *testing.T) {
	mod := NewDllMod("test.dll")

	_, err := mod.convertArg(float32(11.12))
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	_, err = mod.convertArg(float64(11.12))
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	_, err = mod.convertArg(complex64(11.12))
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	_, err = mod.convertArg(complex128(11.12))
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	m := make(map[string]string)
	_, err = mod.convertArg(m)
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	c := make(chan struct{})
	_, err = mod.convertArg(c)
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	s := struct{}{}
	_, err = mod.convertArg(s)
	if err != ErrUnsupportArg {
		t.Error(err)
	}

	_, err = mod.convertArg(interface{}(s))
	if err != ErrUnsupportArg {
		t.Error(err)
	}
}

//func TestDllConvertFunc(t *testing.T) {
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
//}
