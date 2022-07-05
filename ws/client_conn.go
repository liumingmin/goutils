package ws

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

//displace=true，通常在集群环境下，踢掉在其他集群节点建立的连接，当前节点不需要主动调用
func (c *Connection) KickServer(displace bool) {
	if displace {
		c.setDisplaced()
	}

	ServerConnHub.unregister <- c
}

func Connect(ctx context.Context, sId, sUrl string, secureWss bool, header http.Header, opts ...ConnOption) (*Connection, error) {
	baseOpts := []ConnOption{
		ClientIdOption(sId),
		ClientDialWssOption(sUrl, secureWss),
	}
	return DialConnect(ctx, sUrl, header, append(baseOpts, opts...)...)
}

func DialConnect(ctx context.Context, sUrl string, header http.Header, opts ...ConnOption) (*Connection, error) {
	connection := &Connection{
		id:                strconv.FormatInt(time.Now().UnixNano(), 10),
		typ:               CONN_KIND_CLIENT,
		dialer:            websocket.DefaultDialer,
		dialRetryNum:      3,
		dialRetryInterval: time.Second,
	}
	defaultNetParamsOption()(connection)

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connection)
		}
	}

	var conn *websocket.Conn
	for retry := 1; retry <= connection.dialRetryNum; retry++ {
		var err error
		var resp *http.Response
		conn, resp, err = connection.dialer.DialContext(ctx, sUrl, header)
		if err != nil {
			if retry < connection.dialRetryNum {
				log.Warn(ctx, "Failed to connect to server, sleep and try again. retry: %v, error: %v, url: %v", retry, err, sUrl)
				time.Sleep(connection.dialRetryInterval)
				continue
			} else {
				log.Error(ctx, "Failed to connect to server, leave it. retry: %v, error: %v", retry, err)
				return nil, err
			}
		}

		func() {
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
				ioutil.ReadAll(resp.Body)
			}
		}()

		log.Debug(ctx, "Success connect to server. retry: %v", retry)
		break
	}

	connection.conn = conn
	connection.commonData = make(map[string]interface{})

	if connection.pullChannelMap == nil {
		connection.pullChannelMap = make(map[int]chan struct{})
	}
	if connection.sendBuffer == nil {
		SendBufferOption(8)(connection)
	}

	ServerConnHub.register <- connection

	conn.SetCloseHandler(func(code int, text string) error {
		connection.KickServer(false)
		return nil
	})
	return connection, nil
}
