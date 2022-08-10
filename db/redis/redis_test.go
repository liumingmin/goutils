package redis

import (
	"context"
	"testing"
	"time"
)

func TestSentinel(t *testing.T) {
	InitRedises()
	rds := Get("rds-sentinel")
	ctx := context.Background()

	rds.Set(ctx, "test_senti", "test_value", time.Minute)

	value, err := rds.Get(ctx, "test_senti").Result()
	t.Log(value, err)
}

func TestCluster(t *testing.T) {
	InitRedises()
	rds := Get("rds-cluster")
	ctx := context.Background()

	rds.Set(ctx, "test_cluster", "test_value", time.Minute)

	value, err := rds.Get(ctx, "test_cluster").Result()
	t.Log(value, err)
}
