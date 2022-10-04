package ws

import (
	"context"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

var ClientConnHub IHub //服务端管理的来自客户端的连接
var ServerConnHub IHub //客户端管理的连向服务端的连接

//server invoke 服务端调用
func InitServer() {
	InitServerWithOpt(ServerOption{})
}

func InitServerWithOpt(serverOpt ServerOption) {
	ClientConnHub = initServer(serverOpt)
}

//client invoke 客户端调用
func InitClient() {
	ServerConnHub = initClient()
}

type IMessage interface {
	PMsg() *P_MESSAGE
	DataMsg() IDataMessage
}

// P_MESSAGE.Data类型的接口
type IDataMessage interface {
	proto.Message
	Reset()
}

type IConnection interface {
	Id() string
	ConnType() ConnType
	UserId() string
	Type() int
	DeviceId() string
	Version() int
	Charset() int
	ClientIp() string
	Reset()
	IsStopped() bool
	IsDisplaced() bool
	RefreshDeadline()
	SendMsg(ctx context.Context, payload IMessage, sc SendCallback) error

	KickClient(displace bool) //server side invoke
	KickServer()              //client side invoke

	GetPullChannel(pullChannelId int) (chan struct{}, bool)
	SendPullNotify(ctx context.Context, pullChannelId int) error

	GetCommDataValue(key string) (interface{}, bool)
	SetCommDataValue(key string, value interface{})
	RemoveCommDataValue(key string)
	IncrCommDataValueBy(key string, delta int)
}

type IHub interface {
	Find(string) (*Connection, error)
	RangeConnsByFunc(func(string, *Connection) bool)
	ConnectionIds() []string

	registerConn(*Connection)
	unregisterConn(*Connection)
	run()
}

//normal message 普通消息
func NewMessage() IMessage {
	return &Message{
		pMsg: &P_MESSAGE{},
	}
}

//pool message 对象池消息
func GetPoolMessage(protocolId int32) IMessage {
	msg := getPoolMessage()
	msg.pMsg.ProtocolId = protocolId
	msg.dataMsg = getPoolDataMsg(protocolId)
	return msg
}

// 消息发送回调接口
type SendCallback func(ctx context.Context, c IConnection, err error)

// 客户端消息处理函数对象  use RegisterHandler(protocolId, MsgHandler)
type MsgHandler func(context.Context, IConnection, IMessage) error

// 客户端事件处理函数
// ConnEstablishHandlerOption  sync(阻塞主流程)
// ConnClosingHandlerOption   sync(阻塞主流程)
// ConnClosedHandlerOption  async
type EventHandler func(context.Context, IConnection)

// 注册消息处理器
func RegisterHandler(protocolId int32, h MsgHandler) {
	msgHandlers.Store(protocolId, h)
}

// 注册数据消息类型[P_MESSAGE.Data],功能可选，当需要使用框架提供的池功能时使用
func RegisterDataMsgType(protocolId int32, pMsg IDataMessage) {
	typ := reflect.TypeOf(pMsg)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	dataMsgPools[protocolId] = &sync.Pool{
		New: func() interface{} {
			return reflect.New(typ).Interface()
		},
	}
	dataMsgTypes[protocolId] = typ
}
