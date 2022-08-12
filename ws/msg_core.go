package ws

import "google.golang.org/protobuf/proto"

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
