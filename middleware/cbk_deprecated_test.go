package middleware

import (
	"net/http"
	"time"

	"testing"

	"github.com/gin-gonic/gin"

	"encoding/json"

	"fmt"

	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
)

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

func TestCircuitBreaker(t *testing.T) {
	router := gin.New()
	router.Use(CircuitBreakerD(Options{MaxQps: 100, ReqTagFunc: func(c *gin.Context) string {
		ip, _ := utils.ReqHostIp(c)
		return ip
	}}))
	router.GET("/testurl", func(c *gin.Context) {
		time.Sleep(time.Second)
		c.String(http.StatusOK, "ok!!")
	})

	safego.Go(func() {
		router.Run(":8080")
	})

	http.Get("http://127.0.0.1:8080/testurl")
	time.Sleep(time.Second)

	for i := 0; i < 200; i++ {
		safego.Go(func() {
			resp, err := http.Get("http://127.0.0.1:8080/testurl")
			if err != nil {
				fmt.Println(err)
			} else {
				if 200 != resp.StatusCode {
					fmt.Println("点击太快了", resp.StatusCode)
				}
			}

		})
	}

	//w1 := utils.PerformTestRequest("GET", "/testurl", router)
	//if 200 == w1.Code {
	//	fmt.Println("okk")
	//}
	time.Sleep(time.Second * 20)
}
