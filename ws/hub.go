package ws

import (
	"context"
	"fmt"
	"sync"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/safego"
)

//连接管理器
type Hub struct {
	connections *sync.Map        // 客户端连接
	register    chan *Connection // 注册队列
	unregister  chan *Connection // 注销队列
}

var Clients = newHub()

func newHub() *Hub {
	h := Hub{
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: &sync.Map{},
	}
	return &h
}

func (h *Hub) findById(id string) (*Connection, error) {
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

func (h *Hub) ConnectionIds() []string {
	r := make([]string, 0)
	h.connections.Range(func(k, _ interface{}) bool {
		r = append(r, k.(string))
		return true
	})
	return r
}

func (h *Hub) run() {
	for {
		//内部panic保证for循环不挂掉
		func() {
			defer log.Recover(context.Background(), func(e interface{}) string {
				return fmt.Sprintf("Hub run panic. error: %v", e)
			})

			select {
			case conn := <-h.register:
				h.processRegister(conn)
			case conn := <-h.unregister:
				h.processUnregister(conn)
			}
		}()
	}
}

func (h *Hub) processRegister(conn *Connection) {
	ctx := utils.ContextWithTrace()

	if old, err := h.findById(conn.id); err == nil && old != conn {
		// 本进程中已经存在此用户的另外一条连接，踢出老的连接
		log.Debug(ctx, "Repeat register, kick out. id: %v, ptr: %p", conn.id, old)

		old.setDisplaced()
		h.connections.Delete(old.id)
		old.closeRead(ctx)

		old.SendMsg(ctx, &P_MESSAGE{ProtocolId: int32(P_S2C_s2c_err_displace), Data: nil},
			func(cbCtx context.Context, old *Connection, e error) {
				old.Stop(cbCtx)
				old.CloseNormal(cbCtx)
			})

	} else if err == nil && old == conn {
		return
	} else { // 新连接，并且是首次注册
		log.Debug(ctx, "new client register. id: %v", conn.id)
	}

	h.connections.Store(conn.id, conn)

	if conn.connCallback != nil {
		log.Debug(ctx, "Callback ConnFinished. id: %v", conn.id)
		conn.connCallback.ConnFinished(conn.id)
	}
	log.Debug(ctx, "Register ok. id: %v", conn.id)

	safego.Go(conn.readFromClient)
	safego.Go(conn.writeToClient)
}

func (h *Hub) processUnregister(conn *Connection) {
	ctx := utils.ContextWithTrace()

	if c, err := h.findById(conn.id); err == nil && c == conn {
		log.Debug(ctx, "Unregister start. id: %v", c.id)

		h.connections.Delete(c.id)
		defer func() {
			c.Stop(ctx)
			c.CloseNormal(ctx)
		}()

		if conn.connCallback != nil {
			if !conn.IsDisplaced() { //正常断开
				log.Debug(ctx, "Callback DisconnFinished. id: %v", conn.id)
				conn.connCallback.DisconnFinished(conn.id)
			} else {
				log.Debug(ctx, "client displace closed, in hub %v, skipped disconnect callback", conn.id)
			}
		}

		log.Debug(ctx, "Unregister ok. id: %v", c.id)
	}
}

func (h *Hub) Find(id string) (*Connection, error) {
	return h.findById(id)
}

func (h *Hub) RangeConnsByFunc(f func(string, *Connection) bool) {
	h.connections.Range(func(k, v interface{}) bool {
		if a, ok := v.(*Connection); ok {
			return f(k.(string), a)
		}
		return true
	})
}

func init() {
	safego.Go(Clients.run)
}
