package ws

import (
	"context"
	"net/http"

	"github.com/liumingmin/goutils/net/ip"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
)

var (
	defaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (c *Connection) DisplaceClientByIp(ctx context.Context, displaceIp string) {
	c.displaceIp = displaceIp
	c.KickClient(true)
}

func (c *Connection) KickClient(displace bool) {
	if c.typ != CONN_KIND_SERVER {
		return
	}

	if displace {
		c.setDisplaced()
	}

	ClientConnHub.unregisterConn(c)
}

func AcceptGin(ctx *gin.Context, meta ConnectionMeta, opts ...ConnOption) (IConnection, error) {
	meta.clientIp = ctx.ClientIP()
	return Accept(ctx, ctx.Writer, ctx.Request, meta, opts...)
}

func Accept(ctx context.Context, w http.ResponseWriter, r *http.Request, meta ConnectionMeta, opts ...ConnOption) (IConnection, error) {
	if meta.clientIp == "" {
		meta.clientIp = ip.RemoteAddress(r)
	}

	connection := getPoolConnection()

	connection.id = meta.BuildConnId()
	connection.typ = CONN_KIND_SERVER
	connection.meta = meta
	connection.upgrader = defaultUpgrader
	connection.compressionLevel = 1
	connection.maxMessageBytesSize = defaultMaxMessageBytesSize

	defaultNetParamsOption()(connection)

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connection)
		}
	}

	conn, err := connection.upgrader.Upgrade(w, r, nil)
	if err != nil {
		putPoolConnection(connection)
		log.Warn(ctx, "%v connect failed. Header: %v, error: %v", connection.typ, r.Header, err)
		return nil, err
	}
	log.Debug(ctx, "%v connected ok. meta: %#v", connection.typ, meta)

	connection.conn = conn
	connection.conn.SetCompressionLevel(connection.compressionLevel)
	connection.commonData = make(map[string]interface{})
	connection.writeStop = make(chan interface{})
	connection.writeDone = make(chan interface{})
	connection.readDone = make(chan interface{})

	connection.createPullChannelMap()
	if connection.sendBuffer == nil {
		SendBufferOption(8)(connection)
	}

	ClientConnHub.registerConn(connection)

	conn.SetCloseHandler(func(code int, text string) error {
		connection.KickClient(false)
		return nil
	})

	return connection, nil
}
