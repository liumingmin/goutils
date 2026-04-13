# ws 模块 API 文档

## 快速开始

### 服务端

```go
import "github.com/liumingmin/goutils/ws"

// 1. 注册消息处理器
ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    log.Info(ctx, "收到消息: %v", string(msg.GetData()))
    return nil
})

// 2. 初始化服务端
ws.InitServer()

// 3. 在 HTTP Handler 中接受连接
func wsHandler(w http.ResponseWriter, r *http.Request) {
    meta := ws.ConnectionMeta{
        UserId:   "user123",
        Typed:    1,
        DeviceId: "device-abc",
    }
    conn, err := ws.Accept(r.Context(), w, r, meta)
    if err != nil {
        log.Error(r.Context(), "连接失败: %v", err)
        return
    }
    log.Info(r.Context(), "连接成功: %v", conn.Id())
}
```

### 客户端

```go
import "github.com/liumingmin/goutils/ws"

// 1. 注册消息处理器
ws.RegisterHandler(1002, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    log.Info(ctx, "收到服务端消息: %v", string(msg.GetData()))
    return nil
})

// 2. 初始化客户端
ws.InitClient()

// 3. 连接服务端
conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws", nil)
if err != nil {
    log.Error(ctx, "连接失败: %v", err)
    return
}

// 4. 发送消息
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
conn.SendMsg(ctx, msg, nil)
```

---

## 初始化函数

### InitServer()

初始化服务端，使用默认配置。会创建 `ClientConnHub` 管理来自客户端的连接。

```go
ws.InitServer()
```

### InitServerWithOpt(serverOpt ServerOption)

使用选项初始化服务端。

```go
ws.InitServerWithOpt(ws.ServerOption{
    HubOpts: []ws.HubOption{
        ws.HubShardOption(4), // 4 个分片
    },
})
```

### InitClient()

初始化客户端。会创建 `ServerConnHub` 管理连向服务端的连接，并自动注册 `P_DISPLACE` 消息处理器。

```go
ws.InitClient()
```

---

## 连接建立

### Accept

服务端接受 WebSocket 连接。

```go
func Accept(ctx context.Context, w http.ResponseWriter, r *http.Request,
    meta ConnectionMeta, opts ...ConnOption) (IConnection, error)
```

**参数：**

| 参数 | 类型 | 说明 |
|---|---|---|
| `ctx` | `context.Context` | 上下文 |
| `w` | `http.ResponseWriter` | HTTP 响应写入器 |
| `r` | `*http.Request` | HTTP 请求 |
| `meta` | `ConnectionMeta` | 连接元数据 |
| `opts` | `...ConnOption` | 连接选项（可选） |

**返回：** `(IConnection, error)`

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

客户端连接 WebSocket 服务端（单次拨号，带重试）。

```go
func DialConnect(ctx context.Context, sUrl string, header http.Header,
    opts ...ConnOption) (IConnection, error)
```

| 参数 | 类型 | 说明 |
|---|---|---|
| `ctx` | `context.Context` | 上下文（仅用于拨号阶段） |
| `sUrl` | `string` | WebSocket URL |
| `header` | `http.Header` | 自定义请求头（可为 nil） |
| `opts` | `...ConnOption` | 连接选项 |

```go
conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws", nil,
    ws.ClientIdOption("my-client-id"),
    ws.ClientDialRetryOption(5, 2*time.Second),
)
```

### AutoReDialConnect

客户端自动重连，断线后自动重新拨号。

```go
func AutoReDialConnect(ctx context.Context, sUrl string, header http.Header,
    connInterval time.Duration, opts ...ConnOption)
```

| 参数 | 类型 | 说明 |
|---|---|---|
| `ctx` | `context.Context` | 上下文（可取消自动重连） |
| `sUrl` | `string` | WebSocket URL |
| `header` | `http.Header` | 自定义请求头 |
| `connInterval` | `time.Duration` | 重连间隔（默认 5s） |
| `opts` | `...ConnOption` | 连接选项 |

```go
ctx, cancel := context.WithCancel(context.Background())
go ws.AutoReDialConnect(ctx, "ws://localhost:8080/ws", nil, 3*time.Second)

// 需要停止时
cancel()
```

---

## 核心接口

### IConnection

连接接口，封装了 WebSocket 连接的所有操作。

```go
type IConnection interface {
    // 身份信息
    Id() string
    ConnType() ConnType
    UserId() string
    Type() int
    DeviceId() string
    Source() string
    Version() int
    Charset() int
    ClientIp() string

    // 状态
    IsStopped() bool
    IsDisplaced() bool
    RefreshDeadline()

    // 消息发送
    SendMsg(ctx context.Context, payload IMessage, sc SendCallback) error
    SendRequestMsg(ctx context.Context, reqMsg IMessage, sc SendCallback) (IMessage, error)
    SendResponseMsg(ctx context.Context, respMsg IMessage, reqSn uint32, sc SendCallback) error

    // 连接控制
    KickClient(displace bool)                                  // 服务端调用
    KickServer()                                               // 客户端调用
    DisplaceClientByIp(ctx context.Context, displaceIp string) // 服务端调用

    // 消息拉取
    GetPullChannel(pullChannelId int) (chan struct{}, bool)
    SignalPullSend(ctx context.Context, pullChannelId int) error

    // 公共数据存储
    GetCommDataValue(key string) (interface{}, bool)
    SetCommDataValue(key string, value interface{})
    RemoveCommDataValue(key string)
    IncrCommDataValueBy(key string, delta int)
}
```

### IMessage

消息接口。

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

连接管理器接口。

```go
type IHub interface {
    Find(id string) (IConnection, error)
    RangeConnsByFunc(f func(string, IConnection) bool)
    ConnectionIds() []string
}
```

### Puller

消息拉取器接口。

```go
type Puller interface {
    PullSend()
}
```

---

## 全局 Hub

| 变量 | 类型 | 说明 |
|---|---|---|
| `ws.ClientConnHub` | `IHub` | 服务端管理的来自客户端的连接集合 |
| `ws.ServerConnHub` | `IHub` | 客户端管理的连向服务端的连接集合 |

```go
// 服务端遍历所有连接
ws.ClientConnHub.RangeConnsByFunc(func(id string, conn ws.IConnection) bool {
    log.Info(ctx, "在线连接: %v, userId: %v", id, conn.UserId())
    return true // 返回 false 停止遍历
})

// 查找指定连接
conn, err := ws.ClientConnHub.Find("user1-1-device1-android")

// 获取所有连接ID
ids := ws.ClientConnHub.ConnectionIds()
```

---

## 消息处理

### RegisterHandler

注册消息处理器，按 `protocolId` 路由。

```go
func RegisterHandler(protocolId uint32, h MsgHandler)
```

```go
ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
    // 处理消息
    log.Info(ctx, "protocolId: %v, data: %v", msg.GetProtocolId(), string(msg.GetData()))
    return nil
})
```

### RegisterDataMsgType

注册 protobuf 数据消息类型，启用对象池支持。

```go
func RegisterDataMsgType(protocolId uint32, pMsg IDataMessage)
```

```go
// 注册后，GetPoolMessage 会自动创建对应类型的池化 protobuf 消息
ws.RegisterDataMsgType(1001, &pb.MyRequest{})
ws.RegisterDataMsgType(1002, &pb.MyResponse{})
```

### NewMessage

创建普通消息（不使用对象池）。

```go
func NewMessage(protocolId uint32) IMessage
```

```go
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
```

### GetPoolMessage

创建池化消息（使用对象池，发送后自动回收）。

```go
func GetPoolMessage(protocolId uint32) IMessage
```

```go
msg := ws.GetPoolMessage(1001)
msg.SetData([]byte("hello"))
conn.SendMsg(ctx, msg, nil) // 发送后自动回收
```

---

## 类型定义

### ConnectionMeta

连接元数据，用于 `Accept()` 时标识连接。

```go
type ConnectionMeta struct {
    UserId          string // 用户ID
    Typed           int    // 客户端类型枚举
    DeviceId        string // 设备ID
    Source          string // 连接来源（如 "android", "ios", "web"）
    Version         int    // 版本号
    Charset         int    // 字符集 (CHARSET_UTF8 / CHARSET_GBK)
    DisableConnPool bool   // 是否禁用连接对象池
}
```

连接 ID 由 `UserId-Typed-DeviceId-Source` 拼接生成，相同 ID 的新连接会触发顶号。

### 常量

```go
// 连接类型
const (
    CONN_KIND_CLIENT = 0  // 客户端连接
    CONN_KIND_SERVER = 1  // 服务端连接
)

// 字符集
const (
    CHARSET_UTF8 = 0
    CHARSET_GBK  = 1
)
```

### 回调类型

```go
// 消息处理函数
type MsgHandler func(context.Context, IConnection, IMessage) error

// 事件处理函数
type EventHandler func(context.Context, IConnection)

// 发送回调函数
type SendCallback func(ctx context.Context, c IConnection, err error)
```

---

## 连接选项 (ConnOption)

所有选项通过函数式选项模式传递，可自由组合。

### 生命周期回调

| 函数 | 说明 | 同步/异步 |
|---|---|---|
| `ConnEstablishHandlerOption(handler)` | 连接建立后回调 | 同步（阻塞注册流程） |
| `ConnClosingHandlerOption(handler)` | 连接关闭前回调 | 同步（阻塞注销流程） |
| `ConnClosedHandlerOption(handler)` | 连接完全关闭后回调 | 异步 |
| `RecvPingHandlerOption(handler)` | 收到 Ping 时回调 | 同步 |
| `RecvPongHandlerOption(handler)` | 收到 Pong 时回调 | 同步 |

```go
ws.Accept(ctx, w, r, meta,
    ws.ConnEstablishHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Info(ctx, "连接建立: %v", conn.Id())
    }),
    ws.ConnClosedHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Info(ctx, "连接关闭: %v", conn.Id())
    }),
)
```

### 网络参数

| 函数 | 默认值 | 说明 |
|---|---|---|
| `NetMaxFailureRetryOption(n)` | 10 | 网络超时最大重试次数 |
| `NetReadWaitOption(d)` | 60s | 读超时 |
| `NetWriteWaitOption(d)` | 60s | 写超时 |
| `NetTemporaryWaitOption(d)` | 500ms | 网络抖动重试等待间隔 |
| `MaxMessageBytesSizeOption(size)` | 512KB | 单条消息最大字节数 |

```go
ws.Accept(ctx, w, r, meta,
    ws.NetReadWaitOption(30*time.Second),
    ws.NetWriteWaitOption(30*time.Second),
    ws.MaxMessageBytesSizeOption(32*1024*1024), // 32MB
)
```

### 缓冲区与压缩

| 函数 | 默认值 | 说明 |
|---|---|---|
| `SendBufferOption(size)` | 8 | 发送缓冲区 channel 大小 |
| `CompressionLevelOption(level)` | 0 (不压缩) | WebSocket 压缩级别 |

### 调试

| 函数 | 说明 |
|---|---|
| `DebugOption(true)` | 开启 Debug 日志输出（ping/pong/消息详情） |

### 服务端专用

| 函数 | 说明 |
|---|---|
| `SrvUpgraderOption(upgrader)` | 自定义 `websocket.Upgrader` |
| `SrvUpgraderCompressOption(true)` | 启用服务端压缩 |
| `SrvCheckOriginOption(fn)` | 自定义跨域检查函数 |
| `SrvPullChannelsOption([]int{ch1, ch2})` | 注册拉取通知通道 |

```go
ws.Accept(ctx, w, r, meta,
    ws.SrvUpgraderCompressOption(true),
    ws.SrvCheckOriginOption(func(r *http.Request) bool {
        return r.Header.Get("Origin") == "https://example.com"
    }),
    ws.SrvPullChannelsOption([]int{1, 2, 3}),
)
```

### 客户端专用

| 函数 | 说明 |
|---|---|
| `ClientIdOption(id)` | 自定义客户端连接 ID（默认时间戳） |
| `ClientDialOption(dialer)` | 自定义 `websocket.Dialer` |
| `ClientDialWssOption(url, secure)` | WSS 连接配置（secure=false 跳过证书验证） |
| `ClientDialCompressOption(true)` | 启用客户端压缩 |
| `ClientDialHandshakeTimeoutOption(d)` | 握手超时（默认 10s） |
| `ClientDialRetryOption(num, interval)` | 拨号重试次数和间隔（默认 3次, 1s） |
| `ClientDialConnFailedHandlerOption(h)` | 拨号失败回调 |

```go
ws.DialConnect(ctx, "wss://example.com/ws", nil,
    ws.ClientIdOption("client-001"),
    ws.ClientDialWssOption("wss://example.com/ws", false), // 跳过证书验证
    ws.ClientDialRetryOption(5, 2*time.Second),
    ws.ClientDialConnFailedHandlerOption(func(ctx context.Context, conn ws.IConnection) {
        log.Error(ctx, "连接最终失败")
    }),
)
```

### Hub 选项

| 函数 | 说明 |
|---|---|
| `HubShardOption(cnt)` | Hub 分片数量（用于高并发场景） |

```go
ws.InitServerWithOpt(ws.ServerOption{
    HubOpts: []ws.HubOption{
        ws.HubShardOption(8), // 8 个分片
    },
})
```

---

## 消息发送

### SendMsg

发送普通消息。

```go
func (c *Connection) SendMsg(ctx context.Context, payload IMessage, sc SendCallback) error
```

```go
msg := ws.NewMessage(1001)
msg.SetData([]byte("hello"))
err := conn.SendMsg(ctx, msg, func(ctx context.Context, c ws.IConnection, err error) {
    if err != nil {
        log.Error(ctx, "发送失败: %v", err)
    }
})
```

### SendRequestMsg (RPC 调用)

发送请求消息并等待响应。

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
    log.Error(ctx, "RPC 失败: %v", err) // 可能是 ErrWsRpcResponseTimeout
    return
}
log.Info(ctx, "收到响应: %v", string(resp.GetData()))
```

### SendResponseMsg (RPC 响应)

回复 RPC 请求，需传入请求的 sn。

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

## 连接控制

### KickClient

服务端踢出客户端连接。

```go
conn.KickClient(true)  // 顶号模式：会发送 P_DISPLACE 消息
conn.KickClient(false) // 直接断开
```

### KickServer

客户端主动断开连接。

```go
conn.KickServer()
```

### DisplaceClientByIp

按 IP 踢出客户端（集群场景下使用）。

```go
conn.DisplaceClientByIp(ctx, "192.168.1.100")
```

### IsStopped / IsDisplaced

```go
if conn.IsStopped() {
    log.Info(ctx, "连接已停止")
}
if conn.IsDisplaced() {
    log.Info(ctx, "连接被顶号")
}
```

### RefreshDeadline

刷新连接读写超时时间。

```go
conn.RefreshDeadline()
```

---

## 消息拉取 (Puller)

适用于服务端需要主动向客户端推送数据的场景。

### 配置拉取通道

```go
ws.Accept(ctx, w, r, meta,
    ws.SrvPullChannelsOption([]int{1, 2}), // 注册通道 1 和 2
)
```

### 创建 Puller

```go
func NewDefaultPuller(conn IConnection, pullChannelId int,
    firstPullFunc, pullFunc func(context.Context, IConnection)) Puller
```

- `firstPullFunc`：首次连接建立后执行一次
- `pullFunc`：每次收到拉取信号时执行

```go
puller := ws.NewDefaultPuller(conn, 1,
    func(ctx context.Context, c ws.IConnection) {
        // 首次：推送历史数据
        pushHistory(ctx, c)
    },
    func(ctx context.Context, c ws.IConnection) {
        // 每次：推送增量数据
        pushIncrement(ctx, c)
    },
)
puller.PullSend() // 阻塞执行
```

### 触发拉取

```go
// 从任意位置触发消息推送
conn.SignalPullSend(ctx, 1)
```

---

## 公共数据存储

连接对象提供线程安全的 key-value 存储，用于关联业务状态。

```go
conn.SetCommDataValue("loginTime", time.Now())
conn.SetCommDataValue("score", 100)

val, ok := conn.GetCommDataValue("loginTime")
if ok {
    log.Info(ctx, "登录时间: %v", val)
}

conn.IncrCommDataValueBy("score", 10) // score = 110
conn.RemoveCommDataValue("loginTime")
```

---

## 错误常量

```go
var ErrWsRpcResponseTimeout = errors.New("rpc cancel or timeout")
var ErrWsRpcWaitChanClosed  = errors.New("sn channel is closed")
```

---

## 完整示例

### 服务端完整示例

```go
package main

import (
    "context"
    "net/http"

    "github.com/liumingmin/goutils/log"
    "github.com/liumingmin/goutils/ws"
)

func main() {
    // 注册消息处理器
    ws.RegisterHandler(1001, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
        log.Info(ctx, "收到来自 %v 的消息: %v", conn.UserId(), string(msg.GetData()))

        // 回复
        resp := ws.NewMessage(1002)
        resp.SetData([]byte("pong"))
        return conn.SendMsg(ctx, resp, nil)
    })

    // 初始化服务端（带分片）
    ws.InitServerWithOpt(ws.ServerOption{
        HubOpts: []ws.HubOption{ws.HubShardOption(4)},
    })

    // HTTP 路由
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
                log.Info(ctx, "用户 %v 上线", c.UserId())
            }),
            ws.ConnClosedHandlerOption(func(ctx context.Context, c ws.IConnection) {
                log.Info(ctx, "用户 %v 下线", c.UserId())
            }),
            ws.DebugOption(true),
        )
        if err != nil {
            return
        }

        // 触发拉取（示例）
        conn.SignalPullSend(r.Context(), 1)
    })

    log.Info(context.Background(), "WebSocket 服务启动 :8080")
    http.ListenAndServe(":8080", nil)
}
```

### 客户端完整示例

```go
package main

import (
    "context"
    "time"

    "github.com/liumingmin/goutils/log"
    "github.com/liumingmin/goutils/ws"
)

func main() {
    // 注册消息处理器
    ws.RegisterHandler(1002, func(ctx context.Context, conn ws.IConnection, msg ws.IMessage) error {
        log.Info(ctx, "收到服务端消息: %v", string(msg.GetData()))
        return nil
    })

    // 初始化客户端
    ws.InitClient()

    ctx := context.Background()

    // 连接服务端
    conn, err := ws.DialConnect(ctx, "ws://localhost:8080/ws?userId=user1&deviceId=d1", nil,
        ws.ClientDialRetryOption(5, 2*time.Second),
    )
    if err != nil {
        log.Error(ctx, "连接失败: %v", err)
        return
    }

    // 发送消息
    for i := 0; i < 5; i++ {
        msg := ws.NewMessage(1001)
        msg.SetData([]byte("hello"))
        conn.SendMsg(ctx, msg, nil)
        time.Sleep(time.Second)
    }

    // RPC 调用
    reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    req := ws.NewMessage(1001)
    req.SetData([]byte("rpc request"))
    resp, err := conn.SendRequestMsg(reqCtx, req, nil)
    cancel()
    if err != nil {
        log.Error(ctx, "RPC 失败: %v", err)
    } else {
        log.Info(ctx, "RPC 响应: %v", string(resp.GetData()))
    }
}
```
