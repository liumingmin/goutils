package distlock

import (
	"fmt"
	"testing"
	"time"
)

func TestRdsLock(t *testing.T) {
	l, _ := NewRdsLuaLock("accoutId", 4)
	l2, _ := NewRdsLuaLock("accoutId", 4)
	//l.Lock(15)
	//l.Unlock()

	fmt.Println(l.Lock(5))
	fmt.Println("1getlock")
	fmt.Println(l2.Lock(5))
	fmt.Println("2getlock")
	time.Sleep(time.Second * 15)

	//l2, _ := NewRdsLuaLock("accoutId", 15)

	//t.Log(l2.Lock(5))
}
