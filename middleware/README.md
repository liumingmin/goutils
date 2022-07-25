

<!-- toc -->

- [middleware 中间件](#middleware-%E4%B8%AD%E9%97%B4%E4%BB%B6)
  * [limit_conn_test.go 限连接模块](#limit_conn_testgo-%E9%99%90%E8%BF%9E%E6%8E%A5%E6%A8%A1%E5%9D%97)
    + [TestLimitConn](#testlimitconn)
  * [limit_req_test.go 限流模块](#limit_req_testgo-%E9%99%90%E6%B5%81%E6%A8%A1%E5%9D%97)
    + [TestLimitReq](#testlimitreq)
  * [service_handler_test.go service封装器](#service_handler_testgo-service%E5%B0%81%E8%A3%85%E5%99%A8)
    + [TestServiceHandler](#testservicehandler)

<!-- tocstop -->

# middleware 中间件
## limit_conn_test.go 限连接模块
### TestLimitConn
```go

router := gin.New()
lr := NewLimitConn(reqHostIp)

router.Use(lr.Incoming(nil, 10, 4))
router.GET("/testurl", func(c *gin.Context) {
	time.Sleep(time.Second)
	fmt.Println("enter")
	c.String(http.StatusOK, "ok!!")
}, lr.Leaving(nil))

safego.Go(func() {
	router.Run(":8081")
})

time.Sleep(time.Second * 3)

for j := 0; j < 5; j++ {
	time.Sleep(time.Second * 1)
	for i := 0; i < 20; i++ {
		safego.Go(func() {
			resp, err := http.Get("http://127.0.0.1:8081/testurl")
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
```
## limit_req_test.go 限流模块
### TestLimitReq
```go

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
```
## service_handler_test.go service封装器
### TestServiceHandler
```go

router := gin.New()
router.POST("/foo", ServiceHandler(serviceFoo, fooReq{}, nil))

router.Run(":8080")
```
