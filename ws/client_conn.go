package ws

import (
	"context"
	"io/ioutil"
	"net/http"
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
	dialerOpts := []DialerOption{
		DialerWssOption(sUrl, secureWss),
	}

	return DialConnect(ctx, sId, sUrl, header, dialerOpts, opts...)
}

func DialConnect(ctx context.Context, sId, sUrl string, header http.Header, dialerOpts []DialerOption, opts ...ConnOption) (*Connection, error) {
	d := websocket.DefaultDialer
	d.HandshakeTimeout = handshakeTimeout

	if len(dialerOpts) > 0 {
		for _, dialerOpt := range dialerOpts {
			dialerOpt(d)
		}
	}

	var conn *websocket.Conn
	for retry := 1; retry <= connMaxRetry; retry++ {
		var err error
		var resp *http.Response
		conn, resp, err = d.DialContext(ctx, sUrl, header)
		if err != nil {
			if retry < connMaxRetry {
				log.Warn(ctx, "Failed to connect to server, sleep and try again. retry: %v, error: %v, url: %v", retry, err, sUrl)
				time.Sleep(2 * time.Second)
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

	connection := &Connection{
		id:         sId,
		typ:        CONN_KIND_CLIENT,
		conn:       conn,
		commonData: make(map[string]interface{}),
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connection)
		}
	}

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
