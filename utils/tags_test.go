package utils

import (
	"strings"
	"testing"
)

type testUser struct {
	UserId   string
	Nickname string
	Status   *int
	Type     string
	A        int
}

func TestAutoGenTags(t *testing.T) {
	structStrWithTag := AutoGenTags(testUser{}, map[string]TAG_STYLE{
		"json":      TAG_STYLE_SNAKE,
		"bson":      TAG_STYLE_UNDERLINE,
		"form":      TAG_STYLE_ORIG,
		"nonestyle": TAG_STYLE_NONE,
	})

	if !strings.Contains(structStrWithTag, `bson:"user_id"`) {
		t.FailNow()
	}

	if !strings.Contains(structStrWithTag, `form:"UserId"`) {
		t.FailNow()
	}

	if !strings.Contains(structStrWithTag, `json:"userId"`) {
		t.FailNow()
	}

	if !strings.Contains(structStrWithTag, `json:"status"`) {
		t.FailNow()
	}

	if strings.Contains(structStrWithTag, `nonestyle:`) {
		t.FailNow()
	}

	//t.Log(structStrWithTag)
}
