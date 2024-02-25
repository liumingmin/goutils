//go:build windows
// +build windows

package utils

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

var ErrUnsupportArg = errors.New("unsupport argument")

type DllMod struct {
	mod     *syscall.LazyDLL
	funcMap sync.Map
}

func NewDllMod(name string) *DllMod {
	return &DllMod{
		mod: syscall.NewLazyDLL(name),
	}
}

func (d *DllMod) Call(funcName string, args ...any) (retCode uintptr, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Call panic: %v", e)
		}
	}()

	var modFunc *syscall.LazyProc
	value, ok := d.funcMap.Load(funcName)
	if ok {
		modFunc = value.(*syscall.LazyProc)
	} else {
		modFunc = d.mod.NewProc(funcName)
		err = modFunc.Find()
		if err != nil {
			return 0, err
		}
		d.funcMap.Store(funcName, modFunc)
	}

	dllArgs := make([]uintptr, len(args))

	for i, arg := range args {
		dllArg, err := d.convertArg(arg)
		if err != nil {
			return 0, err
		}

		dllArgs[i] = dllArg
	}

	retCode, _, err = modFunc.Call(dllArgs...)
	if err != nil {
		var errno syscall.Errno
		if ok := errors.As(err, &errno); ok && errno == 0 {
			return retCode, nil
		}
	}

	return retCode, err
}

func (d *DllMod) convertArg(arg any) (uintptr, error) {
	argValue := reflect.ValueOf(arg)
	kind := argValue.Kind()
	switch kind {
	case reflect.Bool:
		b := arg.(bool)
		if b {
			return uintptr(1), nil
		} else {
			return uintptr(0), nil
		}
	case reflect.Int:
		return uintptr(arg.(int)), nil
	case reflect.Int8:
		return uintptr(arg.(int8)), nil
	case reflect.Int16:
		return uintptr(arg.(int16)), nil
	case reflect.Int32:
		return uintptr(arg.(int32)), nil
	case reflect.Int64:
		return uintptr(arg.(int64)), nil
	case reflect.Uint:
		return uintptr(arg.(uint)), nil
	case reflect.Uint8:
		return uintptr(arg.(uint8)), nil
	case reflect.Uint16:
		return uintptr(arg.(uint16)), nil
	case reflect.Uint32:
		return uintptr(arg.(uint32)), nil
	case reflect.Uint64:
		return uintptr(arg.(uint64)), nil
	case reflect.Uintptr:
		return uintptr(arg.(uintptr)), nil
	case reflect.Pointer:
		return d.convertArgPtr(argValue)
	case reflect.UnsafePointer:
		return uintptr(arg.(unsafe.Pointer)), nil
	case reflect.String:
		bsptr, _ := syscall.BytePtrFromString(arg.(string))
		return uintptr(unsafe.Pointer(bsptr)), nil
	case reflect.Slice, reflect.Array:
		if bs, ok := arg.([]byte); ok {
			if len(bs) > 0 {
				return uintptr(unsafe.Pointer(&bs[0])), nil
			} else {
				return 0, nil
			}
		}
		return 0, ErrUnsupportArg
	case reflect.Struct, reflect.Interface:
		if argValue.CanAddr() {
			return argValue.Pointer(), nil
		}
		return 0, ErrUnsupportArg
	case reflect.Func:
		return syscall.NewCallback(arg), nil

	case reflect.Float32, //ErrUnsupportArg
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Chan,
		reflect.Map:
	}
	return 0, ErrUnsupportArg
}

func (d *DllMod) convertArgPtr(argValue reflect.Value) (uintptr, error) {
	kind := argValue.Elem().Kind()
	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Struct:
		return argValue.Pointer(), nil
	}
	return 0, ErrUnsupportArg
}

func (d *DllMod) GetCStrFromUintptr(sPtr uintptr) string {
	bytes := make([]byte, 0)
	for {
		b := *(*uint8)(unsafe.Pointer(sPtr))
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
		sPtr += unsafe.Sizeof(b)
	}
	return string(bytes)
}
