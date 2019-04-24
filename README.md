# goutils
```├── cache
│   ├── cached_writer_gin.go
│   ├── cache_func.go             函数接口缓存 
│   ├── cache_page.go             web接口缓存 
│   ├── cache_store.go            缓存接口 
├── conf
│   └── conf.go                   YML读取
├── container
│   ├── bitmap.go                 比特位表
│   ├── buffer_invoker.go         缓冲异步调用
│   ├── chash32.go                一致性HASH32位
│   ├── chash.go                  一致性HASH16位
├── distlock
│   ├── consullock.go             consul实现的分布式锁
│   ├── filelock.go               Linux文件锁
│   ├── lock.go                   锁接口
│   ├── rdslock.go                redis实现分布式锁
├── fsm
│   └── fsm.go                    状态机
├── lighttimer
│   ├── lighttimer.go             轻量级计时器
├── middleware
│   ├── captcha.go                验证码中间件
│   ├── cbk_deprecated.go         熔断(废弃)
│   ├── cbk.go                    熔断
│   ├── limit_conn.go             限连接
│   ├── limit_req.go              限流
│   ├── thumb_image.go            缩略图
├── rpcbi
│   ├── rpcclient.go              多路复用客户端
│   ├── rpccomm.go            
│   ├── rpcserver.go              多路复用服务端
├── rpcpool
│   ├── client_pool.go            rpc连接池
├── rpcpool2
│   ├── rpcclient.go              rpc改进连接池客户端
│   ├── rpcheappool.go            rpc改进连接池(heap实现共享socket)
│   ├── rpcpool.go                rpc改进连接池
├── safego
│   ├── safego.go                 安全的goruntine
│   └── stack.go
├── tcppool
│   ├── conn.go                   tcp连接池
│   └── pool.go                   tcp连接池
└── utils
    ├── async.go                     带超时异步调用
    ├── crc16.go                     查表法crc16
    ├── crc16-kermit.go              算法实现crc16
    ├── httputils.go                 
    ├── math.go
    ├── models.go                    反射创建对象
    ├── packet.go                    二进制网络包封装
    ├── stringutils.go               字符串处理
    └── utils.go                     其他
```