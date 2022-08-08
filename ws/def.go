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
	InitServerWithOpt(&ServerOption{})
}

func InitServerWithOpt(serverOpt *ServerOption) {
	initServer(serverOpt)
}

//client invoke 客户端调用
func InitClient() {
	initClient()
}

type ServerOption struct {
	HubOpts []HubOption
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
	KickServer(displace bool) //client side invoke

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

//不能手动创建，必须使用 NewMessage() 或 GetPoolMessage()
type Message struct {
	pMsg    *P_MESSAGE   // 主消息体,一定不为nil
	dataMsg IDataMessage // 当为nil时,由用户自定义pMsg.Data,当不为nil时,则是池对象 t.pMsg.Data => t.dataMsg
	isPool  bool         // Message是否对象池消息
	sc      SendCallback // 消息发送回调接口
}

func (t *Message) PMsg() *P_MESSAGE {
	return t.pMsg
}

func (t *Message) DataMsg() IDataMessage {
	return t.dataMsg
}

func (t *Message) Marshal() ([]byte, error) {
	var err error
	if len(t.pMsg.Data) == 0 && t.dataMsg != nil {
		t.pMsg.Data, err = proto.Marshal(t.dataMsg)
		if err != nil {
			return nil, err
		}
	}

	return proto.Marshal(t.pMsg)
}

func (t *Message) Unmarshal(payload []byte) error {
	err := proto.Unmarshal(payload, t.pMsg)
	if err != nil {
		return err
	}

	if len(t.pMsg.Data) == 0 {
		return nil
	}

	t.dataMsg = getPoolDataMsg(t.pMsg.ProtocolId)
	if t.dataMsg == nil {
		return nil
	}

	return proto.Unmarshal(t.pMsg.Data, t.dataMsg)
}

//普通消息
func NewMessage() *Message {
	return &Message{
		pMsg: &P_MESSAGE{},
	}
}

//对象池消息
func GetPoolMessage(protocolId int32) *Message {
	msg := getPoolMessage()
	msg.pMsg.ProtocolId = protocolId
	msg.dataMsg = getPoolDataMsg(protocolId)
	return msg
}

// 消息发送回调接口
type SendCallback func(ctx context.Context, c *Connection, err error)

// 客户端消息处理函数对象
// use RegisterHandler(constant...., func(context.Context,*Connection,*Message) error {})
type MsgHandler func(context.Context, *Connection, *Message) error

// 客户端事件处理函数
// ConnEstablishHandlerOption  sync(阻塞主流程)
// ConnClosingHandlerOption   sync(阻塞主流程)
// ConnClosedHandlerOption  async
type EventHandler func(context.Context, *Connection)

// 注册消息处理器
func RegisterHandler(cmd int32, h MsgHandler) {
	Handlers[cmd] = h
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
}

type SrvPuller interface {
	PullSend()
}

func NewDefaultSrvPuller(conn *Connection, pullChannelId int, firstPullFunc, pullFunc func(context.Context, *Connection)) SrvPuller {
	return &defaultSrvPuller{
		conn:          conn,
		pullChannelId: pullChannelId,
		firstPullFunc: firstPullFunc,
		pullFunc:      pullFunc,
	}
}
