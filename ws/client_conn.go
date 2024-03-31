package ws

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

var defaultDialer = &websocket.Dialer{
	Proxy:            http.ProxyFromEnvironment,
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   4096,
	WriteBufferSize:  4096,
}

func (c *Connection) KickServer() {
	if c.typ != CONN_KIND_CLIENT {
		return
	}

	ServerConnHub.unregisterConn(c)
}

// ctx only use for dial phase and stop auto redial
func AutoReDialConnect(ctx context.Context, sUrl string, header http.Header, connInterval time.Duration, opts ...ConnOption) {
	if connInterval == 0 {
		connInterval = time.Second * 5
	}

	for {
		conn, err := DialConnect(ctx, sUrl, header, opts...)
		if err != nil || conn == nil {
			select {
			case <-ctx.Done():
				return
			default:
			}

			time.Sleep(connInterval)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-conn.(*Connection).readDone:
			<-conn.(*Connection).writeDone
			<-conn.(*Connection).writeStop
		}

		time.Sleep(connInterval)
	}
}

// ctx only use for dial phase
func DialConnect(ctx context.Context, sUrl string, header http.Header, opts ...ConnOption) (IConnection, error) {
	connection := &Connection{}
	connection.init()

	connection.typ = CONN_KIND_CLIENT
	connection.dialer = defaultDialer
	connection.snCounter = 0

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
				if connection.dialConnFailedHandler != nil {
					connection.dialConnFailedHandler(ctx, connection)
				}
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
	connection.conn.SetCompressionLevel(connection.compressionLevel)

	if connection.sendBuffer == nil {
		SendBufferOption(8)(connection)
	}

	ServerConnHub.registerConn(connection)

	conn.SetCloseHandler(func(code int, text string) error {
		connection.KickServer()
		return nil
	})
	return connection, nil
}
