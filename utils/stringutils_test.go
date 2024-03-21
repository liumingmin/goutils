package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestStringsReverse(t *testing.T) {
	var strs = []string{"1", "2", "3", "4"}
	revStrs := StringsReverse(strs)

	if !reflect.DeepEqual(revStrs, []string{"4", "3", "2", "1"}) {
		t.FailNow()
	}
}

func TestStringsInArray(t *testing.T) {
	var strs = []string{"1", "2", "3", "4"}
	ok, index := StringsInArray(strs, "3")
	if !ok {
		t.FailNow()
	}

	if index != 2 {
		t.FailNow()
	}

	ok, index = StringsInArray(strs, "5")
	if ok {
		t.FailNow()
	}

	if index != -1 {
		t.FailNow()
	}
}

func TestStringsExcept(t *testing.T) {
	var strs1 = []string{"1", "2", "3", "4"}
	var strs2 = []string{"3", "4", "5", "6"}

	if !reflect.DeepEqual(StringsExcept(strs1, strs2), []string{"1", "2"}) {
		t.FailNow()
	}

	if !reflect.DeepEqual(StringsExcept(strs1, []string{}), []string{"1", "2", "3", "4"}) {
		t.FailNow()
	}

	if !reflect.DeepEqual(StringsExcept([]string{}, strs2), []string{}) {
		t.FailNow()
	}
}

func TestStringsDistinct(t *testing.T) {
	var strs1 = []string{"1", "2", "3", "4", "1", "3"}
	distincted := StringsDistinct(strs1)
	sort.Strings(distincted)
	if !reflect.DeepEqual(distincted, []string{"1", "2", "3", "4"}) {
		t.FailNow()
	}
}
