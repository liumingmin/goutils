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

	CopyStructDefault(f1, &f2)

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

	f3 := &DstFoo{
		A: 1,
		B: []*string{&str1, &str2},
		C: map[string]*int{
			"A": &int1,
			"B": &int2,
		},
		DstUnique:        1,
		SameNameDiffType: time.Now().Format(STRUCT_DATE_TIME_FORMAT_LAYOUT),
	}
	var f4 SrcFoo
	CopyStruct(f3, &f4, BaseConvert)

	if !reflect.DeepEqual(f3.A, f4.A) {
		t.Error(f4)
	}

	if !reflect.DeepEqual(f3.B, f4.B) {
		t.Error(f4)
	}

	if !reflect.DeepEqual(f3.C, f4.C) {
		t.Error(f4)
	}

	f3Time, _ := time.ParseInLocation(STRUCT_DATE_TIME_FORMAT_LAYOUT, f3.SameNameDiffType, time.Local)
	if !reflect.DeepEqual(f3Time, f4.SameNameDiffType) {
		t.Error(f4)
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

	var f3 []*DstFoo
	CopyStructsDefault(f1, &f3)

	if !reflect.DeepEqual(f1[0].A, f3[0].A) {
		t.Error(f3)
	}

	if !reflect.DeepEqual(f1[0].B, f3[0].B) {
		t.Error(f3)
	}

	if !reflect.DeepEqual(f1[0].C, f3[0].C) {
		t.Error(f3)
	}

	if !reflect.DeepEqual(BaseConvert(f1[0].SameNameDiffType, reflect.TypeOf("")), f3[0].SameNameDiffType) {
		t.Error(f3)
	}
}
