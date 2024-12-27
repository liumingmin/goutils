package utils

import (
	"reflect"
	"testing"
)

func TestGBK2UTF8(t *testing.T) {
	src := []byte{206, 210, 202, 199, 103, 111, 117, 116, 105, 108, 115, 49}
	utf8str, err := GBK2UTF8(src)
	if err != nil {
		t.FailNow()
	}

	if string(utf8str) != "我是goutils1" {
		t.FailNow()
	}
}

func TestUTF82GBK(t *testing.T) {
	src := []byte{230, 136, 145, 230, 152, 175, 103, 111, 117, 116, 105, 108, 115, 49}
	gbkStr, err := UTF82GBK(src)
	if err != nil {
		t.FailNow()
	}

	if !reflect.DeepEqual(gbkStr, []byte{206, 210, 202, 199, 103, 111, 117, 116, 105, 108, 115, 49}) {
		t.FailNow()
	}
}
