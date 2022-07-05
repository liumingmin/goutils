package ws

import (
	"context"
	"crypto/tls"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

// 连接动态参数选项
type ConnOption func(*Connection)

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
		conn.sendBuffer = make(chan *Message, bufferSize)
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

func SrvUpgraderCompressOption(compress bool) ConnOption {
	return func(conn *Connection) {
		conn.upgrader.EnableCompression = compress
	}
}

// 客户端专用，Dialer动态参数选项
type DialerOption func(*websocket.Dialer)

func DialerWssOption(sUrl string, secureWss bool) DialerOption {
	u, err := url.Parse(sUrl)
	if err != nil {
		log.Error(context.Background(), "Parse url %s err:%v", sUrl, err)
	}

	return func(dialer *websocket.Dialer) {
		if u != nil && u.Scheme == "wss" && !secureWss {
			dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
}

func DialerCompressOption(compress bool) DialerOption { //, compressLevel int
	return func(dialer *websocket.Dialer) {
		dialer.EnableCompression = compress
	}
}

func DialerHandshakeTimeoutOption(handshakeTimeout time.Duration) DialerOption {
	return func(dialer *websocket.Dialer) {
		dialer.HandshakeTimeout = handshakeTimeout
	}
}
