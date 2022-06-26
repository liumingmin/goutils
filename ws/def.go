package ws

import (
	"context"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Message struct {
	pMsg   *P_MESSAGE   // pb消息体
	isPool bool         // 是否对象池消息
	sc     SendCallback // 消息发送回调接口
}

func (t *Message) PMsg() *P_MESSAGE {
	return t.pMsg
}

func (t *Message) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, t.pMsg)
}

func NewMessage() *Message {
	return &Message{
		pMsg: &P_MESSAGE{},
	}
}

// 消息发送回调接口
type SendCallback func(ctx context.Context, c *Connection, err error)

// 客户端消息处理函数对象
// use RegisterHandler(constant...., func(context.Context,*Connection,*Message) error {})
type Handler func(context.Context, *Connection, *Message) error

//连接回调
type IConnCallback interface {
	ConnFinished(clientId string)
	DisconnFinished(clientId string)
}

//保活回调
type IHeartbeatCallback interface {
	RecvPing(clientId string)
	RecvPong(clientId string) error
}

// 注册消息处理器
func RegisterHandler(cmd int32, h Handler) {
	Handlers[cmd] = h
}

// P_MESSAGE.Data类型的接口
type IDataMessage interface {
	proto.Message
	Reset()
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
