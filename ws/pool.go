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

func getPoolMessage() *Message {
	msg := messagePool.Get().(*Message)
	msg.isPool = true
	return msg
}

func putPoolMessage(msg *Message) {
	if !msg.isPool {
		return
	}

	if msg.dataMsg != nil {
		putPoolDataMsg(msg.pMsg.ProtocolId, msg.dataMsg)
		msg.dataMsg = nil
	}

	msg.pMsg.Reset()
	msg.sc = nil
	messagePool.Put(msg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
