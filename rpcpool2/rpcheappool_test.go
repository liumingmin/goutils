package rpcpool2

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkHeapPool(b *testing.B) {
	pool, _ := NewHeapPool(&Option{Addr: "127.0.0.1:12340", Size: 200, RefSize: 50,
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
