package ws

import (
	"context"

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
// use RegisterHandler(constant...., func(context.Context,*Connection,*P_MESSAGE) error {})
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
