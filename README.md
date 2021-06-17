# github.com/liumingmin/goutils

## 模块树
```
├── cache
│   ├── cached_writer_gin.go
│   ├── cache_func.go             函数接口缓存 
│   ├── cache_page.go             web接口缓存 
│   ├── cache_store.go            缓存接口 
├── cache_func
│   ├── mem_cache.go              内存实现的函数缓存
│   ├── rds_cache.go              redis实现的函数缓存
├── cbk
│   ├── cbk.go                    熔断接口
│   ├── cbk_simple.go             熔断实现
├── conf
│   └── conf.go                   YML读取
├── container
│   ├── bitmap.go                 比特位表
│   ├── buffer_invoker.go         缓冲异步调用
│   ├── const_hash.go             一致性HASH32位
├── distlock
│   ├── consullock.go             consul实现的分布式锁
│   ├── filelock.go               Linux文件锁
│   ├── lock.go                   锁接口
│   ├── rdslock.go                redis实现分布式锁
├── fsm
│   └── fsm.go                    状态机
├── httpx
│   └── httpclientx.go            httpclientx兼容1.x和2.0
├── lighttimer
│   ├── lighttimer.go             轻量级计时器
├── mdb
│   ├── mdb.go                    轻量级内存表
├── middleware
│   ├── captcha.go                验证码中间件
│   ├── limit_conn.go             限连接
│   ├── limit_req.go              限流
│   ├── service_handler.go        封装controller功能
│   ├── thumb_image.go            缩略图
├── safego
│   ├── safego.go                 安全的goruntine
└── utils
    ├── async.go                     带超时异步调用
    ├── crc16.go                     查表法crc16
    ├── crc16-kermit.go              算法实现crc16
    ├── csv_parse.go                 csv解析封装                            
    ├── httputils.go                 httpClient工具
    ├── math.go                      数学库
    ├── models.go                    反射创建对象
    ├── packet.go                    二进制网络包封装
    ├── stringutils.go               字符串处理
    ├── struct.go                    结构体工具(拷贝、合并)
    ├── tags.go                      结构体tag工具     
    └── utils.go                     其他工具类
```


## ws模块

protoc --go_out=. ws/msg.proto