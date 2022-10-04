package ws

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

// 连接动态参数选项
type ConnOption func(*Connection)
type HubOption func(IHub)

type ServerOption struct {
	HubOpts []HubOption
}

func DebugOption(debug bool) ConnOption {
	return func(conn *Connection) {
		conn.debug = debug
	}
}

// callback
func ConnEstablishHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.connEstablishHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v connEstablishHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func ConnClosingHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.connClosingHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v connClosingHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func ConnClosedHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.connClosedHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v connClosedHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func RecvPingHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.recvPingHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v recvPingHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func RecvPongHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.recvPongHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v recvPongHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func SendBufferOption(bufferSize int) ConnOption {
	return func(conn *Connection) {
		conn.sendBuffer = make(chan *Message, bufferSize)
	}
}

func CompressionLevelOption(compressionLevel int) ConnOption {
	return func(conn *Connection) {
		if compressionLevel <= 0 {
			return
		}
		conn.compressionLevel = compressionLevel
	}
}

func defaultNetParamsOption() ConnOption {
	return func(conn *Connection) {
		conn.maxFailureRetry = 10                   //重试次数
		conn.readWait = 60 * time.Second            //读等待
		conn.writeWait = 60 * time.Second           //写等待
		conn.temporaryWait = 500 * time.Millisecond //网络抖动重试等待
	}
}

func NetMaxFailureRetryOption(maxFailureRetry int) ConnOption {
	return func(conn *Connection) {
		if maxFailureRetry < 0 {
			return
		}

		conn.maxFailureRetry = maxFailureRetry
	}
}

func NetReadWaitOption(readWait time.Duration) ConnOption {
	return func(conn *Connection) {
		if readWait <= 0 {
			return
		}

		conn.readWait = readWait
	}
}

func NetWriteWaitOption(writeWait time.Duration) ConnOption {
	return func(conn *Connection) {
		if writeWait <= 0 {
			return
		}

		conn.writeWait = writeWait
	}
}

func NetTemporaryWaitOption(temporaryWait time.Duration) ConnOption {
	return func(conn *Connection) {
		if temporaryWait <= 0 {
			return
		}

		conn.temporaryWait = temporaryWait
	}
}

func MaxMessageBytesSizeOption(size uint32) ConnOption {
	return func(con *Connection) {
		con.maxMessageBytesSize = size
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
func SrvPullChannelsOption(pullChannelIds []int) ConnOption {
	return func(conn *Connection) {
		conn.pullChannelIds = pullChannelIds
	}
}

func SrvUpgraderCompressOption(compress bool) ConnOption {
	return func(conn *Connection) {
		conn.upgrader.EnableCompression = compress
	}
}

func SrvCheckOriginOption(checkOrigin func(r *http.Request) bool) ConnOption {
	return func(conn *Connection) {
		conn.upgrader.CheckOrigin = checkOrigin
	}
}

// 客户端专用
// 默认使用时间戳来记录客户端所连服务器的id
func ClientIdOption(id string) ConnOption {
	return func(conn *Connection) {
		conn.id = id
	}
}

func ClientDialOption(dialer *websocket.Dialer) ConnOption {
	return func(conn *Connection) {
		conn.dialer = dialer
	}
}

func ClientDialWssOption(sUrl string, secureWss bool) ConnOption {
	u, err := url.Parse(sUrl)
	if err != nil {
		log.Error(context.Background(), "Parse url %s err:%v", sUrl, err)
	}

	return func(conn *Connection) {
		if u != nil && u.Scheme == "wss" && !secureWss {
			conn.dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
}

func ClientDialCompressOption(compress bool) ConnOption { //, compressLevel int
	return func(conn *Connection) {
		conn.dialer.EnableCompression = compress
	}
}

func ClientDialHandshakeTimeoutOption(handshakeTimeout time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.dialer.HandshakeTimeout = handshakeTimeout
	}
}

func ClientDialRetryOption(retryNum int, retryInterval time.Duration) ConnOption {
	return func(conn *Connection) {
		conn.dialRetryNum = retryNum
		conn.dialRetryInterval = retryInterval
	}
}

func ClientDialConnFailedHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.dialConnFailedHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v dialConnFailedHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func ClientAutoReconHandlerOption(handler EventHandler) ConnOption {
	return func(conn *Connection) {
		if handler != nil {
			conn.connAutoReconHandler = func(ctx context.Context, c IConnection) {
				defer log.Recover(ctx, func(e interface{}) string {
					return fmt.Sprintf("%v connAutoReconHandler panic, error is: %v", c.ConnType(), e)
				})

				handler(ctx, c)
			}
		}
	}
}

func HubShardOption(cnt uint16) HubOption {
	return func(hub IHub) {
		sHub, ok := hub.(*shardHub)
		if !ok {
			return
		}

		for i := uint16(0); i < cnt; i++ {
			sHub.hubs = append(sHub.hubs, newHub())
		}
	}
}
