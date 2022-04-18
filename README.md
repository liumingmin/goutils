# github.com/liumingmin/goutils
gotuils目标是快速搭建应用的辅助代码库,扫码加讨论群。

<img src="avatar.jpg" width="150" height="150" ></img>

<!-- toc -->

- [ws模块用法](#ws%E6%A8%A1%E5%9D%97%E7%94%A8%E6%B3%95)
- [常用工具库](#%E5%B8%B8%E7%94%A8%E5%B7%A5%E5%85%B7%E5%BA%93)
- [algorithm](#algorithm)
  * [crc16_test.go crc16算法](#crc16_testgo-crc16%E7%AE%97%E6%B3%95)
  * [descartes_test.go 笛卡尔组合](#descartes_testgo-%E7%AC%9B%E5%8D%A1%E5%B0%94%E7%BB%84%E5%90%88)
- [cache 缓存模块](#cache-%E7%BC%93%E5%AD%98%E6%A8%A1%E5%9D%97)
  * [mem_cache_test.go 内存缓存](#mem_cache_testgo-%E5%86%85%E5%AD%98%E7%BC%93%E5%AD%98)
  * [rds_cache_test.go Redis缓存](#rds_cache_testgo-redis%E7%BC%93%E5%AD%98)
- [conf yaml配置模块](#conf-yaml%E9%85%8D%E7%BD%AE%E6%A8%A1%E5%9D%97)
- [container 容器模块](#container-%E5%AE%B9%E5%99%A8%E6%A8%A1%E5%9D%97)
  * [bitmap_test.go 比特位表](#bitmap_testgo-%E6%AF%94%E7%89%B9%E4%BD%8D%E8%A1%A8)
  * [const_hash_test.go 一致性HASH](#const_hash_testgo-%E4%B8%80%E8%87%B4%E6%80%A7hash)
  * [lighttimer_test.go 轻量级计时器](#lighttimer_testgo-%E8%BD%BB%E9%87%8F%E7%BA%A7%E8%AE%A1%E6%97%B6%E5%99%A8)
- [db 数据库](#db-%E6%95%B0%E6%8D%AE%E5%BA%93)
  * [elasticsearch ES搜索引擎](#elasticsearch-es%E6%90%9C%E7%B4%A2%E5%BC%95%E6%93%8E)
  * [kafka kafka消息队列](#kafka-kafka%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  * [mongo mongo数据库](#mongo-mongo%E6%95%B0%E6%8D%AE%E5%BA%93)
  * [redis go-redis](#redis-go-redis)
- [log zap日志库](#log-zap%E6%97%A5%E5%BF%97%E5%BA%93)
  * [zap_test.go](#zap_testgo)
- [middleware 中间件](#middleware-%E4%B8%AD%E9%97%B4%E4%BB%B6)
  * [captcha_test.go 验证码模块](#captcha_testgo-%E9%AA%8C%E8%AF%81%E7%A0%81%E6%A8%A1%E5%9D%97)
  * [limit_conn_test.go 限连接模块](#limit_conn_testgo-%E9%99%90%E8%BF%9E%E6%8E%A5%E6%A8%A1%E5%9D%97)
  * [limit_req_test.go 限流模块](#limit_req_testgo-%E9%99%90%E6%B5%81%E6%A8%A1%E5%9D%97)
  * [service_handler_test.go service封装器](#service_handler_testgo-service%E5%B0%81%E8%A3%85%E5%99%A8)
  * [thumb_image_test.go 缩略图](#thumb_image_testgo-%E7%BC%A9%E7%95%A5%E5%9B%BE)
- [net 网络库](#net-%E7%BD%91%E7%BB%9C%E5%BA%93)
  * [httpx 兼容http1.x和2.0的httpclient](#httpx-%E5%85%BC%E5%AE%B9http1x%E5%92%8C20%E7%9A%84httpclient)
  * [ip](#ip)
  * [packet tcp包model](#packet-tcp%E5%8C%85model)
  * [proxy ssh proxy](#proxy-ssh-proxy)
  * [serverx 兼容http1.x和2.0的http server](#serverx-%E5%85%BC%E5%AE%B9http1x%E5%92%8C20%E7%9A%84http-server)
- [utils 通用工具库](#utils-%E9%80%9A%E7%94%A8%E5%B7%A5%E5%85%B7%E5%BA%93)
  * [buffer_invoker 异步调用](#buffer_invoker-%E5%BC%82%E6%AD%A5%E8%B0%83%E7%94%A8)
  * [cbk 熔断器](#cbk-%E7%86%94%E6%96%AD%E5%99%A8)
  * [csv CSV文件解析为MDB内存表](#csv-csv%E6%96%87%E4%BB%B6%E8%A7%A3%E6%9E%90%E4%B8%BAmdb%E5%86%85%E5%AD%98%E8%A1%A8)
  * [distlock 分布式锁](#distlock-%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
  * [docgen 文档自动生成](#docgen-%E6%96%87%E6%A1%A3%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90)
  * [fsm 有限状态机](#fsm-%E6%9C%89%E9%99%90%E7%8A%B6%E6%80%81%E6%9C%BA)
  * [hc httpclient工具](#hc-httpclient%E5%B7%A5%E5%85%B7)
  * [ismtp 邮件工具](#ismtp-%E9%82%AE%E4%BB%B6%E5%B7%A5%E5%85%B7)
  * [safego 安全的go协程](#safego-%E5%AE%89%E5%85%A8%E7%9A%84go%E5%8D%8F%E7%A8%8B)
  * [snowflake](#snowflake)
- [ws websocket客户端和服务端库](#ws-websocket%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%92%8C%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%BA%93)
  * [js](#js)
  * [wss_test.go](#wss_testgo)

<!-- tocstop -->

## ws模块用法
```shell script
protoc --go_out=. ws/msg.proto

//js  
protoc --js_out=library=protobuf,binary:ws/js  ws/msg.proto
```

## 常用工具库

|文件  |说明    |
|----------|--------|
|async.go|带超时异步调用|
|crc16.go |查表法crc16|
|crc16-kermit.go|算法实现crc16|
|csv_parse.go|csv解析封装|
|httputils.go|httpClient工具|
|math.go|数学库|
|models.go|反射创建对象|
|stringutils.go|字符串处理|
|struct.go|结构体工具(拷贝、合并)|
|tags.go|结构体tag工具 |                     
|utils.go|其他工具类 |  

​                     
## algorithm
### crc16_test.go crc16算法
#### TestCrc16
```go

	t.Log(Crc16([]byte("abcdefg")))
```
### descartes_test.go 笛卡尔组合
#### TestDescartes
```go

	result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})
	for _, item := range result {
		t.Log(item)
	}
```
## cache 缓存模块
### mem_cache_test.go 内存缓存
#### TestMemCacheFunc
```go

	ctx := context.Background()

	const cacheKey = "UT:%v:%v"

	var lCache = cache.New(5*time.Minute, 5*time.Minute)
	result, err := MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = MemCacheFunc(ctx, lCache, 60*time.Second, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	MemCacheDelete(ctx, lCache, cacheKey, "p1", "p2")
```
### rds_cache_test.go Redis缓存
#### TestRdscCacheFunc
```go

	redis.InitRedises()
	ctx := context.Background()

	const cacheKey = "UT:%v:%v"
	const RDSC_DB = "rdscdb"

	rds := redis.Get(RDSC_DB)

	result, err := RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc0, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc1, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc2, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc3, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc4, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc5, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc6, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToArray(result), err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc7, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", drainToMap(result), err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc8, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))

	result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc9, cacheKey, "p1", "p2")
	log.Info(ctx, "%v %v %v", result, err, printKind(result))
	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")

	//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
	//log.Info(ctx, "%v %v %v", result, err, printKind(result))
	//
	//result, err = RdsCacheFunc(ctx, rds, 60, rawGetFunc10, cacheKey, "p1", "p2")
	//log.Info(ctx, "%v %v %v", result, err, printKind(result))

	RdsDeleteCache(ctx, rds, cacheKey, "p1", "p2")
```
#### TestRdsCacheMultiFunc
```go

	redis.InitRedises()
	ctx := context.Background()
	const RDSC_DB = "rdscdb"

	rds := redis.Get(RDSC_DB)
	result, err := RdsCacheMultiFunc(ctx, rds, 30, getThingsByIds, "multikey:%s", []string{"1", "2", "5", "3", "4", "10"})
	if err == nil && result != nil {
		mapValue, ok := result.(map[string]*Thing)
		if ok {
			for key, value := range mapValue {
				log.Info(ctx, "%v===%v", key, value)
			}
		}
	}
```
## conf yaml配置模块
## container 容器模块
### bitmap_test.go 比特位表
#### TestBitmapExists
```go

	bitmap := initTestData()
	t.Log(bitmap)

	t.Log(bitmap.Exists(122))
	t.Log(bitmap.Exists(123))

	//data1 := []byte{1, 2, 4, 7}
	//data2 := []byte{0, 1, 5}

```
#### TestBitmapSet
```go

	bitmap := initTestData()

	t.Log(bitmap.Exists(1256))

	bitmap.Set(1256)

	t.Log(bitmap.Exists(1256))
```
#### TestBitmapUnionOr
```go

	bitmap := initTestData()
	bitmap2 := initTestData()
	bitmap2.Set(256)

	bitmap3 := bitmap.Union(&bitmap2)
	t.Log(bitmap3.Exists(256))

	bitmap3.Set(562)
	t.Log(bitmap3.Exists(562))

	t.Log(bitmap.Exists(562))
```
#### TestBitmapBitInverse
```go

	bitmap := initTestData()

	t.Log(bitmap.Exists(66))

	bitmap.Inverse()

	t.Log(bitmap.Exists(66))

```
### const_hash_test.go 一致性HASH
#### TestConstHash
```go


	var ringchash CHashRing

	var configs []CHashNode
	for i := 0; i < 10; i++ {
		configs = append(configs, TestNode("node"+strconv.Itoa(i)))
	}

	ringchash.Adds(configs)

	fmt.Println(ringchash.Debug())

	fmt.Println("==================================", configs)

	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))

	var configs2 []CHashNode
	for i := 0; i < 2; i++ {
		configs2 = append(configs2, TestNode("node"+strconv.Itoa(10+i)))
	}
	ringchash.Adds(configs2)
	fmt.Println("==================================")
	fmt.Println(ringchash.Debug())
	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))

	ringchash.Del("node0")

	fmt.Println("==================================")
	fmt.Println(ringchash.Debug())
	fmt.Println(ringchash.Get("jjfdsljk:dfdfd:fds"))

	fmt.Println(ringchash.Get("jjfdxxvsljk:dddsaf:xzcv"))
	//
	fmt.Println(ringchash.Get("fcds:cxc:fdsfd"))
	//
	fmt.Println(ringchash.Get("vdsafd:32:fdsfd"))

	fmt.Println(ringchash.Get("xvd:fs:xcvd"))
```
### lighttimer_test.go 轻量级计时器
#### TestStartTicks
```go

	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddTimer(time.Second*time.Duration(2), func(fireSeqNo uint) bool {
		fmt.Println("callback", fireSeqNo, "-")
		if fireSeqNo == 4 {
			return true
		}
		return false
	})

	time.Sleep(time.Hour)
```
#### TestStartTicksDeadline
```go


	//NewLightTimerPool

	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddTimerWithDeadline(time.Second*time.Duration(2), time.Now().Add(time.Second*5), func(seqNo uint) bool {
		fmt.Println("callback", seqNo, "-")
		if seqNo == 4 {
			return true
		}
		return false
	}, func(seqNo uint) bool {
		fmt.Println("end callback", seqNo, "-")
		return true
	})

	time.Sleep(time.Hour)
```
#### TestLtPool
```go

	pool := NewLightTimerPool(10, time.Millisecond)

	for i := 0; i < 100000; i++ {
		tmp := i
		pool.AddTimerWithDeadline(strconv.Itoa(tmp), time.Second*time.Duration(2), time.Now().Add(time.Second*5), func(seqNo uint) bool {
			fmt.Println("callback", tmp, "-", seqNo, "-")
			if seqNo == 4 {
				return true
			}
			return false
		}, func(seqNo uint) bool {
			fmt.Println("end callback", tmp, "-", seqNo, "-")
			return true
		})
	}

	time.Sleep(time.Second * 20)

	fmt.Println(runtime.NumGoroutine())

	time.Sleep(time.Hour)
```
#### TestStartTicks2
```go

	lt := NewLightTimer()
	lt.StartTicks(time.Millisecond)

	lt.AddCallback(time.Second*time.Duration(3), func() {
		fmt.Println("invoke once")
	})

	time.Sleep(time.Hour)
```
## db 数据库
### elasticsearch ES搜索引擎
#### es6 ES6版本API
##### es_test.go
###### TestCreateIndexByModel
```go

	InitClients()

	client := GetEsClient(testUserIndexKey)

	err := client.CreateIndexByModel(context.Background(), testUserIndexName, &MappingModel{
		Mappings: map[string]Mapping{
			testUserTypeName: {
				Dynamic: false,
				Properties: map[string]*elasticsearch.MappingProperty{
					"userId": {
						Type:  "text",
						Index: false,
					},
					"nickname": {
						Type:     "text",
						Analyzer: "standard",
						Fields: map[string]*elasticsearch.MappingProperty{
							"std": {
								Type:     "text",
								Analyzer: "standard",
								ExtProps: map[string]interface{}{
									"term_vector": "with_offsets",
								},
							},
							"keyword": {
								Type: "keyword",
							},
						},
					},
					"status": {
						Type: "keyword",
					},
					"pType": {
						Type: "keyword",
					},
				},
			},
		},
		Settings: Settings{
			IndexMappingIgnoreMalformed: true,
			NumberOfReplicas:            1,
			NumberOfShards:              3,
		},
	})

	t.Log(err)
```
###### TestEsInsert
```go

	InitClients()

	client := GetEsClient(testUserIndexKey)

	for i := 0; i < 100; i++ {
		ptype := "normal"
		if i%10 == 5 {
			ptype = "vip"
		}
		status := "valid"
		if i%30 == 2 {
			status = "invalid"
		}
		id := "000000000" + fmt.Sprint(i)
		err := client.Insert(context.Background(), testUserIndexName, testUserTypeName,
			id, testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
		t.Log(err)
	}
```
###### TestEsBatchInsert
```go

	InitClients()

	client := GetEsClient(testUserIndexKey)

	ids := make([]string, 0)
	items := make([]interface{}, 0)

	for i := 0; i < 100; i++ {
		ptype := "normal"
		if i%10 == 5 {
			ptype = "vip"
		}
		status := "valid"
		if i%30 == 2 {
			status = "invalid"
		}
		id := "x00000000" + fmt.Sprint(i)

		ids = append(ids, id)
		items = append(items, &testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
	}

	err := client.BatchInsert(context.Background(), testUserIndexName, testUserTypeName, ids, items)
	t.Log(err)
```
###### TestEsUpdateById
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(30)

	err := client.UpdateById(context.Background(), testUserIndexName, testUserTypeName,
		id, map[string]interface{}{
			"status": "invalid",
		})
	t.Log(err)

	err = client.UpdateById(context.Background(), testUserIndexName, testUserTypeName,
		id, map[string]interface{}{
			"extField": "ext1234",
		})
	t.Log(err)
```
###### TestDeleteById
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(9)

	err := client.DeleteById(context.Background(), testUserIndexName, testUserTypeName,
		id)
	t.Log(err)
```
###### TestQueryEs
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	bq := elastic.NewBoolQuery()
	bq.Must(elastic.NewMatchQuery("nickname", "超级棒"))

	var users []testUser
	total := int64(0)
	err := client.FindByModel(context.Background(), elasticsearch.QueryModel{
		IndexName: testUserIndexName,
		TypeName:  testUserTypeName,
		Query:     bq,
		Size:      5,
		Results:   &users,
		Total:     &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
```
###### TestQueryEsQuerySource
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	source := `{
		"from":0,
		"size":25,
		"query":{
			"match":{"nickname":"超级"}
		}
	}`

	var users []testUser
	total := int64(0)
	err := client.FindBySource(context.Background(), elasticsearch.SourceModel{
		IndexName: testUserIndexName,
		TypeName:  testUserTypeName,
		Source:    source,
		Results:   &users,
		Total:     &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
```
###### TestAggregateBySource
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	source := `{
		"from": 0,
		"size": 0,
		"_source": {
			"includes": [
				"status",
				"pType",
				"COUNT"
			],
			"excludes": []
		},
		"stored_fields": [
			"status",
			"pType"
		],
		"aggregations": {
			"status": {
				"terms": {
					"field": "status",
					"size": 200,
					"min_doc_count": 1,
					"shard_min_doc_count": 0,
					"show_term_doc_count_error": false,
					"order": [
						{
							"_count": "desc"
						},
						{
							"_key": "asc"
						}
					]
				},
				"aggregations": {
					"pType": {
						"terms": {
							"field": "pType",
							"size": 10,
							"min_doc_count": 1,
							"shard_min_doc_count": 0,
							"show_term_doc_count_error": false,
							"order": [
								{
									"_count": "desc"
								},
								{
									"_key": "asc"
								}
							]
						},
						"aggregations": {
							"statusCnt": {
								"value_count": {
									"field": "_index"
								}
							}
						}
					}
				}
			}
		}
	}`

	var test AggregationTest
	client.AggregateBySource(context.Background(), elasticsearch.AggregateModel{
		IndexName: testUserIndexName,
		TypeName:  testUserTypeName,
		Source:    source,
		AggKeys:   []string{"status"},
	}, &test)
	t.Log(test)
```
#### es7 ES7版本API
##### es_test.go
###### TestCreateIndexByModel
```go

	InitClients()

	client := GetEsClient(testUserIndexKey)

	err := client.CreateIndexByModel(context.Background(), testUserIndexName, &MappingModel{
		Mapping: Mapping{
			Dynamic: false,
			Properties: map[string]*elasticsearch.MappingProperty{
				"userId": {
					Type:  "text",
					Index: false,
				},
				"nickname": {
					Type:     "text",
					Analyzer: "standard",
					Fields: map[string]*elasticsearch.MappingProperty{
						"std": {
							Type:     "text",
							Analyzer: "standard",
							ExtProps: map[string]interface{}{
								"term_vector": "with_offsets",
							},
						},
						"keyword": {
							Type: "keyword",
						},
					},
				},
				"status": {
					Type: "keyword",
				},
				"pType": {
					Type: "keyword",
				},
			},
		},
		Settings: Settings{
			IndexMappingIgnoreMalformed: true,
			NumberOfReplicas:            1,
			NumberOfShards:              3,
		},
	})

	t.Log(err)
```
###### TestEsInsert
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	for i := 0; i < 100; i++ {
		ptype := "normal"
		if i%10 == 5 {
			ptype = "vip"
		}
		status := "valid"
		if i%30 == 2 {
			status = "invalid"
		}
		id := "000000000" + fmt.Sprint(i)
		err := client.Insert(context.Background(), testUserIndexName,
			id, testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
		t.Log(err)
	}
```
###### TestEsBatchInsert
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	ids := make([]string, 0)
	items := make([]interface{}, 0)

	for i := 0; i < 100; i++ {
		ptype := "normal"
		if i%10 == 5 {
			ptype = "vip"
		}
		status := "valid"
		if i%30 == 2 {
			status = "invalid"
		}
		id := "x00000000" + fmt.Sprint(i)

		ids = append(ids, id)
		items = append(items, &testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
	}

	err := client.BatchInsert(context.Background(), testUserIndexName, ids, items)
	t.Log(err)
```
###### TestEsUpdateById
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(30)

	err := client.UpdateById(context.Background(), testUserIndexName,
		id, map[string]interface{}{
			"status": "invalid",
		})
	t.Log(err)

	err = client.UpdateById(context.Background(), testUserIndexName,
		id, map[string]interface{}{
			"extField": "ext1234",
		})
	t.Log(err)
```
###### TestDeleteById
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(9)

	err := client.DeleteById(context.Background(), testUserIndexName, id)
	t.Log(err)
```
###### TestQueryEs
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)

	bq := elastic.NewBoolQuery()
	bq.Must(elastic.NewMatchQuery("nickname", "超级棒"))

	var users []testUser
	total := int64(0)
	err := client.FindByModel(context.Background(), elasticsearch.QueryModel{
		IndexName: testUserIndexName,
		Query:     bq,
		Size:      5,
		Results:   &users,
		Total:     &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
```
###### TestQueryEsQuerySource
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)
	source := `{
		"from":0,
		"size":25,
		"query":{
			"match":{"nickname":"超级"}
		}
	}`

	var users []testUser
	total := int64(0)
	err := client.FindBySource(context.Background(), elasticsearch.SourceModel{
		IndexName: testUserIndexName,
		Source:    source,
		Results:   &users,
		Total:     &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
```
###### TestAggregateBySource
```go

	InitClients()
	client := GetEsClient(testUserIndexKey)
	source := `{
		"from": 0,
		"size": 0,
		"_source": {
			"includes": [
				"status",
				"pType",
				"COUNT"
			],
			"excludes": []
		},
		"stored_fields": [
			"status",
			"pType"
		],
		"aggregations": {
			"status": {
				"terms": {
					"field": "status",
					"size": 200,
					"min_doc_count": 1,
					"shard_min_doc_count": 0,
					"show_term_doc_count_error": false,
					"order": [
						{
							"_count": "desc"
						},
						{
							"_key": "asc"
						}
					]
				},
				"aggregations": {
					"pType": {
						"terms": {
							"field": "pType",
							"size": 10,
							"min_doc_count": 1,
							"shard_min_doc_count": 0,
							"show_term_doc_count_error": false,
							"order": [
								{
									"_count": "desc"
								},
								{
									"_key": "asc"
								}
							]
						},
						"aggregations": {
							"statusCnt": {
								"value_count": {
									"field": "_index"
								}
							}
						}
					}
				}
			}
		}
	}`

	var test AggregationTest
	client.AggregateBySource(context.Background(), elasticsearch.AggregateModel{
		IndexName: testUserIndexName,
		Source:    source,
		AggKeys:   []string{"status"},
	}, &test)
	t.Log(test)
```
### kafka kafka消息队列
#### kafka_test.go
##### TestKafkaProducer
```go

	InitKafka()
	producer := GetProducer("user_producer")
	producer.Produce(&sarama.ProducerMessage{
		Topic: userTopic,
		Key:   sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
	})

	time.Sleep(time.Second * 5)
```
##### TestKafkaConsumer
```go

	InitKafka()

	consumer := GetConsumer("user_consumer")
	go func() {
		consumer.Consume(userTopic, func(msg *sarama.ConsumerMessage) error {
			fmt.Println(string(msg.Key), "=", string(msg.Value))
			return nil
		}, func(err error) {

		})
	}()

	producer := GetProducer("user_producer")
	for i := 0; i < 10; i++ {
		producer.Produce(&sarama.ProducerMessage{
			Topic: userTopic,
			Key:   sarama.ByteEncoder(fmt.Sprint(i)),
			Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		})
	}

	time.Sleep(time.Second * 5)
```
### mongo mongo数据库
#### collection_test.go
##### TestInsert
```go

	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Insert(ctx, testUser{
		UserId:   "1",
		Nickname: "超级棒",
		Status:   "valid",
		Type:     "normal",
	})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
```
##### TestUpdate
```go

	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Update(ctx, bson.M{"user_id": "1"}, bson.M{"$set": bson.M{"nick_name": "超级棒++"}})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
```
##### TestDelete
```go

	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Delete(ctx, bson.M{"user_id": "1"})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
```
### redis go-redis
#### list_test.go Redis List工具库
##### TestList
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	err := ListPush(ctx, rds, "test_list", "stringvalue")
	t.Log(err)
	ListPop(rds, []string{"test_list"}, 3600, 1000, func(key, data string) {
		fmt.Println(key, data)
	})

	err = ListPush(ctx, rds, "test_list", "stringvalue")
	t.Log(err)
	time.Sleep(time.Second * 20)
```
#### lock_test.go Redis 锁工具库
##### TestRdsAllowActionWithCD
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	cd, ok := RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
	cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
	time.Sleep(time.Second * 3)

	cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
	t.Log(cd, ok)
```
##### TestRdsAllowActionByMTs
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	cd, ok := RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
	cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
	time.Sleep(time.Second)

	cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
	t.Log(cd, ok)
```
##### TestRdsLockResWithCD
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	ok := RdsLockResWithCD(ctx, rds, "test:res", "res-1", 3)
	t.Log(ok)
	ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
	t.Log(ok)
	time.Sleep(time.Second * 4)

	ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
	t.Log(ok)
```
#### mq_test.go Redis PubSub工具库
##### TestMqPSubscribe
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()

	MqPSubscribe(ctx, rds, "testkey:*", func(channel string, data string) {
		fmt.Println(channel, data)
	}, 10)

	err := MqPublish(ctx, rds, "testkey:1", "id:1")
	t.Log(err)
	err = MqPublish(ctx, rds, "testkey:2", "id:2")
	t.Log(err)
	err = MqPublish(ctx, rds, "testkey:3", "id:3")
	t.Log(err)

	time.Sleep(time.Second * 3)
```
#### zset_test.go Redis ZSet工具库
##### TestZDescartes
```go

	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()
	dimValues := [][]string{{"dim1a", "dim1b"}, {"dim2a", "dim2b", "dim2c", "dim2d"}, {"dim3a", "dim3b", "dim3c"}}

	dt, err := csv.ReadCsvToDataTable(ctx, "data.csv", ',',
		[]string{"id", "name", "createtime", "dim1", "dim2", "dim3", "member"}, "id", []string{})
	if err != nil {
		t.Log(err)
		return
	}

	err = ZDescartes(ctx, rds, dimValues, func(strs []string) (string, map[string]int64) {
		dimData := make(map[string]int64)
		for _, row := range dt.Rows() {
			if row.String("dim1") == strs[0] &&
				row.String("dim2") == strs[1] &&
				row.String("dim3") == strs[2] {
				dimData[row.String("member")] = row.Int64("createtime")
			}
		}
		return "rds" + strings.Join(strs, "-"), dimData
	}, 1000, 30)

	t.Log(err)
```
## log zap日志库
### zap_test.go
#### TestZap
```go

	ctx := &gin.Context{}
	ctx.Set("__traceId", "aaabbbbbcccc")
	//Info(ctx, "我是日志", "name", "管理员")  //json

	Info(ctx, "我是日志2")

	//Info(ctx, "我是日志3", "name")  //json

	Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
```
#### TestZapJson
```go

	ctx := &gin.Context{}
	ctx.Set("__traceId", "aaabbbbbcccc")
	Info(ctx, "我是日志 %v", "name", "管理员") //json

	Info(ctx, "我是日志3", "管理员") //json
	Error(ctx, "我是日志3")       //json
	Log(ctx, zapcore.ErrorLevel, "日志啊")
```
#### TestPanicLog
```go

	testPanicLog()

	Info(context.Background(), "catch panic")
```
#### TestLevelChange
```go

	traceId := strings.Replace(uuid.New().String(), "-", "", -1)
	ctx := context.WithValue(context.Background(), LOG_TRADE_ID, traceId)
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())

	fmt.Println(LogLess(), "============")

	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())

	fmt.Println(LogMore(), "============")
```
## middleware 中间件
### captcha_test.go 验证码模块
#### TestVerifyCaptcha
```go

	router := gin.New()
	router.GET("/cimage", GetCaptchaImage)

	g := router.Group("/", VerifyCaptcha(func(c *gin.Context) (string, string) {
		return c.DefaultPostForm("cid", ""), c.DefaultPostForm("ccode", "")
	}))
	g.POST("/submit", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	var tplStr = `
<!doctype html>
<html>
 <body>
  <form method="post" action="/submit">
		<div><input type="hidden" name="cid" value="%s"></div>
		<div><input type="image" src="/cimage?id=%s"></div>
		<div><input type="text" name="ccode" value=""></div>
		<div><input type="submit" value="submit"></div>
  </form>
 </body>
</html>
`
	router.GET("/", func(c *gin.Context) {
		cid := GetCaptchaId()
		c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf(tplStr, cid, cid)))
	})
	router.Run(":8080")
```
### limit_conn_test.go 限连接模块
#### TestLimitConn
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
### limit_req_test.go 限流模块
#### TestLimitReq
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
### service_handler_test.go service封装器
#### TestServiceHandler
```go

	router := gin.New()
	router.POST("/foo", ServiceHandler(serviceFoo, fooReq{}, &DefaultServiceResponse{}))

	router.Run(":8080")
```
### thumb_image_test.go 缩略图
#### TestThumbImageServe
```go

	router := gin.New()
	router.Use(ThumbImageServe("/images", GinHttpFs("G:/images", false)))
	router.Run(":8080")
```
## net 网络库
### httpx 兼容http1.x和2.0的httpclient
#### httpclientx_test.go
##### TestHttpXGet
```go

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://127.0.0.1:8049")
		if err != nil {
			t.Fatal(fmt.Errorf("error making request: %v", err))
		}
		t.Log(resp.StatusCode)
		t.Log(resp.Proto)
	}
```
##### TestHttpXPost
```go

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://127.0.0.1:8881")
		if err != nil {
			t.Fatal(fmt.Errorf("error making request: %v", err))
		}
		t.Log(resp.StatusCode)
		t.Log(resp.Proto)
	}
```
### ip
### packet tcp包model
### proxy ssh proxy
#### ssh_client_test.go
##### TestSshClient
```go

	client := getSshClient(t)
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("Create session failed %v", err)
	}
	defer session.Close()

	// run command and capture stdout/stderr
	output, err := session.CombinedOutput("ls -l /data")
	if err != nil {
		t.Fatalf("CombinedOutput failed %v", err)
	}
	t.Log(string(output))
```
##### TestMysqlSshClient
```go

	client := getSshClient(t)
	defer client.Close()

	//test时候，打开，会引入mysql包
	//mysql.RegisterDialContext("tcp", func(ctx context.Context, addr string) (net.Conn, error) {
	//	return client.Dial("tcp", addr)
	//})

	db, err := sql.Open("", "")
	if err != nil {
		t.Fatalf("open db failed %v", err)
	}
	defer db.Close()

	rs, err := db.Query("select  limit 10")
	if err != nil {
		t.Fatalf("open db failed %v", err)
	}
	defer rs.Close()
	for rs.Next() {

	}
```
### serverx 兼容http1.x和2.0的http server
## utils 通用工具库
### buffer_invoker 异步调用
#### buffer_invoker_test.go
##### TestFuncBuffer
```go

	for i := 0; i < 100; i++ {
		item := strconv.Itoa(i)
		safego.Go(func() {
			fb.Invoke("1234", item)
		})
	}

	fmt.Println("for end1")

	time.Sleep(time.Second * 10)

	for i := 0; i < 100; i++ {
		item := strconv.Itoa(i)
		safego.Go(func() {
			fb.Invoke("1234", item)
		})
	}

	fmt.Println("for end2")

	time.Sleep(time.Second * 60)
```
### cbk 熔断器
#### cbk_test.go
##### TestCbkFailed
```go

	InitCbk()

	var ok bool
	var lastBreaked bool
	for j := 0; j < 200; j++ {
		i := j
		//safego.Go(func() {
		err := Impls[SIMPLE].Check("test") //30s 返回一次true尝试
		fmt.Println(i, "Check:", ok)

		if err == nil {
			time.Sleep(time.Millisecond * 10)
			Impls[SIMPLE].Failed("test")

			if i > 105 && lastBreaked {
				Impls[SIMPLE].Succeed("test")
				lastBreaked = false
				fmt.Println(i, "Succeed")
			}
		} else {
			if lastBreaked {
				time.Sleep(time.Second * 10)
			} else {
				lastBreaked = true
			}
		}
		//})
	}
```
### csv CSV文件解析为MDB内存表
#### csv_parse_test.go
##### TestReadCsvToDataTable
```go

	dt, err := ReadCsvToDataTable(context.Background(), `goutils.log`, '\t',
		[]string{"xx", "xx", "xx", "xx"}, "xxx", []string{"xxx"})
	if err != nil {
		t.Log(err)
		return
	}
	for _, r := range dt.Rows() {
		t.Log(r.Data())
	}

	rs := dt.RowsBy("xxx", "869")
	for _, r := range rs {
		t.Log(r.Data())
	}

	t.Log(dt.Row("17"))
```
### distlock 分布式锁
#### consullock_test.go
##### TestAquireConsulLock
```go

	l, _ := NewConsulLock("accountId", 10)
	//l.Lock(15)
	//l.Unlock()
	ctx := context.Background()
	fmt.Println("try lock 1")

	fmt.Println(l.Lock(ctx, 5))
	//time.Sleep(time.Second * 6)

	//fmt.Println("try lock 2")
	//fmt.Println(l.Lock(3))

	l2, _ := NewConsulLock("accountId", 10)
	fmt.Println("try lock 3")
	fmt.Println(l2.Lock(ctx, 15))

	l3, _ := NewConsulLock("accountId", 10)
	fmt.Println("try lock 4")
	fmt.Println(l3.Lock(ctx, 15))

	time.Sleep(time.Minute)
```
#### filelock_test.go
##### TestFileLock
```go

	test_file_path, _ := os.Getwd()
	locked_file := test_file_path

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			flock := NewFileLock(locked_file, false)
			err := flock.Lock()
			if err != nil {
				wg.Done()
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("output : %d\n", num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(2 * time.Second)

```
#### rdslock_test.go
##### TestRdsLock
```go

	redis.InitRedises()
	l, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
	l2, _ := NewRdsLuaLock("rdscdb", "accoutId", 4)
	//l.Lock(15)
	//l.Unlock()
	ctx := context.Background()
	fmt.Println(l.Lock(ctx, 5))
	fmt.Println("1getlock")
	fmt.Println(l2.Lock(ctx, 5))
	fmt.Println("2getlock")
	time.Sleep(time.Second * 15)

	//l2, _ := NewRdsLuaLock("accoutId", 15)

	//t.Log(l2.Lock(5))
```
### docgen 文档自动生成
#### cmd
#### docgen_test.go
##### TestGenDocTestUser
```go

	sb := strings.Builder{}
	sb.WriteString(genDocTestUserQuery())
	sb.WriteString(genDocTestUserCreate())
	sb.WriteString(genDocTestUserUpdate())
	sb.WriteString(genDocTestUserDelete())

	GenDoc(context.Background(), "用户管理", "doc/testuser.md", 2, sb.String())
```
### fsm 有限状态机
### hc httpclient工具
### ismtp 邮件工具
#### ismtp_test.go
##### TestSendEmail
```go

	emailauth := LoginAuth(
		"from",
		"xxxxxx",
		"mailhost.com",
	)

	ctype := fmt.Sprintf("Content-Type: %s; charset=%s", "text/plain", "utf-8")

	msg := fmt.Sprintf("To: %s\r\nCc: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s",
		strings.Join([]string{"target@mailhost.com"}, ";"),
		"",
		"from@mailhost.com",
		"测试",
		ctype,
		"测试")

	err := SendMail("mailhost.com:port", //convert port number from int to string
		emailauth,
		"from@mailhost.com",
		[]string{"target@mailhost.com"},
		[]byte(msg),
	)

	if err != nil {
		t.Log(err)
		return
	}

	return
```
### safego 安全的go协程
### snowflake
## ws websocket客户端和服务端库
### js
### wss_test.go
#### TestWssRun
```go

	InitServer()
	InitClient()

	e := gin.Default()
	e.GET("/join", join)
	go e.Run(":8003")

	connectWss("100")

	time.Sleep(time.Minute * 5)
```
