# ws 模块技术架构文档

## 1. 概述

`ws` 模块是基于 [gorilla/websocket](https://github.com/gorilla/websocket) 封装的高性能 WebSocket 通信框架，支持**服务端**和**客户端**两种角色，提供统一的消息收发、RPC 调用、连接管理、消息拉取等能力。

## 2. 模块结构

```
ws/
├── constant.go      # 常量定义（连接类型、字符集）
├── def.go           # 核心接口与类型定义、初始化函数
├── conn.go          # Connection 实现：读写循环、消息分发、RPC
├── hub.go           # Hub / shardHub 连接注册中心
├── server_conn.go   # 服务端入口：Accept、KickClient、DisplaceClientByIp
├── client_conn.go   # 客户端入口：DialConnect、AutoReDialConnect、KickServer
├── option.go        # 函数式选项（~20个配置项）
├── msg_core.go      # Message 结构与二进制序列化
├── msg.pb.go        # Protobuf 生成代码（P_DISPLACE）
├── pool.go          # 对象池（Message、Connection、DataMessage）
├── puller.go        # Puller 消息拉取机制
└── wss_test.go      # 集成测试
```

## 3. 核心组件架构

```
┌─────────────────────────────────────────────────────────────┐
│                        应用层                                │
│   RegisterHandler(protocolId, handler)                       │
│   RegisterDataMsgType(protocolId, protoMsg)                  │
└───────────────┬─────────────────────────────┬────────────────┘
                │                             │
        ┌───────▼───────┐             ┌───────▼───────┐
        │  InitServer()  │             │  InitClient()  │
        └───────┬───────┘             └───────┬───────┘
                │                             │
        ┌───────▼───────┐             ┌───────▼───────┐
        │ ClientConnHub │             │ ServerConnHub │
        │  (shardHub)   │             │   (Hub)       │
        └───────┬───────┘             └───────┬───────┘
                │                             │
        ┌───────▼───────┐             ┌───────▼───────┐
        │  Accept()     │             │DialConnect()  │
        │  HTTP Upgrade │             │  WS Dial      │
        └───────┬───────┘             └───────┬───────┘
                │                             │
                └──────────┬──────────────────┘
                           │
                   ┌───────▼───────┐
                   │  Connection   │
                   │  ┌──────────┐ │
                   │  │ readLoop │ │  ← goroutine
                   │  │writeLoop │ │  ← goroutine
                   │  │ dispatch │ │
                   │  └──────────┘ │
                   └───────────────┘
```

## 4. 二进制线协议

模块使用自定义二进制协议进行消息传输（非 protobuf 封帧）：

```
┌──────────┬──────────┬──────────────┬──────────────┬─────────────┐
│ MagicFlag│ Length   │ protocolId   │ sn           │ Payload     │
│ 2 bytes  │ 4 bytes  │ 4 bytes LE   │ 4 bytes LE   │ N bytes     │
│ 0xFE 0xEF│ uint32  │ uint32       │ uint32       │ (protobuf)  │
└──────────┴──────────┴──────────────┴──────────────┴─────────────┘
   包头(6B)                      消息头(8B)              载荷

总头部长度: 14 字节
```

- **MagicFlag** (`0xFE, 0xEF`)：用于流同步与包校验
- **Length**：载荷长度 = 消息头(8B) + 数据体长度，小端序
- **protocolId**：消息协议 ID，用于路由到对应 Handler
- **sn**：序列号，用于 RPC 请求/响应匹配。服务端起始值=1（奇数），客户端起始值=0（偶数），每次 +2
- **Payload**：protobuf 序列化数据（可为空）

## 5. 连接生命周期

### 5.1 服务端连接

```
HTTP请求 → Accept() → Upgrader升级为WS → 创建Connection
→ 注册到ClientConnHub → connEstablishHandler回调
→ 启动readLoop + writeLoop goroutine
→ 断开时: Hub注销 → connClosingHandler → 关闭Socket → connClosedHandler
→ Connection回收到对象池
```

### 5.2 客户端连接

```
DialConnect() → Dialer拨号(带重试) → 创建Connection
→ 注册到ServerConnHub → connEstablishHandler回调
→ 启动readLoop + writeLoop goroutine
→ 断开时: 同服务端流程
```

### 5.3 连接置换（顶号）

当相同 ID 的新连接注册时，Hub 自动踢出旧连接：

```
新连接注册 → Hub检测到ID重复 → 旧连接setDisplaced()
→ 从Hub移除旧连接 → connClosingHandler → 发送P_DISPLACE消息给旧连接 → 关闭旧Socket
```

`P_DISPLACE` 消息包含旧IP、新IP、时间戳，协议ID为 `P_BASE_s2c_err_displace (2147483647)`。

## 6. 消息处理流程

### 6.1 发送流程

```
SendMsg() → 写入sendBuffer channel → writeLoop取出
→ marshal()序列化DataMsg → writeHead(6B) + msgHead(8B) + data
→ WebSocket Write → SendCallback回调
```

支持批量发送：writeLoop 在发送一条消息后会尝试 drain sendBuffer 中所有待发消息。

### 6.2 接收流程

```
readLoop → readMessageData()读取二进制帧
→ 校验MagicFlag → 校验Length → 解析msgHead(protocolId, sn)
→ processMsg() → dispatch()
→ 若sn>0: 匹配snChanMap投递RPC响应
→ 查找msgHandlers执行对应Handler
```

### 6.3 RPC 调用

```
发送方: SendRequestMsg() → 分配sn → 注册sn→channel映射 → SendMsg() → 阻塞等待channel
接收方: 收到消息 → Handler处理 → SendResponseMsg()设置reqSn → SendMsg()
发送方: dispatch()匹配sn → channel收到响应 → 返回
```

超时由 context 控制，超时返回 `ErrWsRpcResponseTimeout`。

## 7. Hub 连接管理

### 7.1 单 Hub

使用 `sync.Map` 存储连接，通过 register/unregister channel 串行处理注册注销，避免并发冲突。

### 7.2 分片 Hub (shardHub)

使用 `xxh3.HashString(connId)` 对连接 ID 哈希取模，将连接分散到多个子 Hub 中，降低单个 `sync.Map` 的竞争压力。

```
shardHub
├── Hub[0]  ← connId hash % N == 0
├── Hub[1]  ← connId hash % N == 1
├── ...
└── Hub[N-1]
```

通过 `HubShardOption(cnt)` 配置分片数量。

## 8. 对象池

模块大量使用 `sync.Pool` 减少 GC 压力：

| 池 | 管理对象 | 说明 |
|---|---|---|
| `messagePool` | `*Message` | 消息对象池 |
| `srvConnectionPool` | `*Connection` | 服务端连接对象池 |
| `dataMsgPools[protocolId]` | `IDataMessage` | 按协议ID分池的数据消息对象 |

池化消息通过 `isPool` 标志控制回收，防止重复归还。

## 9. 消息拉取机制 (Puller)

适用于服务端需要主动向客户端推送数据的场景：

```
SignalPullSend(ctx, channelId) → 向pullChannel写入信号
→ Puller.PullSend()收到信号 → 执行pullFunc回调 → 发送消息给客户端
→ 循环阻塞等待下一个信号
```

- `firstPullFunc`：首次连接时执行一次
- `pullFunc`：每次收到信号时执行
- 使用 `atomic.CompareAndSwap` 防止并发执行同一 Puller

## 10. 错误处理策略

| 场景 | 处理方式 |
|---|---|
| 网络超时 (`net.Error.Timeout()`) | 重试最多 `maxFailureRetry` 次，间隔 `temporaryWait` |
| WebSocket 正常关闭 | Debug 级别日志，静默退出 |
| WebSocket 异常关闭 | Warn 级别日志 |
| 消息过大 (`> maxMessageBytesSize`) | 返回错误，丢弃消息 |
| 包头 MagicFlag 不匹配 | 返回 "packet head flag error" |
| RPC 超时 | 返回 `ErrWsRpcResponseTimeout` |
| goroutine panic | `log.Recover()` 捕获并记录 |

## 11. 并发模型

- 每个 Connection 启动 2 个 goroutine：`readLoop` 和 `writeLoop`
- 读写通过 `sendBuffer` channel 解耦
- 连接注册/注销通过 Hub 的 channel 串行化
- `stopped` 和 `displaced` 使用 `atomic` 操作保证并发安全
- `commonData` 使用 `RWMutex` 保护
- `snChanMap` 使用 `sync.Map` 存储 RPC 等待通道

## 12. 依赖

| 依赖 | 用途 |
|---|---|
| `github.com/gorilla/websocket` | WebSocket 协议实现 |
| `github.com/zeebo/xxh3` | Hub 分片哈希算法 |
| `google.golang.org/protobuf` | Protobuf 序列化 |
| `github.com/liumingmin/goutils/log` | 日志（带 traceId） |
| `github.com/liumingmin/goutils/utils` | 工具函数（SafeCloseChan） |
| `github.com/liumingmin/goutils/utils/safego` | 安全 goroutine 启动 |
| `github.com/liumingmin/goutils/net/ip` | HTTP 远程 IP 提取 |

## 13. 关键设计决策

1. **自定义二进制协议而非纯 protobuf 封帧**：MagicFlag 提供流同步能力，避免粘包问题
2. **Hub channel 串行化**：注册/注销通过 channel 避免 `sync.Map` 的细粒度锁竞争
3. **服务端/客户端 SN 分奇偶**：服务端从 1 开始(奇数)，客户端从 0 开始(偶数)，天然避免冲突
4. **对象池化 Connection**：服务端连接默认池化，通过 `reset()` 复用，减少 GC
5. **函数式选项模式**：所有配置通过 `ConnOption` 函数组合，灵活且向后兼容
