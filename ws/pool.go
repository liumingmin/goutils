package ws

import (
	"sync"
)

var (
	messagePool = sync.Pool{
		New: func() interface{} {
			msg := NewMessage()
			msg.isPool = true
			return msg
		},
	}
)

func GetPoolMessage() *Message {
	msg := messagePool.Get().(*Message)
	msg.isPool = true
	return msg
}

func PutPoolMessage(msg *Message) {
	if !msg.isPool {
		return
	}

	msg.pMsg.Reset()
	msg.sc = nil
	messagePool.Put(msg)
}
