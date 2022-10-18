package ws

import (
	"reflect"
	"sync"
)

var (
	messagePool = sync.Pool{
		New: func() interface{} {
			msg := &Message{}
			msg.isPool = true
			return msg
		},
	}

	dataMsgPools = make(map[uint32]*sync.Pool)
	dataMsgTypes = make(map[uint32]reflect.Type)

	srvConnectionPool = sync.Pool{
		New: func() interface{} {
			return &Connection{}
		},
	}
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
		putPoolDataMsg(msg.protocolId, msg.dataMsg)
		msg.dataMsg = nil
	}

	msg.protocolId = 0
	msg.data = nil
	msg.sc = nil
	messagePool.Put(msg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func getPoolDataMsg(protocolId uint32) IDataMessage {
	pool, ok := dataMsgPools[protocolId]
	if !ok {
		return nil
	}
	return pool.Get().(IDataMessage)
}

func getDataMsg(protocolId uint32) IDataMessage {
	typ, ok := dataMsgTypes[protocolId]
	if !ok {
		return nil
	}
	return reflect.New(typ).Interface().(IDataMessage)
}

func putPoolDataMsg(protocolId uint32, dataMsg IDataMessage) {
	pool, ok := dataMsgPools[protocolId]
	if !ok {
		return
	}

	dataMsg.Reset()
	pool.Put(dataMsg)
}

func getPoolConnection() *Connection {
	conn := srvConnectionPool.Get().(*Connection)
	conn.isPool = true
	return conn
}

func putPoolConnection(conn *Connection) {
	if !conn.isPool {
		return
	}

	conn.Reset()
	srvConnectionPool.Put(conn)
}
