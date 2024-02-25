package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	err := ListPush(ctx, rds, "test_list", "stringvalue")
	if err != nil {
		t.Error(err)
	}
	ListPop(rds, []string{"test_list"}, 3600, 1000, func(key, data string) {
		fmt.Println(key, data)
	})

	err = ListPush(ctx, rds, "test_list", "stringvalue")
	if err != nil {
		t.Error(err)
	}
}
