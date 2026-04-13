# ws Module Technical Architecture

## 1. Overview

The `ws` module is a high-performance WebSocket communication framework built on [gorilla/websocket](https://github.com/gorilla/websocket). It supports both **server** and **client** roles, providing unified capabilities for message sending/receiving, RPC calls, connection management, and message pulling.

## 2. Module Structure

```
ws/
├── constant.go      # Constants (connection types, charsets)
├── def.go           # Core interfaces, type definitions, initialization functions
├── conn.go          # Connection implementation: read/write loops, message dispatch, RPC
├── hub.go           # Hub / shardHub connection registry
├── server_conn.go   # Server entry: Accept, KickClient, DisplaceClientByIp
├── client_conn.go   # Client entry: DialConnect, AutoReDialConnect, KickServer
├── option.go        # Functional options (~20 configuration items)
├── msg_core.go      # Message struct and binary serialization
├── msg.pb.go        # Protobuf generated code (P_DISPLACE)
├── pool.go          # Object pools (Message, Connection, DataMessage)
├── puller.go        # Puller message pulling mechanism
└── wss_test.go      # Integration tests
```

## 3. Core Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       Application Layer                      │
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

## 4. Binary Wire Protocol

The module uses a custom binary protocol for message transmission (not protobuf framing):

```
┌──────────┬──────────┬──────────────┬──────────────┬─────────────┐
│ MagicFlag│ Length   │ protocolId   │ sn           │ Payload     │
│ 2 bytes  │ 4 bytes  │ 4 bytes LE   │ 4 bytes LE   │ N bytes     │
│ 0xFE 0xEF│ uint32  │ uint32       │ uint32       │ (protobuf)  │
└──────────┴──────────┴──────────────┴──────────────┴─────────────┘
  Header(6B)                      Msg Head(8B)              Body

Total header length: 14 bytes
```

- **MagicFlag** (`0xFE, 0xEF`): Used for stream synchronization and packet validation
- **Length**: Payload length = message head (8B) + data body length, little-endian
- **protocolId**: Message protocol ID, used for routing to the corresponding handler
- **sn**: Sequence number, used for RPC request/response matching. Server starts at 1 (odd), client starts at 0 (even), increments by 2
- **Payload**: Protobuf serialized data (can be empty)

## 5. Connection Lifecycle

### 5.1 Server Connection

```
HTTP Request → Accept() → Upgrader upgrades to WS → Create Connection
→ Register to ClientConnHub → connEstablishHandler callback
→ Start readLoop + writeLoop goroutines
→ On disconnect: Hub unregister → connClosingHandler → Close Socket → connClosedHandler
→ Connection returned to object pool
```

### 5.2 Client Connection

```
DialConnect() → Dialer with retry → Create Connection
→ Register to ServerConnHub → connEstablishHandler callback
→ Start readLoop + writeLoop goroutines
→ On disconnect: Same as server flow
```

### 5.3 Connection Displacement

When a new connection with the same ID registers, the Hub automatically kicks the old connection:

```
New connection registers → Hub detects duplicate ID → Old connection setDisplaced()
→ Remove old connection from Hub → connClosingHandler → Send P_DISPLACE message to old connection → Close old Socket
```

The `P_DISPLACE` message contains old IP, new IP, and timestamp, with protocol ID `P_BASE_s2c_err_displace (2147483647)`.

## 6. Message Processing Flow

### 6.1 Send Flow

```
SendMsg() → Write to sendBuffer channel → writeLoop retrieves
→ marshal() serializes DataMsg → writeHead(6B) + msgHead(8B) + data
→ WebSocket Write → SendCallback callback
```

Supports batch sending: after sending one message, writeLoop attempts to drain all pending messages from sendBuffer.

### 6.2 Receive Flow

```
readLoop → readMessageData() reads binary frame
→ Validate MagicFlag → Validate Length → Parse msgHead(protocolId, sn)
→ processMsg() → dispatch()
→ If sn>0: Match snChanMap to deliver RPC response
→ Look up msgHandlers to execute corresponding handler
```

### 6.3 RPC Call

```
Sender: SendRequestMsg() → Allocate sn → Register sn→channel mapping → SendMsg() → Block waiting on channel
Receiver: Receive message → Handler processes → SendResponseMsg() sets reqSn → SendMsg()
Sender: dispatch() matches sn → Channel receives response → Return
```

Timeout is controlled by context; timeout returns `ErrWsRpcResponseTimeout`.

## 7. Hub Connection Management

### 7.1 Single Hub

Uses `sync.Map` to store connections, serializes registration/unregistration through register/unregister channels to avoid concurrency conflicts.

### 7.2 Sharded Hub (shardHub)

Uses `xxh3.HashString(connId)` to hash connection IDs modulo the number of shards, distributing connections across multiple sub-Hubs to reduce contention on a single `sync.Map`.

```
shardHub
├── Hub[0]  ← connId hash % N == 0
├── Hub[1]  ← connId hash % N == 1
├── ...
└── Hub[N-1]
```

Configure shard count via `HubShardOption(cnt)`.

## 8. Object Pools

The module extensively uses `sync.Pool` to reduce GC pressure:

| Pool | Managed Object | Description |
|---|---|---|
| `messagePool` | `*Message` | Message object pool |
| `srvConnectionPool` | `*Connection` | Server connection object pool |
| `dataMsgPools[protocolId]` | `IDataMessage` | Per-protocol-ID data message object pool |

Pooled messages use the `isPool` flag to control recycling, preventing double returns.

## 9. Message Pulling Mechanism (Puller)

Suitable for scenarios where the server needs to actively push data to clients:

```
SignalPullSend(ctx, channelId) → Write signal to pullChannel
→ Puller.PullSend() receives signal → Execute pullFunc callback → Send message to client
→ Loop blocking for next signal
```

- `firstPullFunc`: Executed once on first connection
- `pullFunc`: Executed each time a signal is received
- Uses `atomic.CompareAndSwap` to prevent concurrent execution of the same Puller

## 10. Error Handling Strategy

| Scenario | Handling |
|---|---|
| Network timeout (`net.Error.Timeout()`) | Retry up to `maxFailureRetry` times, with `temporaryWait` interval |
| WebSocket normal close | Debug level log, silent exit |
| WebSocket abnormal close | Warn level log |
| Message too large (`> maxMessageBytesSize`) | Return error, discard message |
| Packet header MagicFlag mismatch | Return "packet head flag error" |
| RPC timeout | Return `ErrWsRpcResponseTimeout` |
| goroutine panic | `log.Recover()` catches and logs |

## 11. Concurrency Model

- Each Connection spawns 2 goroutines: `readLoop` and `writeLoop`
- Read and write are decoupled via `sendBuffer` channel
- Connection registration/unregistration is serialized through Hub channels
- `stopped` and `displaced` use `atomic` operations for concurrency safety
- `commonData` is protected by `RWMutex`
- `snChanMap` uses `sync.Map` to store RPC waiting channels

## 12. Dependencies

| Dependency | Purpose |
|---|---|
| `github.com/gorilla/websocket` | WebSocket protocol implementation |
| `github.com/zeebo/xxh3` | Hub sharding hash algorithm |
| `google.golang.org/protobuf` | Protobuf serialization |
| `github.com/liumingmin/goutils/log` | Logging (with traceId) |
| `github.com/liumingmin/goutils/utils` | Utility functions (SafeCloseChan) |
| `github.com/liumingmin/goutils/utils/safego` | Safe goroutine launching |
| `github.com/liumingmin/goutils/net/ip` | HTTP remote IP extraction |

## 13. Key Design Decisions

1. **Custom binary protocol instead of pure protobuf framing**: MagicFlag provides stream synchronization, avoiding packet sticking issues
2. **Hub channel serialization**: Registration/unregistration via channels avoids fine-grained lock contention on `sync.Map`
3. **Server/client SN parity separation**: Server starts at 1 (odd), client starts at 0 (even), naturally avoiding conflicts
4. **Connection object pooling**: Server connections are pooled by default, reused via `reset()`, reducing GC
5. **Functional options pattern**: All configuration composed via `ConnOption` functions, flexible and backward compatible
