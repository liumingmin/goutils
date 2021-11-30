package ws

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

//displace=true，通常在集群环境下，踢掉在其他集群节点建立的连接，当前节点不需要主动调用
func (c *Connection) KickServer(displace bool) {
	if displace {
		c.setDisplaced()
	}

	Servers.unregister <- c
}

func Connect(ctx context.Context, sId, sUrl string, secureWss bool, header http.Header, opts ...ConnOption) (*Connection, error) {
	u, err := url.Parse(sUrl)
	if err != nil {
		log.Error(ctx, "Parse url %s err:%v", sUrl, err)
		return nil, err
	}

	var d *websocket.Dialer
	if u.Scheme == "wss" && !secureWss {
		d = &websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	} else {
		d = websocket.DefaultDialer
	}

	d.HandshakeTimeout = handshakeTimeout

	var conn *websocket.Conn
	for retry := 1; retry <= connMaxRetry; retry++ {
		var err error
		var resp *http.Response
		conn, resp, err = d.Dial(sUrl, header)
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
		id:             sId,
		typ:            CONN_TYPE_CLIENT,
		conn:           conn,
		commonData:     make(map[string]interface{}),
		pullChannelMap: make(map[int]chan struct{}),
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connection)
		}
	}

	if connection.sendBuffer == nil {
		SendBufferOption(8)(connection)
	}

	Servers.register <- connection

	conn.SetCloseHandler(func(code int, text string) error {
		connection.KickServer(false)
		return nil
	})
	return connection, nil
}
