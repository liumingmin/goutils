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
	connections *sync.Map        // 连接容器
	register    chan *Connection // 注册队列
	unregister  chan *Connection // 注销队列
}

var ClientConnHub = newHub()
var ServerConnHub = newHub()

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
		log.Debug(ctx, "%v Repeat register, kick out. id: %v, ptr: %p", conn.typ, conn.id, old)

		old.setDisplaced()
		h.connections.Delete(old.id)
		old.closeRead(ctx)

		message := GetPMessage()
		message.ProtocolId = int32(P_S2C_s2c_err_displace)

		old.SendMsg(ctx, message,
			func(cbCtx context.Context, old *Connection, e error) {
				old.setStop(cbCtx)
				old.closeSocket(cbCtx)
			})

	} else if err == nil && old == conn {
		return
	} else { // 新连接，并且是首次注册
		log.Debug(ctx, "%v new register. id: %v", conn.typ, conn.id)
	}

	h.connections.Store(conn.id, conn)

	if conn.connCallback != nil {
		log.Debug(ctx, "%v Callback ConnFinished. id: %v", conn.typ, conn.id)
		conn.connCallback.ConnFinished(conn.id)
	}
	log.Debug(ctx, "%v Register ok. id: %v", conn.typ, conn.id)

	safego.Go(conn.readFromConnection)
	safego.Go(conn.writeToConnection)
}

func (h *Hub) processUnregister(conn *Connection) {
	ctx := utils.ContextWithTrace()

	if c, err := h.findById(conn.id); err == nil && c == conn {
		log.Debug(ctx, "%v unregister start. id: %v", c.typ, c.id)

		h.connections.Delete(c.id)
		defer func() {
			c.setStop(ctx)
			c.closeSocket(ctx)
		}()

		if conn.connCallback != nil {
			if !conn.IsDisplaced() { //正常断开
				log.Debug(ctx, "%v disconnect callback. id: %v", conn.typ, conn.id)
				conn.connCallback.DisconnFinished(conn.id)
			} else {
				log.Debug(ctx, "%v displaced, skipped disconnect callback. id: %v", conn.typ, conn.id)
			}
		}

		log.Debug(ctx, "%v unregister finish. id: %v", c.typ, c.id)
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

func InitServer() {
	safego.Go(ClientConnHub.run)
}

func InitClient() {
	safego.Go(ServerConnHub.run)
}
