package ws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/safego"
	"google.golang.org/protobuf/proto"
)

// 消息发送回调接口
type SendCallback func(ctx context.Context, c *Connection, err error)

// 客户端消息处理函数对象
type Handler func(context.Context, *Connection, *P_MESSAGE) error

// 连接动态参数选项
type ConnOption func(*Connection)

//配置项
var (
	dispatcherNum   = conf.ExtInt("ws.dispatcherNum", 16)              //并发处理消息数量
	maxFailureRetry = conf.ExtInt("ws.maxFailureRetry", 10)            //重试次数
	ReadWait        = conf.ExtDuration("ws.readWait", 60*time.Second)  //读等待
	WriteWait       = conf.ExtDuration("ws.writeWait", 60*time.Second) //写等待
	PingPeriod      = WriteWait * 4 / 10                               //ping间隔应该小于写等待时间

	NetTemporaryWait = 500 * time.Millisecond //网络抖动重试等待
)

var (
	Handlers = make(map[int32]Handler)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//连接回调
type IConnCallback interface {
	ConnFinished(clientId string)
	DisconnFinished(clientId string)
}

//保活回调
type IHeartbeatCallback interface {
	RecvPing(clientId string)
	RecvPong(clientId string) error
}

//websocket连接封装
type Connection struct {
	id         string
	meta       *ConnectionMeta //连接信息
	conn       *websocket.Conn // websocket connection
	sendBuffer chan *Message   //

	connCallback      IConnCallback
	heartbeatCallback IHeartbeatCallback

	commonDataLock sync.RWMutex
	commonData     map[string]interface{}

	stopped   int32 //连接断开
	displaced int32 //连接被顶号

	pullChannelMap map[int]chan struct{} //新消息通知通道
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

type Message struct {
	pbMessage *P_MESSAGE   // 消息体
	sc        SendCallback // 消息发送回调接口
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

func (c *Connection) Stop(ctx context.Context) {
	atomic.CompareAndSwapInt32(&c.stopped, 0, 1)

	c.closeWrite(ctx)
	c.closeRead(ctx)
}

func (c *Connection) Displaced() bool {
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

	c.sendBuffer <- &Message{
		pbMessage: payload,
		sc:        sc,
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

	select {
	case _, isok := <-c.sendBuffer:
		if isok {
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

func (c *Connection) CloseNormal(ctx context.Context) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("Close connection panic, error is: %v", e)
	})

	return c.conn.Close()
}

func (c *Connection) readMsgFromWs() {
	failedRetry := 0

	var pLimit chan struct{}

	pLimit = make(chan struct{}, dispatcherNum)

	for {
		ctx := utils.ContextWithTrace()

		c.conn.SetReadDeadline(time.Now().Add(ReadWait))
		t, data, err := c.conn.ReadMessage()
		if err != nil {
			if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
				log.Debug(ctx, "Read failure. retryTimes: %v, id: %v, ptr: %p messageType: %v, error: %v",
					failedRetry, c.id, c, t, errNet)

				failedRetry++
				if failedRetry < maxFailureRetry {
					time.Sleep(NetTemporaryWait)
					continue
				}

				log.Warn(ctx, "Read failure and reach max times. id: %v, ptr: %p messageType: %v, error: %v",
					c.id, c, t, errNet)
				break
			}

			log.Warn(ctx, "Conn closed or Read failed. id: %v, ptr: %p, msgType: %v, err: %v",
				c.id, c, t, err)
			break
		}

		c.processMsg(ctx, pLimit, data)
	}
}

func (c *Connection) processMsg(ctx context.Context, pLimit chan struct{}, msgData []byte) {
	log.Debug(ctx, "receive raw message. data len: %v, cid: %s", len(msgData), c.id)

	var message P_MESSAGE
	err := proto.Unmarshal(msgData, &message)
	if err != nil {
		log.Error(ctx, "Unmarshal pb failed. data: %v, err: %v, cid: %s", msgData, err, c.id)
		return
	}

	log.Debug(ctx, "receive ws message. data: %#v, cid: %s", message, c.id)

	pLimit <- struct{}{}
	safego.Go(func() {
		defer func() { <-pLimit }()
		c.dispatch(ctx, &message)
	})
}

func (c *Connection) sendMsgToWs(ctx context.Context, message *Message) error {
	err := c.doSendMsgToWs(ctx, message)
	c.callback(ctx, message.sc, err)
	return err
}

func (c *Connection) doSendMsgToWs(ctx context.Context, message *Message) error {
	c.conn.SetWriteDeadline(time.Now().Add(WriteWait))

	data, err := proto.Marshal(message.pbMessage)
	if err != nil {
		log.Error(ctx, "Marshal Message to pb failed. error: %v", err)
		return err
	}

	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Warn(ctx, "Unable to get next writer of connection. error: %v", err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		log.Warn(ctx, "Write Message to writer failed. message: %v, error: %v", message, err)
		return err
	}

	failedRetry := 0
	for {
		if err := w.Close(); err == nil {
			log.Debug(ctx, "finish write message to connect. cid: %v, message: %v", c.id, message)
			return nil
		}

		if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
			log.Debug(ctx, "Write close failed. retryTimes: %v, id: %v, ptr: %p, error: %v",
				failedRetry, c.id, c, errNet)

			failedRetry++
			if failedRetry < maxFailureRetry {
				time.Sleep(NetTemporaryWait)
				continue
			}

			log.Warn(ctx, "Write close failed and reach max times. id: %v, ptr: %p, error: %v",
				failedRetry, c.id, c, errNet)
			return errors.New("writer close failed")
		}

		if e, ok := err.(*websocket.CloseError); ok {
			log.Debug(ctx, "Websocket close error. client id: %v, ptr: %p, error: %v",
				c.id, c, e.Code)
		} else {
			log.Warn(ctx, "Writer close failed. id: %v, ptr: %p, error: %v", c.id, c, err)
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
		log.Error(ctx, "No handler. CMD: %d, Body: %s", msg.ProtocolId, msg.Data)
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

func ConnectCbOption(connCallback IConnCallback) ConnOption {
	return func(conn *Connection) {
		conn.connCallback = connCallback
	}
}

func HeartbeatCbOption(heartbeatCallback IHeartbeatCallback) ConnOption {
	return func(conn *Connection) {
		conn.heartbeatCallback = heartbeatCallback
	}
}
