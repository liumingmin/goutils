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

	dataMsgPools = make(map[int32]*sync.Pool)
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

func getPoolDataMsg(protocolId int32) IDataMessage {
	pool, ok := dataMsgPools[protocolId]
	if !ok {
		return nil
	}
	return pool.Get().(IDataMessage)
}

func putPoolDataMsg(protocolId int32, dataMsg IDataMessage) {
	pool, ok := dataMsgPools[protocolId]
	if !ok {
		return
	}

	dataMsg.Reset()
	pool.Put(dataMsg)
}
