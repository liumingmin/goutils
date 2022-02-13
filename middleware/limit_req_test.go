package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/liumingmin/goutils/utils/safego"

	"github.com/gin-gonic/gin"
)

func TestLimitReq(t *testing.T) {
	router := gin.New()
	lr := NewLimitReq(reqHostIp)

	router.Use(lr.Incoming(nil, 10, 4))
	router.GET("/testurl", func(c *gin.Context) {
		time.Sleep(time.Second)
		fmt.Println("enter")
		c.String(http.StatusOK, "ok!!")
	})

	safego.Go(func() {
		router.Run(":8080")
	})

	time.Sleep(time.Second * 3)

	for j := 0; j < 5; j++ {
		time.Sleep(time.Second * 1)
		for i := 0; i < 20; i++ {
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
	}

	//w1 := utils.PerformTestRequest("GET", "/testurl", router)
	//if 200 == w1.Code {
	//	fmt.Println("okk")
	//}
	time.Sleep(time.Minute * 20)
}
