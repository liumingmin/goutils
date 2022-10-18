package ws

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

//不能手动创建，必须使用 NewMessage() 或 GetPoolMessage()
type Message struct {
	protocolId uint32       // 消息协议ID
	data       []byte       // 内容-自定义消息
	dataMsg    IDataMessage // 当为nil时,由用户自定义t.data,当不为nil时,则是池对象 t.data => t.dataMsg
	isPool     bool         // Message是否对象池消息
	sc         SendCallback // 消息发送回调接口
}

func (t *Message) GetProtocolId() uint32 {
	return t.protocolId
}

func (t *Message) GetData() []byte {
	return t.data
}

func (t *Message) SetData(data []byte) {
	t.data = data
}

func (t *Message) DataMsg() IDataMessage {
	return t.dataMsg
}

func (t *Message) protoIdToLEBytes() [4]byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], t.protocolId)
	return bytes
}

func (t *Message) marshal() error {
	var err error
	if len(t.data) == 0 && t.dataMsg != nil {
		t.data, err = proto.Marshal(t.dataMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Message) unmarshal(payload []byte) error {
	t.protocolId = binary.LittleEndian.Uint32(payload[:4])
	t.data = payload[4:]

	if len(t.data) == 0 {
		return nil
	}

	if t.isPool {
		t.dataMsg = getPoolDataMsg(t.protocolId)
	} else {
		t.dataMsg = getDataMsg(t.protocolId)
	}
	if t.dataMsg == nil {
		return nil
	}

	return proto.Unmarshal(t.data, t.dataMsg)
}
