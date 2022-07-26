package ws

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		old.setStop(ctx)

		old.handleClosing(ctx)

		message := GetPoolMessage(int32(P_S2C_s2c_err_displace))
		dataMsg := message.DataMsg()
		if dataMsg != nil {
			displaceMsg := dataMsg.(*P_DISPLACE)
			displaceMsg.OldIp = []byte(old.ClientIp())
			displaceMsg.NewIp = []byte(conn.ClientIp())
			displaceMsg.Ts = time.Now().UnixNano()
		}
		message.sc = func(cbCtx context.Context, old *Connection, e error) {
			old.closeSocket(cbCtx)
		}
		old.sendMsgToWs(ctx, message)
	} else if err == nil && old == conn {
		return
	} else { // 新连接，并且是首次注册
		log.Debug(ctx, "%v new register. id: %v", conn.typ, conn.id)
	}

	h.connections.Store(conn.id, conn)

	if conn.connEstablishHandler != nil {
		log.Debug(ctx, "%v connEstablishHandler. id: %v", conn.typ, conn.id)
		conn.connEstablishHandler(conn)
	}
	log.Debug(ctx, "%v Register ok. id: %v", conn.typ, conn.id)

	safego.Go(conn.readFromConnection)
	safego.Go(conn.writeToConnection)
}

func (h *Hub) processUnregister(conn *Connection) {
	ctx := utils.ContextWithTrace()

	if c, err := h.findById(conn.id); err == nil && c == conn {
		log.Debug(ctx, "%v unregister start. id: %v", c.typ, c.id)

		h.connections.Delete(conn.id)
		conn.setStop(ctx)

		conn.handleClosing(ctx)
		conn.closeSocket(ctx)

		log.Debug(ctx, "%v unregister finish. id: %v", conn.typ, conn.id)
	}
	//not in hub conn is displaced connect,do not process it
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
	RegisterDataMsgType(int32(P_S2C_s2c_err_displace), &P_DISPLACE{})

	safego.Go(ClientConnHub.run)
}

func InitClient() {
	RegisterDataMsgType(int32(P_S2C_s2c_err_displace), &P_DISPLACE{})

	RegisterHandler(int32(P_S2C_s2c_err_displace), func(ctx context.Context, conn *Connection, message *Message) error {
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

	safego.Go(ServerConnHub.run)
}
