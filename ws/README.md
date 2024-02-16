**Read this in other languages: [English](README.md), [中文](README_zh.md).**


# Introduction to the ws Module
* Interface-oriented programming, automatic connection management
* Built-in support for multiple platforms, devices, versions, and character sets
* Supports asynchronous messaging and synchronous RPC calls
* Extremely efficient binary-level protocol for minimal data transfer
* Message body object pool
* Zero-copy message data
* Multi-language support (Go/JavaScript/TypeScript/C++)

# Module Usage
## plug and play
Using GoLang to write a simple server:

0. Define request and response packet protocol IDs for the server and client.
```go
const C2S_REQ  = 2
const S2C_RESP = 3
```

1. Register a server-side connection reception route.
```go
ws.InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) 
```

2. Register a server-side message reception handler, which sends a response packet after processing.
```go
ws.RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "server recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    packet := GetPoolMessage(S2C_RESP)
    packet.SetData([]byte("server response"))
    connection.SendMsg(ctx, packet, nil)
    return nil
})
```

3. Create a listening service and start it.

```go
http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
    connMeta := ws.ConnectionMeta{
        UserId:   r.URL.Query().Get("uid"),
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

## Using GoLang to write a simple client:

1. Register a client-side connection reception route.
```go
ws.InitClient()
```

2. Register a client-side message reception handler for receiving messages from the server.
```go
ws.RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
    log.Info(ctx, "client recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
    return nil
})
```

3. Connect to the established server.
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

4. In the callback after the connection is established, use ConnEstablishHandlerOption to send messages to the server.
```go
packet := ws.GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client request"))
conn.SendMsg(context.Background(), packet, nil)
```

5. Example of sending request-response RPC calls based on WebSocket (ws).
```go
packet := GetPoolMessage(C2S_REQ)
packet.SetData([]byte("client rpc req info"))
resp, err := conn.SendRequestMsg(context.Background(), packet, nil)
if err == nil {
    log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
}
```

6. Example of sending request-response RPC calls with a timeout based on WebSocket (ws).
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

## Advanced Usage
### About protobuf
Protobuf definitions generate corresponding source code. The Git repository already includes the generated results, so this step can be skipped.


The .pb files only define top-level related structure definitions; the framework communication protocol does not use protobuf implementation. Business message structures can choose to be implemented using protobuf or JSON.


If protobuf is used, the framework can support object pool functionality.


```shell script
protoc --go_out=. ws/msg.proto
```

### Available Callable Interfaces
Following the principle of interface-oriented design, implementation is separated from definition. The def.go file contains all the functions and interfaces that users need to use.

## Other Language Clients
### JavaScript Client Usage
Supports lib mode and CommonJS mode.

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

### TypeScript Client Usage in CommonJS Mode
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

### C++ Client
```shell script
#1. unzip cpp/protobuf.zip (download from https://github.com/protocolbuffers/protobuf/releases  sourcecode: protobuf-cpp-3.21.12.zip then build)
#2. gen compatible protobuf cpp code
cpp\protobuf\bin\protoc --cpp_out=cpp/QWS msg.proto

#build sln
```

## More Comprehensive Demo Examples
```go
InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke 
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

## GoLang client-side Demo Examples
```go
InitClient()        
                                            //client invoke 
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