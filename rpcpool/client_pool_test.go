package rpcpool

import (
	//"fmt"

	"net"
	"sync"
	"testing"
)

type Args struct {
	A, B int
}

func BenchmarkPool(b *testing.B) {
	pool := &Pool{}
	pool.Init(Option{5, 1000, true},
		func() (net.Conn, error) {
			return net.Dial("tcp", "127.0.0.1:12340")
		})

	args := &Args{7, 8}
	var reply int

	b.StopTimer()
	b.StartTimer()

	b.N = 10000

	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c, err := pool.Get(); err == nil {
				err = c.CallWithTimeout("SArith.Multiply", args, &reply)
				if err != nil {
					b.Log("arith error:", err)
				}
				c.Release()
			}
		}()
	}

	wg.Wait()
	b.StopTimer()

	//fmt.Println(len(pool.clientIdles))
}

func TestPool1(t *testing.T) {
	pool := &Pool{}
	pool.Init(Option{3, 5, true},
		func() (net.Conn, error) {
			return net.Dial("tcp", "127.0.0.1:12340")
		})

	t.Log("init end")
	args := &Args{7, 8}
	var reply int

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c, err := pool.Get(); err == nil {
				err = c.CallWithTimeout("SArith.Multiply", args, &reply)
				if err != nil {
					t.Log("arith error:", err)
				}
				c.Release()
			}
		}()
	}

	wg.Wait()
}
