package main

import (
	"time"
	"github.com/liumingmin/goutils/async"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/liumingmin/goutils/circuitbreaker"
	"encoding/json"
	"strings"
)

func main(){
	testCircuitBreaker()
	
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

func testCircuitBreaker(){
	router := gin.New()
	router.Use(circuitbreaker.CircuitBreaker(circuitbreaker.CircuitBreakerOptions{MaxQps:100,ReqTagFunc:reqTag1}))
	router.GET("/testurl", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.String(http.StatusOK,"ok!!")
	})

	router.Run(":8080")
}

func reqTag(c *gin.Context) string {
	keyValue := ""
	reqMap := make(map[string]interface{})
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&reqMap); err == nil {
		if value,isok := reqMap["personId"];isok{
			keyValue=value.(string)
		}
	}

	return keyValue
}

func reqTag1(c *gin.Context) string {

	return strings.Split(c.Request.RemoteAddr,":")[0]
}