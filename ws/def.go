package ws

import "context"

// 消息发送回调接口
type SendCallback func(ctx context.Context, c *Connection, err error)

// 客户端消息处理函数对象
// use RegisterHandler(constant...., func(context.Context,*Connection,*P_MESSAGE) error {})
type Handler func(context.Context, *Connection, *P_MESSAGE) error

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
