package ws

import "sync"

var (
	pbMessagePool = sync.Pool{
		New: func() interface{} {
			return &P_MESSAGE{}
		},
	}
)

func GetPMessage() *P_MESSAGE {
	return pbMessagePool.Get().(*P_MESSAGE)
}

func PutPMessage(msg *P_MESSAGE) {
	msg.Reset()
	pbMessagePool.Put(msg)
}

func PutPMessageIntfs(message interface{}) {
	msg, ok := message.(*P_MESSAGE) //优先判断
	if ok {
		PutPMessage(msg)
	} else {
		w := message.(*msgSendWrapper)
		PutPMessage(w.pbMessage)
	}
}
