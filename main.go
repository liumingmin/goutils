package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/async"
	"github.com/liumingmin/goutils/cbk"
	"net/http"
	"strings"
	"time"
	//"github.com/liumingmin/goutils/log4go"
	//"sync/atomic"
)

func main() {
	//testWeb()
	//fmt.Println(rune("?"[0]))
	//time.Sleep(time.Hour)

	//fmt.Println(1111)
	//log4go.Finest("dddd %v",2321)

	//var isOutPool int32 =2
	//result := atomic.CompareAndSwapInt32(&isOutPool, 1, 0)
	//fmt.Println(result,isOutPool)
	//time.Sleep(time.Second)

	var x uint = 1
	fmt.Println(x >> 3)
}

func testAsyncInvoke() {
	var respInfos []string
	result := async.AsyncInvokeWithTimeout(time.Second*4, func() {
		time.Sleep(time.Second * 2)
		respInfos = []string{"we add1", "we add2"}
		fmt.Println("1done")
	}, func() {
		time.Sleep(time.Second * 3)
		fmt.Println("willpanic")
		panic(nil)
		//respInfos = append(respInfos,"we add3")
		fmt.Println("2done")
	})

	fmt.Println("1alldone:", result)
}

func testCircuitBreaker() {
	router := gin.New()
	router.Use(cbk.CircuitBreaker(cbk.Options{MaxQps: 100, ReqTagFunc: reqTag1}))
	router.GET("/testurl", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.String(http.StatusOK, "ok!!")
	})

	router.Run(":8080")
}

func testWeb() {
	router := gin.New()
	//router.Use(func(c *gin.Context) {
	//	c.Abort()
	//	c.String(http.StatusServiceUnavailable,"To many requests in a second")
	//	return
	//})
	router.GET("/testurl1", func(context *gin.Context) {
		fmt.Println("brefore testurl1")
	}, func(c *gin.Context) {
		c.String(http.StatusOK, "2222")
	}, func(context *gin.Context) {
		fmt.Println("end testurl1")
	})
	router.Run(":8080")
}

func reqTag(c *gin.Context) string {
	keyValue := ""
	reqMap := make(map[string]interface{})
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&reqMap); err == nil {
		if value, isok := reqMap["personId"]; isok {
			keyValue = value.(string)
		}
	}

	return keyValue
}

func reqTag1(c *gin.Context) string {

	return strings.Split(c.Request.RemoteAddr, ":")[0]
}
