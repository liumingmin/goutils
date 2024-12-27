package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func AsyncInvokesWithTimeout(timeout time.Duration, fs []func()) bool {
	return AsyncInvokeWithTimeout(timeout, fs...)
}

// usage:
// var respInfos []string
//
//	result := AsyncInvokeWithTimeout(time.Second*4, func() {
//		time.Sleep(time.Second*2)
//		respInfos = []string{"we add1","we add2"}
//		fmt.Println("1done")
//	},func() {
//
//		time.Sleep(time.Second*1)
//		//respInfos = append(respInfos,"we add3")
//		fmt.Println("2done")
//	})
//
// fmt.Println("1alldone:",result)
func AsyncInvokeWithTimeout(timeout time.Duration, args ...func()) bool {
	if len(args) == 0 {
		return false
	}

	wg := &sync.WaitGroup{}

	for _, arg := range args {
		f := arg
		wg.Add(1)
		go func() {
			defer func() {
				recover()
			}()
			defer wg.Done()

			f()
		}()
	}

	return waitInvokeTimeout(wg, timeout)
}

func waitInvokeTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer func() {
			recover()
		}()
		defer close(c)
		wg.Wait()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-c:
		return true // completed normally
	case <-timer.C:
		return false // timed out
	}
}

func Sleep(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}

var globalCloseChanMutex sync.Mutex

func SafeCloseChan[T any](mutex *sync.Mutex, ch chan T) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("close chan panic, error is: %v", e)
		}
	}()

	if mutex == nil {
		mutex = &globalCloseChanMutex
	}

	mutex.Lock()
	defer mutex.Unlock()

	select {
	case _, isok := <-ch:
		if isok {
			close(ch)
		}
	default:
		close(ch)
	}

	return nil
}
