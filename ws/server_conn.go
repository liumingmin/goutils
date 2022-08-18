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
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//displace=true，通常在集群环境下，踢掉在其他集群节点建立的连接，当前节点不需要主动调用
func (c *Connection) KickClient(displace bool) {
	if displace {
		c.setDisplaced()
	}

	ClientConnHub.unregisterConn(c)
}

func AcceptGin(ctx *gin.Context, meta ConnectionMeta, opts ...ConnOption) (*Connection, error) {
	meta.ip = ctx.ClientIP()
	return Accept(ctx, ctx.Writer, ctx.Request, meta, opts...)
}

func Accept(ctx context.Context, w http.ResponseWriter, r *http.Request, meta ConnectionMeta, opts ...ConnOption) (*Connection, error) {
	if meta.ip == "" {
		meta.ip = ip.RemoteAddress(r)
	}

	connection := getPoolSrvConnection()

	connection.id = meta.BuildConnId()
	connection.typ = CONN_KIND_SERVER
	connection.meta = meta
	connection.upgrader = defaultUpgrader
	connection.compressionLevel = 1

	defaultNetParamsOption()(connection)

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connection)
		}
	}

	conn, err := connection.upgrader.Upgrade(w, r, nil)
	if err != nil {
		putPoolSrvConnection(connection)
		log.Warn(ctx, "%v connect failed. Header: %v, error: %v", connection.typ, r.Header, err)
		return nil, err
	}
	log.Debug(ctx, "%v connected ok. meta: %#v", connection.typ, meta)

	connection.conn = conn
	connection.conn.SetCompressionLevel(connection.compressionLevel)
	connection.commonData = make(map[string]interface{})

	if connection.pullChannelMap == nil {
		connection.pullChannelMap = make(map[int]chan struct{})
	}
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
