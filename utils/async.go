package utils

import (
	"sync"
	"time"

	"goutils/safego"
)

func AsyncInvokesWithTimeout(timeout time.Duration, fs []func()) bool {
	return AsyncInvokeWithTimeout(timeout, fs...)
}

//usage:
// var respInfos []string
//result := AsyncInvokeWithTimeout(time.Second*4, func() {
//	time.Sleep(time.Second*2)
//	respInfos = []string{"we add1","we add2"}
//	fmt.Println("1done")
//},func() {
//	time.Sleep(time.Second*1)
//	//respInfos = append(respInfos,"we add3")
//	fmt.Println("2done")
//})
//fmt.Println("1alldone:",result)
//
func AsyncInvokeWithTimeout(timeout time.Duration, args ...func()) bool {
	if len(args) == 0 {
		return false
	}

	wg := &sync.WaitGroup{}

	for _, arg := range args {
		f := arg
		wg.Add(1)
		safego.Go(func() {
			defer wg.Done()
			f()
		})
	}

	return waitInvokeTimeout(wg, timeout)
}

func waitInvokeTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	safego.Go(func() {
		defer close(c)
		wg.Wait()
	})
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}
