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
	"github.com/liumingmin/goutils/utils"
)

func (c *Connection) writeToServer() {
	defer func() {
		log.Debug(context.Background(), "Finish writing to server. id: %v, ptr: %p", c.id, c)
		c.closeSocket(context.Background())
	}()

	for {
		ctx := utils.ContextWithTrace()

		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				log.Debug(ctx, "Send channel closed. id: %v", c.id)
				return
			}

			if err := c.sendMsgToWs(ctx, message); err != nil {
				log.Warn(ctx, "send message failed. id: %v, error: %v", c.id, err)
				return
			}
		}
	}
}

func (c *Connection) readFromServer() {
	ctx := utils.ContextWithTrace()

	defer func() {
		if err := recover(); err != nil {
			log.Error(ctx, "readFromServer  id: %v, ptr: %p, panic :%v", c.id, c, err)
		} else {
			log.Debug(ctx, "readFromServer finished. id: %v, ptr: %p", c.id, c)
		}

		c.setStop(ctx)
		Servers.unregister <- c
	}()

	c.conn.SetReadDeadline(time.Now().Add(ReadWait))

	pingHand := c.conn.PingHandler()
	c.conn.SetPingHandler(func(message string) error {
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		if c.heartbeatCallback != nil {
			c.heartbeatCallback.RecvPing(c.id)
		}
		return pingHand(message)
	})

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		if c.heartbeatCallback != nil {
			c.heartbeatCallback.RecvPong(c.id)
		}
		return nil
	})

	c.conn.SetCloseHandler(func(code int, text string) error {
		log.Debug(context.Background(), "Connection closed. code: %v, id: %v, ptr: %p", code, c.id, c)
		return nil
	})
	c.readMsgFromWs()
}

func Connect(ctx context.Context, sUrl string, secureWss bool, header http.Header, meta *ConnectionMeta, opts ...ConnOption) (*Connection, error) {
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
				log.Warn(ctx, "Failed to connect to server, sleep and try again. retry: %v, error: %v, url\r: %v", retry, err, sUrl)
				time.Sleep(2 * time.Second)
			} else {
				log.Error(ctx, "Failed to connect to server, leave it. retry: %v, error: %v", retry, err)
				return nil, err
			}
		}

		func() {
			if resp != nil && resp.Body != nil {
				defer resp.Body.Close()
				if content, err := ioutil.ReadAll(resp.Body); err == nil {
					log.Debug(ctx, "Dial resp. resp: %v", string(content))
				}
			}
		}()

		log.Debug(ctx, "Success connect to server. retry: %v", retry)
		break
	}

	connection := &Connection{
		id:             meta.BuildConnId(),
		meta:           meta,
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
		SendBufferOption(256)(connection)
	}

	Servers.register <- connection
	log.Debug(ctx, "Client register ok. id: %v, ptr: %p", connection.id, connection)
	return connection, nil
}
