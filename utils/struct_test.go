package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestCopyStruct(t *testing.T) {
	type SrcFoo struct {
		A                int
		B                []*string
		C                map[string]*int
		SrcUnique        string
		SameNameDiffType time.Time
	}
	type DstFoo struct {
		A                int
		B                []*string
		C                map[string]*int
		DstUnique        int
		SameNameDiffType string
	}

	// Create the initial value
	str1 := "hello"
	str2 := "bye bye"
	int1 := 1
	int2 := 2
	f1 := &SrcFoo{
		A: 1,
		B: []*string{&str1, &str2},
		C: map[string]*int{
			"A": &int1,
			"B": &int2,
		},
		SrcUnique:        "unique",
		SameNameDiffType: time.Now(),
	}
	var f2 DstFoo

	CopyStruct(f1, &f2, BaseConvert)

	if !reflect.DeepEqual(f1.A, f2.A) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(f1.B, f2.B) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(f1.C, f2.C) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(BaseConvert(f1.SameNameDiffType, reflect.TypeOf("")), f2.SameNameDiffType) {
		t.Error(f2)
	}
}

func TestCopyStructs(t *testing.T) {
	type SrcFoo struct {
		A                int
		B                []*string
		C                map[string]*int
		SrcUnique        string
		SameNameDiffType time.Time
	}
	type DstFoo struct {
		A                int
		B                []*string
		C                map[string]*int
		DstUnique        int
		SameNameDiffType string
	}

	// Create the initial value
	str1 := "hello"
	str2 := "bye bye"
	int1 := 1
	int2 := 2
	f1 := []SrcFoo{{
		A: 1,
		B: []*string{&str1, &str2},
		C: map[string]*int{
			"A": &int1,
			"B": &int2,
		},
		SrcUnique:        "unique",
		SameNameDiffType: time.Now(),
	}}
	var f2 []DstFoo
	CopyStructs(f1, &f2, BaseConvert)

	if !reflect.DeepEqual(f1[0].A, f2[0].A) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(f1[0].B, f2[0].B) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(f1[0].C, f2[0].C) {
		t.Error(f2)
	}

	if !reflect.DeepEqual(BaseConvert(f1[0].SameNameDiffType, reflect.TypeOf("")), f2[0].SameNameDiffType) {
		t.Error(f2)
	}
}
