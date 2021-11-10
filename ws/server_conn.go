package ws

import (
	"context"
	"net/http"

	"github.com/liumingmin/goutils/log"
)

//displace=true，通常在集群环境下，踢掉在其他集群节点建立的连接，当前节点不需要主动调用
func (c *Connection) KickClient(displace bool) {
	if displace {
		c.setDisplaced()
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
		typ:            CONN_TYPE_CLIENT,
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
