package distlock

import (
	"testing"
	"time"
)

func TestRdsLock(t *testing.T) {
	l, _ := NewRdsLuaLock("accoutId", 15)
	//l.Lock(15)
	//l.Unlock()

	t.Log(l.Lock(5))
	time.Sleep(time.Second * 5)

	l2, _ := NewRdsLuaLock("accoutId", 15)

	t.Log(l2.Lock(5))
}
