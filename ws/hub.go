package ws

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/liumingmin/goutils/algorithm"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/safego"
)

//连接管理器
type Hub struct {
	connections *sync.Map        // 连接容器
	register    chan *Connection // 注册队列
	unregister  chan *Connection // 注销队列
}

func newHub() *Hub {
	h := &Hub{
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: &sync.Map{},
	}
	return h
}

func (h *Hub) Find(id string) (*Connection, error) {
	if v, exists := h.connections.Load(id); exists {
		if conn, ok := v.(*Connection); ok {
			return conn, nil
		} else {
			return nil, fmt.Errorf("conn not found: %v", id)
		}
	} else {
		return nil, fmt.Errorf("conn not found: %v", id)
	}
}

func (h *Hub) RangeConnsByFunc(f func(string, *Connection) bool) {
	h.connections.Range(func(k, v interface{}) bool {
		if a, ok := v.(*Connection); ok {
			return f(k.(string), a)
		}
		return true
	})
}

func (h *Hub) ConnectionIds() []string {
	r := make([]string, 0)
	h.connections.Range(func(k, _ interface{}) bool {
		r = append(r, k.(string))
		return true
	})
	return r
}

func (h *Hub) registerConn(conn *Connection) {
	h.register <- conn
}

func (h *Hub) unregisterConn(conn *Connection) {
	h.unregister <- conn
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.processRegister(conn)
		case conn := <-h.unregister:
			h.processUnregister(conn)
		}
	}
}

func (h *Hub) processRegister(conn *Connection) {
	ctx := utils.ContextWithTsTrace()
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("processRegister. error: %v", e)
	})

	if old, err := h.Find(conn.id); err == nil && old != conn {
		// 本进程中已经存在此用户的另外一条连接，踢出老的连接
		log.Debug(ctx, "%v Repeat register, kick out. id: %v, ptr: %p", conn.typ, conn.id, old)

		old.setDisplaced()
		h.connections.Delete(old.id)
		old.setStop(ctx)

		if old.connClosingHandler != nil {
			old.connClosingHandler(ctx, old)
		}
		safego.Go(func() {
			defer old.closeSocket(ctx)
			h.sendDisplace(ctx, old, conn.ClientIp())
		})
	} else if err == nil && old == conn {
		return
	} else { // 新连接，并且是首次注册
		log.Debug(ctx, "%v new register. id: %v", conn.typ, conn.id)
	}

	h.connections.Store(conn.id, conn)

	if conn.connEstablishHandler != nil {
		conn.connEstablishHandler(ctx, conn)
	}
	log.Debug(ctx, "%v Register ok. id: %v", conn.typ, conn.id)

	safego.Go(conn.readFromConnection)
	safego.Go(conn.writeToConnection)
}

func (h *Hub) processUnregister(conn *Connection) {
	ctx := utils.ContextWithTsTrace()
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("processUnregister. error: %v", e)
	})

	if c, err := h.Find(conn.id); err == nil && c == conn {
		log.Debug(ctx, "%v unregister start. id: %v", c.typ, c.id)

		h.connections.Delete(conn.id)
		conn.setStop(ctx)

		if conn.connClosingHandler != nil {
			conn.connClosingHandler(ctx, conn)
		}
		safego.Go(func() {
			defer conn.closeSocket(ctx)
			if conn.IsDisplaced() {
				h.sendDisplace(ctx, conn, conn.displaceIp)
			}
		})

		log.Debug(ctx, "%v unregister finish. id: %v", conn.typ, conn.id)
	}
	//not in hub conn is displaced connect,do not process it
}

func (h *Hub) sendDisplace(ctx context.Context, old *Connection, newIp string) {
	if old.typ != CONN_KIND_SERVER {
		return
	}

	<-old.writeDone

	message := GetPoolMessage(int32(P_BASE_s2c_err_displace))
	dataMsg := message.DataMsg()
	if dataMsg != nil {
		displaceMsg := dataMsg.(*P_DISPLACE)
		displaceMsg.OldIp = []byte(old.ClientIp())
		displaceMsg.NewIp = []byte(newIp)
		displaceMsg.Ts = time.Now().UnixNano()
	}

	old.sendMsgToWs(ctx, message.(*Message))
}

//init server
func initServer(serverOpt ServerOption) IHub {
	RegisterDataMsgType(int32(P_BASE_s2c_err_displace), &P_DISPLACE{})

	connHub := newShardHub(serverOpt.HubOpts)
	safego.Go(connHub.run)
	return connHub
}

func newShardHub(opts []HubOption) IHub {
	sHub := &shardHub{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(sHub)
		}
	}

	if len(sHub.hubs) == 0 {
		return newHub()
	}

	return sHub
}

//init client
func initClient() IHub {
	RegisterDataMsgType(int32(P_BASE_s2c_err_displace), &P_DISPLACE{})

	RegisterHandler(int32(P_BASE_s2c_err_displace), func(ctx context.Context, conn IConnection, message IMessage) error {
		dataMsg := message.DataMsg()
		if dataMsg != nil {
			displaceMsg := dataMsg.(*P_DISPLACE)
			log.Info(ctx, "client: %v displaced by %v at %v", string(displaceMsg.OldIp), string(displaceMsg.NewIp),
				time.Unix(0, displaceMsg.Ts))
			return nil
		}

		log.Info(ctx, "client displaced")
		return nil
	})

	connHub := newHub()
	safego.Go(connHub.run)
	return connHub
}

//shard hub
type shardHub struct {
	hubs []*Hub
}

func (h *shardHub) Find(id string) (*Connection, error) {
	idx := algorithm.Crc16s(id) % uint16(len(h.hubs))
	return h.hubs[idx].Find(id)
}

func (h *shardHub) RangeConnsByFunc(rangFunc func(string, *Connection) bool) {
	for _, hub := range h.hubs {
		hub.RangeConnsByFunc(rangFunc)
	}
}

func (h *shardHub) ConnectionIds() []string {
	connIds := make([]string, 0)
	for _, hub := range h.hubs {
		connIds = append(connIds, hub.ConnectionIds()...)
	}
	return connIds
}

func (h *shardHub) registerConn(conn *Connection) {
	idx := algorithm.Crc16s(conn.Id()) % uint16(len(h.hubs))
	h.hubs[idx].registerConn(conn)
}

func (h *shardHub) unregisterConn(conn *Connection) {
	idx := algorithm.Crc16s(conn.Id()) % uint16(len(h.hubs))
	h.hubs[idx].unregisterConn(conn)
}

func (h *shardHub) run() {
	for _, hub := range h.hubs {
		safego.Go(hub.run)
	}
}
