package rpcpool2

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkPool(b *testing.B) {
	pool, _ := NewPool(&Option{Addr: "127.0.0.1:12340", Size: 200,
		KeepAlive: time.Hour * 4})

	args := &Args{7, 8}
	var reply int

	b.StopTimer()
	b.StartTimer()

	b.N = 100000

	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c, err := pool.Get(); err == nil {
				err = c.CallWithTimeout("SArith.Multiply", args, &reply)
				if err != nil {
					//b.Log("arith error:", err)
				}
				c.Release()
			}
		}()
	}

	wg.Wait()
	b.StopTimer()

	//fmt.Println(len(pool.clientIdles))
}

type Args struct {
	A, B int
}

func TestPool1(t *testing.T) {
	pool, _ := NewPool(&Option{Addr: "127.0.0.1:12340", Size: 5,
		KeepAlive: time.Hour * 4})

	t.Log("init end")
	args := &Args{7, 8}
	var reply int

	for {
		wg := &sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if c, err := pool.Get(); err == nil {
					err = c.CallWithTimeout("SArith.Multiply", args, &reply)
					if err != nil {
						fmt.Println("arith error:", err)
					} else {
						fmt.Println(reply)
					}
					c.Release()
				}

			}()
		}

		wg.Wait()
	}
}
