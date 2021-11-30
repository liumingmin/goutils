package ws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"google.golang.org/protobuf/proto"
)

var (
	Handlers = make(map[int32]Handler)
)

//websocket连接封装
type Connection struct {
	id         string
	typ        ConnType
	meta       ConnectionMeta   //连接信息
	conn       *websocket.Conn  //websocket connection
	sendBuffer chan interface{} //发送缓冲区  *msgSendWrapper or *P_MESSAGE

	connCallback      IConnCallback
	heartbeatCallback IHeartbeatCallback

	commonDataLock sync.RWMutex
	commonData     map[string]interface{}

	stopped   int32 //连接断开
	displaced int32 //连接被顶号

	pullChannelMap map[int]chan struct{} //新消息通知通道

	upgrader *websocket.Upgrader //可自定义upgrader
}

type ConnectionMeta struct {
	UserId   string //userId
	Typed    int    //客户端类型枚举
	DeviceId string //设备ID
	Version  int    //版本
	Charset  int    //客户端使用的字符集
}

func (m *ConnectionMeta) BuildConnId() string {
	return fmt.Sprintf("%v-%v-%v", m.UserId, m.Typed, m.DeviceId)
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) UserId() string {
	return c.meta.UserId
}

func (c *Connection) Type() int {
	return c.meta.Typed
}

func (c *Connection) Version() int {
	return c.meta.Version
}

func (c *Connection) Charset() int {
	return c.meta.Charset
}

func (c *Connection) GetPullChannel(notifyType int) (chan struct{}, bool) {
	v, ok := c.pullChannelMap[notifyType]
	return v, ok
}

func (c *Connection) IsStopped() bool {
	return atomic.LoadInt32(&c.stopped) == 1
}

func (c *Connection) setStop(ctx context.Context) {
	atomic.CompareAndSwapInt32(&c.stopped, 0, 1)

	c.closeWrite(ctx)
	c.closeRead(ctx)
}

func (c *Connection) setDisplaced() bool {
	return atomic.CompareAndSwapInt32(&c.displaced, 0, 1)
}

func (c *Connection) IsDisplaced() bool {
	return atomic.LoadInt32(&c.displaced) == 1
}

func (c *Connection) RefreshDeadline() {
	t := time.Now()
	c.conn.SetReadDeadline(t.Add(ReadWait))
	c.conn.SetWriteDeadline(t.Add(WriteWait))
}

func (c *Connection) SendMsg(ctx context.Context, payload *P_MESSAGE, sc SendCallback) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("send msg failed, sendBuffer chan is closed. error: %v", err)
	})

	if c.IsStopped() {
		return errors.New("connect is stopped")
	}

	if sc != nil {
		c.sendBuffer <- &msgSendWrapper{
			pbMessage: payload,
			sc:        sc,
		}
	} else {
		c.sendBuffer <- payload
	}

	return nil
}

//通知指定消息通道转发消息
func (c *Connection) SendPullNotify(ctx context.Context, pullChannel int) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err, _ = e.(error)
		return fmt.Sprintf("SendPullNotify err: %v", e)
	})

	if !c.IsStopped() {
		pullChannel, ok := c.pullChannelMap[pullChannel]
		if !ok {
			return
		}

		select {
		case pullChannel <- struct{}{}:
		default:
		}
	}
	return nil
}

func (c *Connection) closeWrite(ctx context.Context) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("Close writer panic, error is: %v", e)
	})

	c.commonDataLock.Lock()
	defer c.commonDataLock.Unlock()

	size := len(c.sendBuffer)
	for i := 0; i < size; i++ {
		msg, isok := <-c.sendBuffer
		if isok {
			PutPMessageIntfs(msg)
		} else {
			break
		}
	}

	select {
	case msg, isok := <-c.sendBuffer:
		if isok {
			PutPMessageIntfs(msg)
			close(c.sendBuffer)
		}
		break
	default:
		close(c.sendBuffer)
	}
}

func (c *Connection) closeRead(ctx context.Context) {
	for _, pullChannel := range c.pullChannelMap {
		func() {
			defer log.Recover(ctx, func(e interface{}) string {
				return fmt.Sprintf("Close reader panic, error is: %v", e)
			})

			select {
			case _, isok := <-pullChannel:
				if isok {
					close(pullChannel)
				}
				break
			default:
				close(pullChannel)
			}
		}()
	}
}

func (c *Connection) closeSocket(ctx context.Context) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("Close connection panic, error is: %v", e)
	})

	c.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(WriteWait))
	return c.conn.Close()
}

func (c *Connection) writeToConnection() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		log.Debug(context.Background(), "%v write finish. id: %v, ptr: %p", c.typ, c.id, c)
		ticker.Stop()

		if c.typ == CONN_TYPE_CLIENT {
			c.KickServer(false)
		} else if c.typ == CONN_TYPE_SERVER {
			c.KickClient(false)
		}
	}()

	for {
		ctx := utils.ContextWithTrace()

		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				log.Debug(ctx, "%v send channel closed. id: %v", c.typ, c.id)
				return
			}

			if e := c.sendMsgToWs(ctx, message); e != nil {
				log.Warn(ctx, "%v send message failed. id: %v, error: %v", c.typ, c.id, e)
				return
			}
		case <-ticker.C:
			log.Debug(ctx, "%v send Ping. id: %v, ptr: %p", c.typ, c.id, c)
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WriteWait)); err != nil {
				if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
					log.Debug(ctx, "%v send Ping. timeout. id: %v, error: %v", c.typ, c.id, errNet)

					time.Sleep(NetTemporaryWait)
					continue
				}

				log.Info(ctx, "%v send Ping failed. id: %v, error: %v", c.typ, c.id, c, err)
				return
			}
		}
	}
}

func (c *Connection) readFromConnection() {
	defer func() {
		log.Debug(context.Background(), "%v read finish. id: %v, ptr: %p", c.typ, c.id, c)
		if c.typ == CONN_TYPE_CLIENT {
			c.KickServer(false)
		} else if c.typ == CONN_TYPE_SERVER {
			c.KickClient(false)
		}
	}()

	c.conn.SetReadDeadline(time.Now().Add(ReadWait))

	pingHandler := c.conn.PingHandler()
	c.conn.SetPingHandler(func(message string) error {
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		err := pingHandler(message)

		if c.heartbeatCallback != nil {
			c.heartbeatCallback.RecvPing(c.id)
		}
		return err
	})
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		if c.heartbeatCallback != nil {
			c.heartbeatCallback.RecvPong(c.id)
		}
		return nil
	})
	c.readMsgFromWs()
}

func (c *Connection) readMsgFromWs() {
	failedRetry := 0

	for {
		ctx := utils.ContextWithTrace()

		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		t, data, err := c.conn.ReadMessage()
		if err != nil {
			if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
				log.Debug(ctx, "%v Read failure. retryTimes: %v, id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, failedRetry, c.id, c, t, errNet)

				failedRetry++
				if failedRetry < maxFailureRetry {
					time.Sleep(NetTemporaryWait)
					continue
				}

				log.Warn(ctx, "%v Read failure and reach max times. id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, c.id, c, t, errNet)
				break
			}

			log.Warn(ctx, "%v Conn closed or Read failed. id: %v, ptr: %p, msgType: %v, err: %v",
				c.typ, c.id, c, t, err)
			break
		}

		c.processMsg(ctx, data)
	}
}

func (c *Connection) processMsg(ctx context.Context, msgData []byte) {
	log.Debug(ctx, "%v receive raw message. data len: %v, cid: %s", c.typ, len(msgData), c.id)

	message := GetPMessage()
	defer PutPMessage(message)

	err := proto.Unmarshal(msgData, message)
	if err != nil {
		log.Error(ctx, "%v Unmarshal pb failed. data: %v, err: %v, cid: %s", c.typ, msgData, err, c.id)
		return
	}

	log.Debug(ctx, "%v receive ws message. data: %#v, cid: %s", c.typ, message, c.id)

	c.dispatch(ctx, message)
}

func (c *Connection) sendMsgToWs(ctx context.Context, message interface{}) error {
	var err error
	msg, ok := message.(*P_MESSAGE) //优先判断
	if ok {
		defer PutPMessage(msg)
		err = c.doSendMsgToWs(ctx, msg)

	} else {
		w := message.(*msgSendWrapper)
		defer PutPMessage(w.pbMessage)

		err = c.doSendMsgToWs(ctx, w.pbMessage)
		c.callback(ctx, w.sc, err)
	}

	return err
}

func (c *Connection) doSendMsgToWs(ctx context.Context, message *P_MESSAGE) error {
	c.conn.SetWriteDeadline(time.Now().Add(WriteWait))

	data, err := proto.Marshal(message)
	if err != nil {
		log.Error(ctx, "%v Marshal msgSendWrapper to pb failed. error: %v", c.typ, err)
		return err
	}

	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Warn(ctx, "%v Unable to get next writer of connection. error: %v", c.typ, err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		log.Warn(ctx, "%v Write msgSendWrapper to writer failed. message: %v, error: %v", c.typ, message, err)
		return err
	}

	failedRetry := 0
	for {
		if err := w.Close(); err == nil {
			log.Debug(ctx, "%v finish write message. cid: %v, message: %v", c.typ, c.id, message)
			return nil
		}

		if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
			log.Debug(ctx, "%v Write close failed. retryTimes: %v, id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, errNet)

			failedRetry++
			if failedRetry < maxFailureRetry {
				time.Sleep(NetTemporaryWait)
				continue
			}

			log.Warn(ctx, "%v Write close failed and reach max times. id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, errNet)
			return errors.New("writer close failed")
		}

		if e, ok := err.(*websocket.CloseError); ok {
			log.Debug(ctx, "%v Websocket close error. client id: %v, ptr: %p, error: %v",
				c.typ, c.id, c, e.Code)
		} else {
			log.Warn(ctx, "%v Writer close failed. id: %v, ptr: %p, error: %v", c.typ, c.id, c, err)
		}
		return errors.New("writer close failed")
	}

	return nil
}

func (c *Connection) callback(ctx context.Context, sc SendCallback, e error) {
	if sc != nil {
		sc(ctx, c, e)
	}
}

// 消息分发器，分发器会根据消息的协议ID查找对应的Handler。
func (c *Connection) dispatch(ctx context.Context, msg *P_MESSAGE) error {
	if h, exist := Handlers[msg.ProtocolId]; exist {
		return h(ctx, c, msg)
	} else {
		log.Error(ctx, "%v No handler. CMD: %d, Body: %s", c.typ, msg.ProtocolId, msg.Data)
		return errors.New("no handler")
	}
}

// 注册消息处理器
func RegisterHandler(cmd int32, h Handler) {
	Handlers[cmd] = h
}

//连接数据存储结构
func (c *Connection) GetCommDataValue(key string) (interface{}, bool) {
	c.commonDataLock.RLock()
	defer c.commonDataLock.RUnlock()

	if value, ok := c.commonData[key]; ok {
		return value, true
	}

	return nil, false
}

func (c *Connection) SetCommDataValue(key string, value interface{}) {
	c.commonDataLock.Lock()
	defer c.commonDataLock.Unlock()

	c.commonData[key] = value
}

func (c *Connection) RemoveCommDataValue(key string) {
	c.commonDataLock.Lock()
	defer c.commonDataLock.Unlock()

	delete(c.commonData, key)
}

func (c *Connection) IncrCommDataValueBy(key string, delta int) {
	c.commonDataLock.Lock()
	defer c.commonDataLock.Unlock()

	if value, ok := c.commonData[key]; ok {
		iValue, _ := value.(int)
		iValue += delta
		c.commonData[key] = iValue
		return
	}

	c.commonData[key] = delta
}
