**其他语言版本: [English](README.md), [中文](README_zh.md).**


# ws模块简介

* 面向接口编程、连接自动化管理
* 多平台、多设备、多版本、多字符集内建支持
* 支持异步消息和同步RPC调用
* 二进制底层协议极限省流
* 消息包体对象池
* 消息数据零拷贝
* 多语言支持(go/js/ts/cpp)

# ws模块用法
## 开箱即用
使用GO编写一个简单的服务端步骤:

0. 在服务端和客户端定义请求和响应包协议ID
```go
const C2S_REQ  = 2
const S2C_RESP = 3
```

1. 注册一个服务端连接接收路由
```go
ws.InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) 
```

2. 注册一个服务端接收消息的处理器,处理完毕后发送回包
```go
ws.RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "server recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    packet := GetPoolMessage(S2C_RESP)
    packet.SetData([]byte("server response"))
    connection.SendMsg(ctx, packet, nil)
    return nil
})
```

3. 创建一个监听服务并启动

```go
http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
    connMeta := ws.ConnectionMeta{
        UserId:   ctx.DefaultQuery("uid", ""),
    }
    _, err := ws.Accept(ctx, w, r, connMeta, 
        ws.ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
        }),
        ws.ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
        }),
        ws.ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
        }))
    if err != nil {
        log.Error(ctx, "Accept client connection failed. error: %v", err)
        return
    }
})
http.ListenAndServe(":8003", nil)
```

## 使用GO编写一个简单的客户端步骤:

1. 注册一个客户端连接接收路由
```go
ws.InitClient()
```

2. 注册一个客户端接收服务端消息的处理器
```go
ws.RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "client recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    return nil
})
```

3. 连接到已经创建好的服务器
```go
url := "ws://127.0.0.1:8003/join?uid=100"
conn, _ := ws.DialConnect(context.Background(), url, http.Header{},
    ws.ClientIdOption("server1"),
    ws.ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
    }),
    ws.ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn closing: %v, %p", conn.Id(), conn)
    }),
    ws.ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn closed: %v, %p", conn.Id(), conn)
    }),
)
log.Info(ctx, "%v", conn)
```

4. 连接建立后的回调中ConnEstablishHandlerOption可以向服务端发送消息
```go
packet := ws.GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client request"))
conn.SendMsg(context.Background(), packet, nil)
```

5. 基于ws发送请求响应rpc调用案例
```go
packet := GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client rpc req info"))
resp, err := conn.SendRequestMsg(context.Background(), packet, nil)
if err == nil {
    log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
}
```

6. 基于ws发送带超时的请求响应rpc调用案例
```go
timeoutCtx, _ := context.WithTimeout(ctx, time.Second*5)
packet := GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client rpc req info timeout"))
resp, err := conn.SendRequestMsg(timeoutCtx, packet, nil)
if err == nil {
    log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
} else {
    log.Error(ctx, "client recv err: %v", err)
}
```

## 进阶使用
### 关于protobuf
protobuf定义生成对应的源码，Git仓库已经包含生成的结果，可跳过该步骤。

pb的文件中仅定义了顶号相关的结构定义，框架通讯的协议并不使用pb实现，业务的消息结构可选择使用pb或json来实现

如果使用pb，框架可以支持对象池功能
```shell script
protoc --go_out=. ws/msg.proto
```

### 有哪些可调用接口
根据面向接口设计原则，实现与定义分离，def.go文件包含用户需要使用的所有函数和接口。

## 其他语言客户端
### js客户端使用
lib模式和commonjs模式

https://www.npmjs.com/package/google-protobuf

```shell script
npm i google-protobuf

//lib js  (msg_pb_libs.js+google-protobuf.js)
protoc --js_out=library=msg_pb_libs,binary:ws/js  ws/msg.proto

//commonjs  (msg_pb_dist.js or msg_pb_dist.min.js)
cd ws
protoc --js_out=import_style=commonjs,binary:js  msg.proto

cd js
npm i -g browserify
npm i -g minifier
browserify msg_pb.js <custom_pb.js> -o  msg_pb_dist.js
minify msg_pb_dist.js   //msg_pb_dist.min.js

http://127.0.0.1:8003/js/demo.html
```

### ts客户端使用commonjs模式
```shell script
npm i protobufjs
npm i -g protobufjs-cli

cd ts

pbjs -t static-module -w commonjs -o dist/msg_pb.js ../msg.proto
pbts -o msg_pb.d.ts dist/msg_pb.js

tsc -p tsconfig.json
node demo.js //const WebSocket = require("ws");

npm i -g browserify

browserify dist/msg_pb.js dist/wsc.js dist/demo.js  -o dist/bundle.js

http://127.0.0.1:8003/ts/demo.html
```

### cpp客户端
```shell script
#1. unzip cpp/protobuf.zip (download from https://github.com/protocolbuffers/protobuf/releases  sourcecode: protobuf-cpp-3.21.12.zip then build)
#2. gen compatible protobuf cpp code
cpp\protobuf\bin\protoc --cpp_out=cpp/QWS msg.proto

#build sln
```

## 更加丰富的案例demo
```go
InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke 服务端调用
ctx := context.Background()

const (
    C2S_REQ  = 2
    S2C_RESP = 3
)

//server reg handler
RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "server recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    packet := GetPoolMessage(S2C_RESP)
    packet.SetData([]byte("server response"))
    connection.SendMsg(ctx, packet, nil)
    return nil
})

http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
    connMeta := ConnectionMeta{
        UserId:   r.URL.Query().Get("uid"),
        Typed:    0,
        DeviceId: "",
        Version:  0,
        Charset:  0,
    }
    _, err := Accept(ctx, w, r, connMeta, DebugOption(true),
        SrvUpgraderCompressOption(true),
        CompressionLevelOption(2),
        ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
        }),
        ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
        }),
        ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
            log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
        }))
    if err != nil {
        log.Error(ctx, "Accept client connection failed. error: %v", err)
        return
    }
})
http.ListenAndServe(":8003", nil)
```

## go客户端demo
```go
InitClient()        
                                            //client invoke 客户端调用
const (
    C2S_REQ  = 2
    S2C_RESP = 3
)

//client reg handler
RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "client recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    return nil
})

//client connect
uid := "100"
url := "ws://127.0.0.1:8003/join?uid=" + uid
conn, _ := DialConnect(context.Background(), url, http.Header{},
    DebugOption(true),
    ClientIdOption("server1"),
    ClientDialWssOption(url, false),
    ClientDialCompressOption(true),
    CompressionLevelOption(2),
    ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
    }),
    ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn closing: %v, %p", conn.Id(), conn)
    }),
    ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
        log.Info(ctx, "client conn closed: %v, %p", conn.Id(), conn)
    }),
)

log.Info(ctx, "%v", conn)
time.Sleep(time.Second * 5)

packet := GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client request"))
conn.SendMsg(context.Background(), packet, nil)
```