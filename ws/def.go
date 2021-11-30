package ws

import "context"

// 消息发送回调接口
type SendCallback func(ctx context.Context, c *Connection, err error)

// 客户端消息处理函数对象
type Handler func(context.Context, *Connection, *P_MESSAGE) error

// 连接动态参数选项
type ConnOption func(*Connection)

type ConnType int8

func (t ConnType) String() string {
	if t == CONN_TYPE_CLIENT {
		return "client"
	}
	if t == CONN_TYPE_SERVER {
		return "server"
	}
	return ""
}

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

type msgSendWrapper struct {
	pbMessage *P_MESSAGE   // 消息体
	sc        SendCallback // 消息发送回调接口
}
