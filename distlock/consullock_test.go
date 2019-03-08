package distlock

import (
	"fmt"
	"testing"
	"time"
)

func TestAquireConsulLock(t *testing.T) {
	l, _ := NewConsulLock("accountId", 10)
	//l.Lock(15)
	//l.Unlock()

	fmt.Println("try lock 1")

	fmt.Println(l.Lock(5))
	//time.Sleep(time.Second * 6)

	//fmt.Println("try lock 2")
	//fmt.Println(l.Lock(3))

	l2, _ := NewConsulLock("accountId", 10)
	fmt.Println("try lock 3")
	fmt.Println(l2.Lock(15))

	l3, _ := NewConsulLock("accountId", 10)
	fmt.Println("try lock 4")
	fmt.Println(l3.Lock(15))

	time.Sleep(time.Minute)
}
