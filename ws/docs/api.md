# ws Module API Documentation

## Quick Start

### Server

```go
import "github.com/liumingmin/goutils/ws"

// 1. Register message handler
ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    log.Info(ctx, "Received message: %v", string(msg.GetData()))
    return nil
})

// 2. Initialize server
ws.InitServer()

// 3. Accept connections in HTTP Handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
    meta := ws.ConnectionMeta{
        UserId:   "user123",
        Typed:    1,
        DeviceId: "device-abc",
    }
    conn, err := ws.Accept(r.Context(), w, r, meta)
    if err != nil {
        log.Error(r.Context(), "Connection failed: %v", err)
        return
    }
    log.Info(r.Context(), "Connected: %v", conn.Id())
}
```

### Client

```go
import "github.com/liumingmin/goutils/ws"

// 1. Register message handler
ws.RegisterHandler(1002, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    log.Info(ctx, "Received server message: %v", string(msg.GetData()))
    return nil
})

// 2. Initialize client
ws.InitClient()

// 3. Connect to server
conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws", nil)
if err != nil {
    log.Error(ctx, "Connection failed: %v", err)
    return
}

// 4. Send message
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
conn.SendMsg(ctx, msg, nil)
```

---

## Initialization Functions

### InitServer()

Initialize the server with default configuration. Creates `ClientConnHub` to manage connections from clients.

```go
ws.InitServer()
```

### InitServerWithOpt(serverOpt ServerOption)

Initialize the server with options.

```go
ws.InitServerWithOpt(ws.ServerOption{
    HubOpts: []ws.HubOption{
        ws.HubShardOption(4), // 4 shards
    },
})
```

### InitClient()

Initialize the client. Creates `ServerConnHub` to manage connections to the server, and automatically registers the `P_DISPLACE` message handler.

```go
ws.InitClient()
```

---

## Connection Establishment

### Accept

Server accepts a WebSocket connection.

```go
func Accept(ctx context.Context, w http.ResponseWriter, r *http.Request,
    meta ConnectionMeta, opts ...ConnOption) (IConnection, error)
```

**Parameters:**

| Parameter | Type | Description |
|---|---|---|
| `ctx` | `context.Context` | Context |
| `w` | `http.ResponseWriter` | HTTP response writer |
| `r` | `*http.Request` | HTTP request |
| `meta` | `ConnectionMeta` | Connection metadata |
| `opts` | `...ConnOption` | Connection options (optional) |

**Returns:** `(IConnection, error)`

```go
conn, err := ws.Accept(ctx, w, r, ws.ConnectionMeta{
    UserId:   "user1",
    Typed:    1,
    DeviceId: "device1",
    Source:   "android",
    Charset:  ws.CHARSET_UTF8,
})
```

### DialConnect

Client connects to WebSocket server (single dial with retry).

```go
func DialConnect(ctx context.Context, sUrl string, header http.Header,
    opts ...ConnOption) (IConnection, error)
```

| Parameter | Type | Description |
|---|---|---|
| `ctx` | `context.Context` | Context (used for dial phase only) |
| `sUrl` | `string` | WebSocket URL |
| `header` | `http.Header` | Custom request headers (can be nil) |
| `opts` | `...ConnOption` | Connection options |

```go
conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws", nil,
    ws.ClientIdOption("my-client-id"),
    ws.ClientDialRetryOption(5, 2*time.Second),
)
```

### AutoReDialConnect

Client auto-reconnect, automatically re-dials after disconnection.

```go
func AutoReDialConnect(ctx context.Context, sUrl string, header http.Header,
    connInterval time.Duration, opts ...ConnOption)
```

| Parameter | Type | Description |
|---|---|---|
| `ctx` | `context.Context` | Context (can cancel auto-reconnect) |
| `sUrl` | `string` | WebSocket URL |
| `header` | `http.Header` | Custom request headers |
| `connInterval` | `time.Duration` | Reconnect interval (default 5s) |
| `opts` | `...ConnOption` | Connection options |

```go
ctx, cancel := context.WithCancel(context.Background())
go ws.AutoReDialConnect(ctx, "ws://localhost:8080/ws", nil, 3*time.Second)

// To stop
cancel()
```

---

## Core Interfaces

### IConnection

Connection interface, encapsulating all WebSocket connection operations.

```go
type IConnection interface {
    // Identity
    Id() string
    ConnType() ConnType
    UserId() string
    Type() int
    DeviceId() string
    Source() string
    Version() int
    Charset() int
    ClientIp() string

    // State
    IsStopped() bool
    IsDisplaced() bool
    RefreshDeadline()

    // Messaging
    SendMsg(ctx context.Context, payload IMessage, sc SendCallback) error
    SendRequestMsg(ctx context.Context, reqMsg IMessage, sc SendCallback) (IMessage, error)
    SendResponseMsg(ctx context.Context, respMsg IMessage, reqSn uint32, sc SendCallback) error

    // Connection control
    KickClient(displace bool)                                  // Server side
    KickServer()                                               // Client side
    DisplaceClientByIp(ctx context.Context, displaceIp string) // Server side

    // Message pulling
    GetPullChannel(pullChannelId int) (chan struct{}, bool)
    SignalPullSend(ctx context.Context, pullChannelId int) error

    // Common data storage
    GetCommDataValue(key string) (interface{}, bool)
    SetCommDataValue(key string, value interface{})
    RemoveCommDataValue(key string)
    IncrCommDataValueBy(key string, delta int)
}
```

### IMessage

Message interface.

```go
type IMessage interface {
    GetProtocolId() uint32
    GetSn() uint32
    GetData() []byte
    SetData(data []byte)
    DataMsg() IDataMessage
}
```

### IHub

Connection manager interface.

```go
type IHub interface {
    Find(id string) (IConnection, error)
    RangeConnsByFunc(f func(string, IConnection) bool)
    ConnectionIds() []string
}
```

### Puller

Message puller interface.

```go
type Puller interface {
    PullSend()
}
```

---

## Global Hubs

| Variable | Type | Description |
|---|---|---|
| `ws.ClientConnHub` | `IHub` | Server-managed collection of connections from clients |
| `ws.ServerConnHub` | `IHub` | Client-managed collection of connections to servers |

```go
// Server iterates all connections
ws.ClientConnHub.RangeConnsByFunc(func(id string, conn ws.IConnection) bool {
    log.Info(ctx, "Online connection: %v, userId: %v", id, conn.UserId())
    return true // Return false to stop iteration
})

// Find a specific connection
conn, err := ws.ClientConnHub.Find("user1-1-device1-android")

// Get all connection IDs
ids := ws.ClientConnHub.ConnectionIds()
```

---

## Message Handling

### RegisterHandler

Register a message handler, routed by `protocolId`.

```go
func RegisterHandler(protocolId uint32, h MsgHandler)
```

```go
ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    // Handle message
    log.Info(ctx, "protocolId: %v, data: %v", msg.GetProtocolId(), string(msg.GetData()))
    return nil
})
```

### RegisterDataMsgType

Register a protobuf data message type, enabling object pool support.

```go
func RegisterDataMsgType(protocolId uint32, pMsg IDataMessage)
```

```go
// After registration, GetPoolMessage will automatically create pooled protobuf messages of the corresponding type
ws.RegisterDataMsgType(1001, &pb.MyRequest{})
ws.RegisterDataMsgType(1002, &pb.MyResponse{})
```

### NewMessage

Create a normal message (without object pool).

```go
func NewMessage(protocolId uint32) IMessage
```

```go
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
```

### GetPoolMessage

Create a pooled message (uses object pool, automatically recycled after sending).

```go
func GetPoolMessage(protocolId uint32) IMessage
```

```go
msg := ws.GetPoolMessage(1001)
msg.SetData([]byte("hello"))
conn.SendMsg(ctx, msg, nil) // Automatically recycled after sending
```

---

## Type Definitions

### ConnectionMeta

Connection metadata, used to identify connections during `Accept()`.

```go
type ConnectionMeta struct {
    UserId          string // User ID
    Typed           int    // Client type enum
    DeviceId        string // Device ID
    Source          string // Connection source (e.g. "android", "ios", "web")
    Version         int    // Version number
    Charset         int    // Charset (CHARSET_UTF8 / CHARSET_GBK)
    DisableConnPool bool   // Whether to disable connection object pool
}
```

Connection ID is generated by concatenating `UserId-Typed-DeviceId-Source`. A new connection with the same ID triggers displacement.

### Constants

```go
// Connection types
const (
    CONN_KIND_CLIENT = 0  // Client connection
    CONN_KIND_SERVER = 1  // Server connection
)

// Charsets
const (
    CHARSET_UTF8 = 0
    CHARSET_GBK  = 1
)
```

### Callback Types

```go
// Message handler function
type MsgHandler func(context.Context, IConnection, IMessage) error

// Event handler function
type EventHandler func(context.Context, IConnection)

// Send callback function
type SendCallback func(ctx context.Context, c IConnection, err error)
```

---

## Connection Options (ConnOption)

All options are passed via functional options pattern and can be freely combined.

### Lifecycle Callbacks

| Function | Description | Sync/Async |
|---|---|---|
| `ConnEstablishHandlerOption(handler)` | Callback after connection established | Sync (blocks registration flow) |
| `ConnClosingHandlerOption(handler)` | Callback before connection closes | Sync (blocks unregistration flow) |
| `ConnClosedHandlerOption(handler)` | Callback after connection fully closed | Async |
| `RecvPingHandlerOption(handler)` | Callback when Ping received | Sync |
| `RecvPongHandlerOption(handler)` | Callback when Pong received | Sync |

```go
ws.Accept(ctx, w, r, meta,
    ws.ConnEstablishHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Info(ctx, "Connection established: %v", conn.Id())
    }),
    ws.ConnClosedHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Info(ctx, "Connection closed: %v", conn.Id())
    }),
)
```

### Network Parameters

| Function | Default | Description |
|---|---|---|
| `NetMaxFailureRetryOption(n)` | 10 | Max network timeout retry count |
| `NetReadWaitOption(d)` | 60s | Read timeout |
| `NetWriteWaitOption(d)` | 60s | Write timeout |
| `NetTemporaryWaitOption(d)` | 500ms | Network jitter retry wait interval |
| `MaxMessageBytesSizeOption(size)` | 512KB | Max single message size in bytes |

```go
ws.Accept(ctx, w, r, meta,
    ws.NetReadWaitOption(30*time.Second),
    ws.NetWriteWaitOption(30*time.Second),
    ws.MaxMessageBytesSizeOption(32*1024*1024), // 32MB
)
```

### Buffer and Compression

| Function | Default | Description |
|---|---|---|
| `SendBufferOption(size)` | 8 | Send buffer channel size |
| `CompressionLevelOption(level)` | 0 (no compression) | WebSocket compression level |

### Debug

| Function | Description |
|---|---|
| `DebugOption(true)` | Enable Debug log output (ping/pong/message details) |

### Server-Specific

| Function | Description |
|---|---|
| `SrvUpgraderOption(upgrader)` | Custom `websocket.Upgrader` |
| `SrvUpgraderCompressOption(true)` | Enable server-side compression |
| `SrvCheckOriginOption(fn)` | Custom origin check function |
| `SrvPullChannelsOption([]int{ch1, ch2})` | Register pull notification channels |

```go
ws.Accept(ctx, w, r, meta,
    ws.SrvUpgraderCompressOption(true),
    ws.SrvCheckOriginOption(func(r *http.Request) bool {
        return r.Header.Get("Origin") == "https://example.com"
    }),
    ws.SrvPullChannelsOption([]int{1, 2, 3}),
)
```

### Client-Specific

| Function | Description |
|---|---|
| `ClientIdOption(id)` | Custom client connection ID (default: timestamp) |
| `ClientDialOption(dialer)` | Custom `websocket.Dialer` |
| `ClientDialWssOption(url, secure)` | WSS connection config (secure=false skips certificate verification) |
| `ClientDialCompressOption(true)` | Enable client-side compression |
| `ClientDialHandshakeTimeoutOption(d)` | Handshake timeout (default 10s) |
| `ClientDialRetryOption(num, interval)` | Dial retry count and interval (default 3 times, 1s) |
| `ClientDialConnFailedHandlerOption(h)` | Dial failure callback |

```go
ws.DialConnect(ctx, "wss://example.com/ws", nil,
    ws.ClientIdOption("client-001"),
    ws.ClientDialWssOption("wss://example.com/ws", false), // Skip certificate verification
    ws.ClientDialRetryOption(5, 2*time.Second),
    ws.ClientDialConnFailedHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Error(ctx, "Connection ultimately failed")
    }),
)
```

### Hub Options

| Function | Description |
|---|---|
| `HubShardOption(cnt)` | Hub shard count (for high-concurrency scenarios) |

```go
ws.InitServerWithOpt(ws.ServerOption{
    HubOpts: []ws.HubOption{
        ws.HubShardOption(8), // 8 shards
    },
})
```

---

## Message Sending

### SendMsg

Send a normal message.

```go
func (c *Connection) SendMsg(ctx context.Context, payload IMessage, sc SendCallback) error
```

```go
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
err := conn.SendMsg(ctx, msg, func(ctx context.Context, c ws.IConnection, err error) {
    if err != nil {
        log.Error(ctx, "Send failed: %v", err)
    }
})
```

### SendRequestMsg (RPC Call)

Send a request message and wait for response.

```go
func (c *Connection) SendRequestMsg(ctx context.Context, reqMsg IMessage,
    sc SendCallback) (IMessage, error)
```

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req := ws.NewMessage(1001)
req.SetData([]byte("request"))
resp, err := conn.SendRequestMsg(ctx, req, nil)
if err != nil {
    log.Error(ctx, "RPC failed: %v", err) // May be ErrWsRpcResponseTimeout
    return
}
log.Info(ctx, "Received response: %v", string(resp.GetData()))
```

### SendResponseMsg (RPC Response)

Reply to an RPC request, must pass the request's sn.

```go
func (c *Connection) SendResponseMsg(ctx context.Context, respMsg IMessage,
    reqSn uint32, sc SendCallback) error
```

```go
ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    resp := ws.NewMessage(1002)
    resp.SetData([]byte("response"))
    return conn.SendResponseMsg(ctx, resp, msg.GetSn(), nil)
})
```

---

## Connection Control

### KickClient

Server kicks a client connection.

```go
conn.KickClient(true)  // Displacement mode: sends P_DISPLACE message
conn.KickClient(false) // Direct disconnect
```

### KickServer

Client actively disconnects.

```go
conn.KickServer()
```

### DisplaceClientByIp

Kick client by IP (used in cluster scenarios).

```go
conn.DisplaceClientByIp(ctx, "192.168.1.100")
```

### IsStopped / IsDisplaced

```go
if conn.IsStopped() {
    log.Info(ctx, "Connection is stopped")
}
if conn.IsDisplaced() {
    log.Info(ctx, "Connection was displaced")
}
```

### RefreshDeadline

Refresh connection read/write timeout.

```go
conn.RefreshDeadline()
```

---

## Message Pulling (Puller)

Suitable for scenarios where the server needs to actively push data to clients.

### Configure Pull Channels

```go
ws.Accept(ctx, w, r, meta,
    ws.SrvPullChannelsOption([]int{1, 2}), // Register channels 1 and 2
)
```

### Create Puller

```go
func NewDefaultPuller(conn IConnection, pullChannelId int,
    firstPullFunc, pullFunc func(context.Context, IConnection)) Puller
```

- `firstPullFunc`: Executed once after first connection established
- `pullFunc`: Executed each time a pull signal is received

```go
puller := ws.NewDefaultPuller(conn, 1,
    func(ctx context.Context, c ws.IConnection) {
        // First time: push historical data
        pushHistory(ctx, c)
    },
    func(ctx context.Context, c ws.IConnection) {
        // Each time: push incremental data
        pushIncrement(ctx, c)
    },
)
puller.PullSend() // Blocking execution
```

### Trigger Pull

```go
// Trigger message push from anywhere
conn.SignalPullSend(ctx, 1)
```

---

## Common Data Storage

Connection objects provide thread-safe key-value storage for associating business state.

```go
conn.SetCommDataValue("loginTime", time.Now())
conn.SetCommDataValue("score", 100)

val, ok := conn.GetCommDataValue("loginTime")
if ok {
    log.Info(ctx, "Login time: %v", val)
}

conn.IncrCommDataValueBy("score", 10) // score = 110
conn.RemoveCommDataValue("loginTime")
```

---

## Error Constants

```go
var ErrWsRpcResponseTimeout = errors.New("rpc cancel or timeout")
var ErrWsRpcWaitChanClosed  = errors.New("sn channel is closed")
```

---

## Complete Examples

### Server Complete Example

```go
package main

import (
    "context"
    "net/http"

    "github.com/liumingmin/goutils/log"
    "github.com/liumingmin/goutils/ws"
)

func main() {
    // Register message handler
    ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
        log.Info(ctx, "Received from %v: %v", conn.UserId(), string(msg.GetData()))

        // Reply
        resp := ws.NewMessage(1002)
        resp.SetData([]byte("pong"))
        return conn.SendMsg(ctx, resp, nil)
    })

    // Initialize server (with sharding)
    ws.InitServerWithOpt(ws.ServerOption{
        HubOpts: []ws.HubOption{ws.HubShardOption(4)},
    })

    // HTTP route
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        meta := ws.ConnectionMeta{
            UserId:   r.URL.Query().Get("userId"),
            Typed:    1,
            DeviceId: r.URL.Query().Get("deviceId"),
            Source:   "web",
            Charset:  ws.CHARSET_UTF8,
        }

        conn, err := ws.Accept(r.Context(), w, r, meta,
            ws.ConnEstablishHandlerOption(func(ctx context.Context, c ws.IConnection) {
                log.Info(ctx, "User %v online", c.UserId())
            }),
            ws.ConnClosedHandlerOption(func(ctx context.Context, c ws.IConnection) {
                log.Info(ctx, "User %v offline", c.UserId())
            }),
            ws.DebugOption(true),
        )
        if err != nil {
            return
        }

        // Trigger pull (example)
        conn.SignalPullSend(r.Context(), 1)
    })

    log.Info(context.Background(), "WebSocket server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Client Complete Example

```go
package main

import (
    "context"
    "time"

    "github.com/liumingmin/goutils/log"
    "github.com/liumingmin/goutils/ws"
)

func main() {
    // Register message handler
    ws.RegisterHandler(1002, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
        log.Info(ctx, "Received server message: %v", string(msg.GetData()))
        return nil
    })

    // Initialize client
    ws.InitClient()

    ctx := context.Background()

    // Connect to server
    conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws?userId=user1&deviceId=d1", nil,
        ws.ClientDialRetryOption(5, 2*time.Second),
    )
    if err != nil {
        log.Error(ctx, "Connection failed: %v", err)
        return
    }

    // Send messages
    for i := 0; i < 5; i++ {
        msg := ws.NewMessage(1001)
        msg.SetData([]byte("hello"))
        conn.SendMsg(ctx, msg, nil)
        time.Sleep(time.Second)
    }

    // RPC call
    reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    req := ws.NewMessage(1001)
    req.SetData([]byte("rpc request"))
    resp, err := conn.SendRequestMsg(reqCtx, req, nil)
    cancel()
    if err != nil {
        log.Error(ctx, "RPC failed: %v", err)
    } else {
        log.Info(ctx, "RPC response: %v", string(resp.GetData()))
    }
}
```
