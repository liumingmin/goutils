
# ws模块用法

## go客户端
```shell script
protoc --go_out=. ws/msg.proto
```


## js客户端使用lib模式和commonjs模式
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

https://www.npmjs.com/package/google-protobuf

## ts客户端使用commonjs模式
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

## cpp客户端
```shell script
#1. unzip cpp/protobuf.zip (download from https://github.com/protocolbuffers/protobuf/releases  sourcecode: protobuf-cpp-3.21.12.zip then build)
#2. gen compatible protobuf cpp code
cpp\protobuf\bin\protoc --cpp_out=cpp/QWS msg.proto

#build sln
```

## go服务端demo
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

e := gin.New()
e.GET("/join", func(ctx *gin.Context) {
    connMeta := ConnectionMeta{
        UserId:   ctx.DefaultQuery("uid", ""),
        Typed:    0,
        DeviceId: "",
        Version:  0,
        Charset:  0,
    }
    _, err := AcceptGin(ctx, connMeta, DebugOption(true),
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
go e.Run(":8003")
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