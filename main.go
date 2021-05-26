package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/utils"
)

const (
	x1 = iota + 100
	x2
	x3
)

func main() {
	//testWeb()
	//fmt.Println(rune("?"[0]))
	//time.Sleep(time.Hour)

	//fmt.Println(1111)
	//log.Finest("dddd %v",2321)

	//var isOutPool int32 =2
	//result := atomic.CompareAndSwapInt32(&isOutPool, 1, 0)
	//fmt.Println(result,isOutPool)
	//time.Sleep(time.Second)

	//	sort.Reverse()

	//var x uint = 1
	//fmt.Println(x3)

	//var is = []interface{}{"1", "2"}

	//var iis interface{}
	//iis = ""
	//
	//v, ok := iis.([]interface{})
	//fmt.Println(fmt.Sprintf("%v:%v", "aaa", "bbb", "cccc"))

	fmt.Println(fmt.Sprint(222))
}

func testAsyncInvoke() {
	var respInfos []string
	result := utils.AsyncInvokeWithTimeout(time.Second*4, func() {
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
