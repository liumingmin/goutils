package distlock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/liumingmin/goutils/redis"
)

func TestRdsLock(t *testing.T) {
	redis.InitRedises()
	l, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
	l2, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
	//l.Lock(15)
	//l.Unlock()
	ctx := context.Background()
	fmt.Println(l.Lock(ctx, 5))
	fmt.Println("1getlock")
	fmt.Println(l2.Lock(ctx, 5))
	fmt.Println("2getlock")
	time.Sleep(time.Second * 15)

	//l2, _ := NewRdsLuaLock("accoutId", 15)

	//t.Log(l2.Lock(5))
}
