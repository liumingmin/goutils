package ws

import "github.com/gorilla/websocket"

//通用option
func ConnectCbOption(connCallback IConnCallback) ConnOption {
	return func(conn *Connection) {
		conn.connCallback = connCallback
	}
}

func HeartbeatCbOption(heartbeatCallback IHeartbeatCallback) ConnOption {
	return func(conn *Connection) {
		conn.heartbeatCallback = heartbeatCallback
	}
}

func SendBufferOption(bufferSize int) ConnOption {
	return func(conn *Connection) {
		conn.sendBuffer = make(chan interface{}, bufferSize)
	}
}

//服务端特有
//upgrader定制
func SrvUpgraderOption(upgrader *websocket.Upgrader) ConnOption {
	return func(conn *Connection) {
		conn.upgrader = upgrader
	}
}

//为每种消息拉取逻辑分别注册不同的通道
func SrvPullChannelsOption(channels []int) ConnOption {
	return func(conn *Connection) {
		pullChannelMap := make(map[int]chan struct{})
		for _, channel := range channels {
			pullChannelMap[channel] = make(chan struct{}, 2)
		}

		conn.pullChannelMap = pullChannelMap
	}
}
