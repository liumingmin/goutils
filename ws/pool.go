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
	pMsg := pbMessagePool.Get().(*P_MESSAGE)
	pMsg.Type = P_MSG_POOL
	return pMsg
}

func PutPMessage(msg *P_MESSAGE) {
	if msg.Type != P_MSG_POOL {
		return
	}

	msg.Reset()
	pbMessagePool.Put(msg)
}
