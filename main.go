package main

import (
	"time"
	"github.com/liumingmin/goutils/async"
	"fmt"
)

func main(){
	testAsyncInvoke()

	time.Sleep(time.Hour)
}

func testAsyncInvoke(){
	var respInfos []string
	result := async.AsyncInvokeWithTimeout(time.Second*4, func() {
		time.Sleep(time.Second*2)
		respInfos = []string{"we add1","we add2"}
		fmt.Println("1done")
	},func() {
		time.Sleep(time.Second*3)
		fmt.Println("willpanic")
		panic(nil)
		//respInfos = append(respInfos,"we add3")
		fmt.Println("2done")
	})

	fmt.Println("1alldone:",result)
}