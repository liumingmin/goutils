package redis

import (
	"context"
	"testing"
	"time"
)

func TestRdsAllowActionWithCD(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	cd, ok := RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
	cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
	time.Sleep(time.Second * 3)

	cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
}

func TestRdsAllowActionByMTs(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	cd, ok := RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
	cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
	time.Sleep(time.Second)

	cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
}

func TestRdsLockResWithCD(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	ok := RdsLockResWithCD(ctx, rds, "test:res", "res-1", 3)
	t.Log(ok)
	ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
	t.Log(ok)
	time.Sleep(time.Second * 4)

	ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
	t.Log(ok)
}
