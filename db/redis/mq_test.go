package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestMqPSubscribe(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	MqPSubscribe(ctx, rds, "testkey:*", func(channel string, data string) {
		fmt.Println(channel, data)
	}, 10)

	err := MqPublish(ctx, rds, "testkey:1", "id:1")
	if err != nil {
		t.Error(err)
	}
	err = MqPublish(ctx, rds, "testkey:2", "id:2")
	if err != nil {
		t.Error(err)
	}
	err = MqPublish(ctx, rds, "testkey:3", "id:3")
	if err != nil {
		t.Error(err)
	}
}
