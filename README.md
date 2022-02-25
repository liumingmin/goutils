# github.com/liumingmin/goutils

## 模块树
```
├── cache
│   ├── mem_cache.go                 内存实现的函数缓存
│   ├── rds_cache.go                 goredis实现的函数缓存
├── conf
│   └── conf.go                      YAML读取
├── container
│   ├── bitmap.go                    比特位表
│   ├── buffer_invoker.go            缓冲异步调用
│   ├── const_hash.go                一致性HASH32位
│   ├── mdb.go                       轻量级内存表
│   ├── lighttimer.go                轻量级计时器
├── distlock
│   ├── consullock.go                consul实现的分布式锁
│   ├── filelock.go                  Linux文件锁
│   ├── lock.go                      锁接口
│   ├── rdslock.go                   redis实现分布式锁
├── middleware
│   ├── captcha.go                   验证码中间件
│   ├── limit_conn.go                限连接
│   ├── limit_req.go                 限流
│   ├── service_handler.go           封装controller功能
│   ├── thumb_image.go               缩略图
├── db
│   ├── mongo
│   │   ├── client.go                    官方client封装
│   │   ├── collection.go                官方主从方式collection封装
│   ├── redis
│   │   ├── redis.go                     goredis封装
├── net
│   ├── httpx
│   │   └── httpclientx.go           httpclientx兼容1.x和2.0
└── utils
│   ├── cbk
│   │   ├── cbk.go                   熔断接口
│   │   ├── cbk_simple.go            熔断实现
│   ├── fsm
│   │   └── fsm.go                   状态机
│   ├── safego
│   │   ├── safego.go                安全的goruntine
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

js

    protoc --js_out=library=protobuf,binary:ws/js  ws/msg.proto