package utils

import (
	"strings"
	"testing"
)

func TestParseContentByTag(t *testing.T) {
	tagStr := "<a href='goutils' ></a>"
	str, nextOffset := ParseContentByTag(tagStr, "<a", ">")

	if strings.TrimSpace(str) != "href='goutils'" {
		t.Error(str)
	}

	if tagStr[nextOffset:] != "</a>" {
		t.Error(tagStr[nextOffset:])
	}
}

func TestCheckKeyValueExpected(t *testing.T) {
	keyValues := make(map[string]string)
	keyValues["gotuils1"] = "tim"
	keyValues["gotuils2"] = "jack"
	if !CheckKeyValueExpected(keyValues, "gotuils1", "eric", []string{"tim"}) {
		t.FailNow()
	}

	if CheckKeyValueExpected(keyValues, "gotuils3", "eric", []string{"tim"}) {
		t.FailNow()
	}

	if !CheckKeyValueExpected(keyValues, "gotuils3", "tim", []string{"tim"}) {
		t.FailNow()
	}

	if CheckKeyValueExpected(keyValues, "gotuils1", "eric", []string{"carry"}) {
		t.FailNow()
	}

	if CheckKeyValueExpected(keyValues, "gotuils3", "eric", []string{"carry"}) {
		t.FailNow()
	}

	if !CheckKeyValueExpected(keyValues, "gotuils3", "carry", []string{"carry"}) {
		t.FailNow()
	}
}
