package utils

import (
	"fmt"
	"testing"
)

type testUser struct {
	UserId   string
	Nickname string
	Status   string
	Type     string
}

func TestAutoGenTags(t *testing.T) {
	fmt.Println(AutoGenTags(testUser{}, map[string]TAG_STYLE{
		"json": TAG_STYLE_SNAKE,
		"bson": TAG_STYLE_UNDERLINE,
		"form": TAG_STYLE_ORIG,
	}))
}
