package distlock

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/liumingmin/goutils/db/redis"
)

func TestRdsLock(t *testing.T) {
	redis.InitRedises()
	l, err := NewRdsLuaLock("rdscdb", "accoutId", 4)
	if err != nil {
		t.Error(err)
	}

	l2, err := NewRdsLuaLock("rdscdb", "accoutId", 4)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	if !l.Lock(ctx, 1) {
		t.Error("can not get lock")
	}

	time.Sleep(time.Millisecond * 300)
	if l2.Lock(ctx, 1) {
		t.Error("except get lock")
	}
	l.Unlock(ctx)

	time.Sleep(time.Millisecond * 100)

	if !l2.Lock(ctx, 1) {
		t.Error("can not get lock")
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
