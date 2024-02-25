package redis

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestSentinel(t *testing.T) {
	InitRedises()
	rds := Get("rds-sentinel")
	ctx := context.Background()

	rds.Set(ctx, "test_senti", "test_value", time.Minute)

	_, err := rds.Get(ctx, "test_senti").Result()
	if err != nil {
		t.Error(err)
	}
}

func TestCluster(t *testing.T) {
	InitRedises()
	rds := Get("rds-cluster")
	ctx := context.Background()

	rds.Set(ctx, "test_cluster", "test_value", time.Minute)

	_, err := rds.Get(ctx, "test_cluster").Result()
	if err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", time.Second*2)
	if err != nil {
		fmt.Println("Please install redis on local and start at port: 6379, then run test.")
		return
	}
	conn.Close()

	m.Run()
}
