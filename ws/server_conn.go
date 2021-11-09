package ws

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
)

func (c *Connection) writeToClient() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		log.Debug(context.Background(), "Finish writing to client. id: %v, ptr: %p", c.id, c)
		ticker.Stop()

		c.KickClient(false)
	}()

	for {
		ctx := utils.ContextWithTrace()

		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if e := c.sendMsgToWs(ctx, message); e != nil {
				log.Warn(ctx, "sendBuffer message to client failed. id: %v, error: %v", c.id, e)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(WriteWait)); err != nil {
				log.Warn(ctx, "Set write deadline to client failed. id：%v, error: %v", c.id, err)
			}

			log.Debug(ctx, "Send Ping. id: %v, ptr: %v", c.id, c)

			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
					log.Debug(ctx, "Ping timeout. id: %v, error: %v", c.id, errNet)

					time.Sleep(NetTemporaryWait)
					continue
				}

				log.Info(ctx, "Ping failed. id: %v, ptr: %p, error: %v", c.id, c, err)
				return
			}
		}
	}
}

func (c *Connection) readFromClient() {
	defer func() {
		log.Debug(context.Background(), "Finish reading from client. id: %v, ptr: %p", c.id, c)
		c.KickClient(false)
	}()

	c.conn.SetReadDeadline(time.Now().Add(ReadWait))

	pingHandler := c.conn.PingHandler()
	c.conn.SetPingHandler(func(message string) error {
		log.Debug(context.Background(), "Receive ping. id: %v, ptr: %p", c.id, c)
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		err := pingHandler(message)

		if c.heartbeatCallback != nil {
			c.heartbeatCallback.RecvPing(c.id)
		}
		return err
	})
	c.conn.SetPongHandler(func(string) error {
		log.Debug(context.Background(), "Receive pong. id: %v, ptr: %p", c.id, c)
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

func (c *Connection) KickClient(displace bool) {
	if displace {
		c.Displaced()
	}

	Clients.unregister <- c
}

func Accept(ctx context.Context, w http.ResponseWriter, r *http.Request, meta *ConnectionMeta, opts ...ConnOption) (*Connection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(ctx, "Client connect failed. Header: %v, error: %v", r.Header, err)
		return nil, err
	}
	log.Debug(ctx, "Client connected. meta: %#v", meta)

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

	Clients.register <- connection
	log.Debug(ctx, "Client register ok. id: %v, ptr: %p", connection.id, connection)
	return connection, nil
}

func PullChannelsOption(channels []int) ConnOption {
	return func(conn *Connection) {
		//为每种消息拉取逻辑分别注册不同的通道
		pullChannelMap := make(map[int]chan struct{})
		for _, channel := range channels {
			pullChannelMap[channel] = make(chan struct{}, 2)
		}

		conn.pullChannelMap = pullChannelMap
	}
}
