package ws

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
	"github.com/liumingmin/goutils/utils/safego"
)

var (
	msgHandlers sync.Map // map[int32]MsgHandler
	msgHeadFlag = [2]byte{0xFE, 0xEE}
)

type ConnType int8

func (t ConnType) String() string {
	if t == CONN_KIND_CLIENT {
		return "client"
	}
	if t == CONN_KIND_SERVER {
		return "server"
	}
	return ""
}

//websocket连接封装
type Connection struct {
	id             string
	typ            ConnType
	meta           ConnectionMeta        //连接信息
	conn           *websocket.Conn       //websocket connection
	stopped        int32                 //flag connection stopped and will disconnect
	writeStop      chan interface{}      //writeConnection loop stop
	writeDone      chan interface{}      //writeConnection finished
	readDone       chan interface{}      //readFromConnection finished
	displaced      int32                 //连接被顶号
	displaceIp     string                //displaced by ip(cluster use) 顶号IP(集群下使用)
	sendBuffer     chan *Message         //发送缓冲区
	pullChannelIds []int                 //to construct pullChannelMap
	pullChannelMap map[int]chan struct{} //pullSendNotify 拉取通知通道
	debug          bool                  //debug日志输出
	isPool         bool                  //poolObject 池对象

	connEstablishHandler  EventHandler
	connClosingHandler    EventHandler
	connClosedHandler     EventHandler
	connAutoReconHandler  EventHandler
	recvPingHandler       EventHandler
	recvPongHandler       EventHandler
	dialConnFailedHandler EventHandler

	commonDataLock sync.RWMutex
	commonData     map[string]interface{}

	//net params
	maxFailureRetry     int           //重试次数
	readWait            time.Duration //读等待
	writeWait           time.Duration //写等待
	temporaryWait       time.Duration //网络抖动重试等待
	compressionLevel    int
	maxMessageBytesSize uint32

	//server internal param
	upgrader *websocket.Upgrader //custome upgrader

	//client internal param
	dialer            *websocket.Dialer
	dialRetryNum      int
	dialRetryInterval time.Duration
}

type ConnectionMeta struct {
	UserId   string //userId
	Typed    int    //客户端类型枚举
	DeviceId string //设备ID
	Version  int    //版本
	Charset  int    //客户端使用的字符集

	//inner set
	clientIp string
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

func (c *Connection) DeviceId() string {
	return c.meta.DeviceId
}

func (c *Connection) Version() int {
	return c.meta.Version
}

func (c *Connection) Charset() int {
	return c.meta.Charset
}

func (c *Connection) ClientIp() string {
	return c.meta.clientIp
}

func (c *Connection) ConnType() ConnType {
	return c.typ
}

func (c *Connection) Reset() {
	c.id = ""
	c.meta = ConnectionMeta{}
	c.conn = nil
	c.connEstablishHandler = nil
	c.connClosingHandler = nil
	c.connClosedHandler = nil
	c.connAutoReconHandler = nil
	c.recvPingHandler = nil
	c.recvPongHandler = nil
	c.dialConnFailedHandler = nil
	c.commonData = nil
	c.stopped = 0
	c.writeStop = nil
	c.writeDone = nil
	c.readDone = nil
	c.displaced = 0
	c.displaceIp = ""

	c.sendBuffer = nil
	c.pullChannelIds = nil
	c.pullChannelMap = nil
	c.compressionLevel = 0
	c.debug = false
	c.isPool = false

	c.maxFailureRetry = 0
	c.readWait = 0
	c.writeWait = 0
	c.temporaryWait = 0

	c.upgrader = nil

	c.dialer = nil
	c.dialRetryNum = 0
	c.dialRetryInterval = 0
	c.maxMessageBytesSize = defaultMaxMessageBytesSize
}

func (c *Connection) createPullChannelMap() {
	pullChannelMap := make(map[int]chan struct{})

	if len(c.pullChannelIds) > 0 {
		for _, channel := range c.pullChannelIds {
			pullChannelMap[channel] = make(chan struct{}, 1)
		}
	}

	c.pullChannelMap = pullChannelMap
}

func (c *Connection) GetPullChannel(pullChannelId int) (chan struct{}, bool) {
	v, ok := c.pullChannelMap[pullChannelId]
	return v, ok
}

func (c *Connection) IsStopped() bool {
	return atomic.LoadInt32(&c.stopped) == 1
}

func (c *Connection) setStop(ctx context.Context) {
	ok := atomic.CompareAndSwapInt32(&c.stopped, 0, 1)
	if !ok {
		return
	}

	close(c.writeStop)
	c.closePull(ctx)
}

func (c *Connection) setDisplaced() bool {
	return atomic.CompareAndSwapInt32(&c.displaced, 0, 1)
}

func (c *Connection) IsDisplaced() bool {
	return atomic.LoadInt32(&c.displaced) == 1
}

func (c *Connection) RefreshDeadline() {
	t := time.Now()
	c.conn.SetReadDeadline(t.Add(c.readWait))
	c.conn.SetWriteDeadline(t.Add(c.writeWait))
}

func (c *Connection) SendMsg(ctx context.Context, payload IMessage, sc SendCallback) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("%v send msg failed, sendBuffer chan is closed. error: %v", c.typ, err)
	})

	if c.IsStopped() {
		return errors.New("connect is stopped")
	}

	message := payload.(*Message)
	message.sc = sc

	c.sendBuffer <- message
	return nil
}

//通知指定消息通道转发消息
func (c *Connection) SendPullNotify(ctx context.Context, pullChannelId int) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err, _ = e.(error)
		return fmt.Sprintf("%v SendPullNotify err: %v", c.typ, e)
	})

	if c.IsStopped() {
		return errors.New("connect is stopped")
	}

	pullChannel, ok := c.pullChannelMap[pullChannelId]
	if !ok {
		return
	}

	select {
	case pullChannel <- struct{}{}:
	default:
	}

	return nil
}

func (c *Connection) closeWrite(ctx context.Context) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v Close writer panic, error is: %v", c.typ, e)
	})

	size := len(c.sendBuffer)
	for i := 0; i < size; i++ {
		msg, isok := <-c.sendBuffer
		if isok {
			putPoolMessage(msg)
		} else {
			return //buffer has closed
		}
	}
	close(c.sendBuffer)
}

func (c *Connection) closePull(ctx context.Context) {
	if c.pullChannelMap == nil {
		return
	}

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

func (c *Connection) handleClosed(ctx context.Context) {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v handleClosed panic, error is: %v", c.typ, e)
	})

	<-c.writeDone
	<-c.readDone

	defer putPoolConnection(c)

	if c.connClosedHandler != nil {
		c.connClosedHandler(ctx, c)
	}

	if c.connAutoReconHandler != nil {
		c.connAutoReconHandler(ctx, c)
	}
}

func (c *Connection) closeSocket(ctx context.Context) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v Close connection panic, error is: %v", c.typ, e)
	})

	defer c.handleClosed(ctx)

	c.conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(c.writeWait))

	var err error
	for i := 0; i < 3; i++ {
		err = c.conn.Close()
		if err == nil {
			return nil
		}
	}
	return err
}

func (c *Connection) writeToConnection() {
	ticker := time.NewTicker(c.writeWait * 4 / 10) //PingPeriod
	defer func() {
		ticker.Stop()

		ctx := utils.ContextWithTsTrace()
		c.setStop(ctx)
		c.closeWrite(ctx)
		close(c.writeDone)

		log.Debug(ctx, "%v write finish. id: %v, ptr: %p", c.typ, c.id, c)

		if c.typ == CONN_KIND_CLIENT {
			c.KickServer()
		} else if c.typ == CONN_KIND_SERVER {
			c.KickClient(false)
		}
	}()

	pingPayload := []byte{}

	for {
		ctx := utils.ContextWithTsTrace()

		if c.IsStopped() {
			log.Info(ctx, "%v connection stopped, finish write. id: %v, ptr: %p", c.typ, c.id, c)
			return
		}

		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				return
			}

			if err := c.sendMsgToWs(ctx, message); err != nil {
				log.Debug(ctx, "%v send message failed. id: %v, error: %v", c.typ, c.id, err)
				return
			}

			ok = c.batchSendMsgToWs(ctx)
			if !ok {
				return
			}

		case <-ticker.C:
			if c.debug {
				pingPayload = []byte(c.typ.String() + utils.NanoTsBase36() + utils.RandBase36())
				log.Debug(ctx, "%v send ping. pingId: %v, ptr: %p", c.typ, string(pingPayload), c)
			}

			if err := c.conn.WriteControl(websocket.PingMessage, pingPayload, time.Now().Add(c.writeWait)); err != nil {
				if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
					log.Debug(ctx, "%v send ping. timeout. id: %v, error: %v", c.typ, c.id, errNet)
					continue
				}

				log.Warn(ctx, "%v send Ping failed. id: %v, error: %v", c.typ, c.id, c, err)
				return
			}
		case <-c.writeStop:
			return
		}
	}
}

func (c *Connection) readFromConnection() {
	defer func() {
		close(c.readDone)
		log.Debug(context.Background(), "%v read finish. id: %v, ptr: %p", c.typ, c.id, c)

		if c.typ == CONN_KIND_CLIENT {
			c.KickServer()
		} else if c.typ == CONN_KIND_SERVER {
			c.KickClient(false)
		}
	}()

	c.conn.SetReadDeadline(time.Now().Add(c.readWait))

	pingHandler := c.conn.PingHandler()
	c.conn.SetPingHandler(func(message string) error {
		if c.IsStopped() {
			log.Info(context.Background(), "%v recv ping. pingId: %v, ptr: %p, connect is stopped", c.typ, message, c)
			return nil
		}

		c.conn.SetReadDeadline(time.Now().Add(c.readWait))
		err := pingHandler(message)

		if c.recvPingHandler != nil {
			c.recvPingHandler(context.Background(), c)
		}

		if c.debug {
			log.Debug(context.Background(), "%v recv ping. pingId: %v, ptr: %p", c.typ, message, c)
		}
		return err
	})
	c.conn.SetPongHandler(func(message string) error {
		if c.IsStopped() {
			log.Info(context.Background(), "%v recv pong. pingId: %v, ptr: %p, connect is stopped", c.typ, message, c)
			return nil
		}

		c.conn.SetReadDeadline(time.Now().Add(c.readWait))
		if c.recvPongHandler != nil {
			c.recvPongHandler(context.Background(), c)
		}

		if c.debug {
			log.Debug(context.Background(), "%v recv pong. pingId: %v, ptr: %p", c.typ, message, c)
		}
		return nil
	})
	c.readMsgFromWs()
}

func (c *Connection) readMsgFromWs() {
	defer log.Recover(context.Background(), func(e interface{}) string {
		return fmt.Sprintf("%v readMsgFromWs failed, error: %v", c.typ, e)
	})

	failedRetry := 0

	for {
		ctx := utils.ContextWithTsTrace()
		if c.IsStopped() {
			log.Info(ctx, "%v connect is stopped. id: %v, ptr: %p", c.typ, c.id, c)
			return
		}

		c.conn.SetReadDeadline(time.Now().Add(c.readWait))
		messageType, data, err := c.readMessageData() //c.conn.ReadMessage()

		if err != nil {
			if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
				log.Debug(ctx, "%v Read failure. retryTimes: %v, id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, failedRetry, c.id, c, messageType, errNet)

				failedRetry++
				if failedRetry < c.maxFailureRetry {
					time.Sleep(c.temporaryWait)
					continue
				}

				log.Info(ctx, "%v Read failure and reach max times. id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, c.id, c, messageType, errNet)
				return
			}

			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived) || c.IsStopped() {
				log.Debug(ctx, "%v Conn closed. id: %v, ptr: %p, msgType: %v, err: %v",
					c.typ, c.id, c, messageType, err)
			} else {
				log.Warn(ctx, "%v Read failed. id: %v, ptr: %p, msgType: %v, err: %v",
					c.typ, c.id, c, messageType, err)
			}
			return
		}

		if messageType == websocket.BinaryMessage && len(data) > 0 {
			safego.Go(func() {
				c.processMsg(ctx, data)
			})
		}
	}
}

func (c *Connection) isErrEOF(err error) bool {
	if err == io.EOF {
		return true
	}
	return false
}

func (c *Connection) readMessageData() (int, []byte, error) {
	var reader io.Reader
	messageType, reader, err := c.conn.NextReader()
	if err != nil && !c.isErrEOF(err) {
		return messageType, nil, err
	}

	var headBytes [6]byte
	_, err = io.ReadAtLeast(reader, headBytes[:], 6)
	if err != nil && !c.isErrEOF(err) {
		return messageType, nil, err
	}

	if headBytes[0] != msgHeadFlag[0] || headBytes[1] != msgHeadFlag[1] {
		return messageType, nil, errors.New("packet head flag error")
	}

	lengthSlice := headBytes[2:6]
	var length uint32
	binary.Read(bytes.NewReader(lengthSlice), binary.LittleEndian, &length)

	if length > c.maxMessageBytesSize {
		return messageType, nil, errors.New("packet size exceed max")
	}

	dataBuffer := make([]byte, length)
	_, err = io.ReadAtLeast(reader, dataBuffer, int(length))
	if err != nil && !c.isErrEOF(err) {
		return messageType, nil, err
	}
	return messageType, dataBuffer, nil
}

func (c *Connection) processMsg(ctx context.Context, msgData []byte) {
	if c.debug {
		log.Debug(ctx, "%v receive raw message. data len: %v, cid: %v", c.typ, msgData, c.id)
	}

	message := getPoolMessage()
	defer putPoolMessage(message)

	err := message.Unmarshal(msgData)
	if err != nil {
		log.Error(ctx, "%v Unmarshal pb failed. data: %v, err: %v, cid: %v", c.typ, msgData, err, c.id)
		return
	}

	if c.debug {
		log.Debug(ctx, "%v receive ws message. data: %#v, cid: %s", c.typ, message, c.id)
	}

	c.dispatch(ctx, message)
}

func (c *Connection) batchSendMsgToWs(ctx context.Context) bool {
	size := len(c.sendBuffer)
	for i := 0; i < size; i++ {
		message, ok := <-c.sendBuffer
		if !ok {
			return false
		}
		if err := c.sendMsgToWs(ctx, message); err != nil {
			log.Debug(ctx, "%v batchSendMsgToWs failed. id: %v, error: %v", c.typ, c.id, err)
			return false
		}
	}
	return true
}

func (c *Connection) sendMsgToWs(ctx context.Context, message *Message) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v sendMsgToWs failed, error: %v", c.typ, e)
	})

	defer putPoolMessage(message)

	if message.PMsg().ProtocolId == int32(P_BASE_raw_bytes_msg) {
		err := c.doSendMsgToWs(ctx, message.PMsg().Data)
		c.callback(ctx, message.sc, err)
		return err
	}

	msgData, err := message.Marshal()
	if err != nil {
		log.Warn(ctx, "%v Marshal message to pb failed. error: %v", c.typ, err)
		c.callback(ctx, message.sc, err)
		return err
	}

	err = c.doSendMsgToWs(ctx, msgData)
	c.callback(ctx, message.sc, err)
	return err
}

func (c *Connection) doSendMsgToWs(ctx context.Context, data []byte) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v doSendMsgToWs failed, error: %v", c.typ, e)
	})

	c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))

	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Warn(ctx, "%v Unable to get next writer of connection. error: %v", c.typ, err)
		return err
	}

	var headBytes [6]byte
	headBytes[0] = msgHeadFlag[0]
	headBytes[1] = msgHeadFlag[1]
	binary.LittleEndian.PutUint32(headBytes[2:6], uint32(len(data)))

	_, err = w.Write(headBytes[:])
	if err != nil {
		log.Warn(ctx, "%v Write packet head to writer failed. error: %v", c.typ, err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		log.Warn(ctx, "%v Write payload to writer failed. message: %v, error: %v", c.typ, len(data), err)
		return err
	}

	failedRetry := 0
	for {
		if err := w.Close(); err == nil {
			log.Debug(ctx, "%v finish write message. cid: %v, message: %v", c.typ, c.id, len(data))
			return nil
		}

		if errNet, ok := err.(net.Error); (ok && errNet.Timeout()) || (ok && errNet.Temporary()) {
			log.Debug(ctx, "%v Write close failed. retryTimes: %v, id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, errNet)

			failedRetry++
			if failedRetry < c.maxFailureRetry {
				time.Sleep(c.temporaryWait)
				continue
			}

			log.Info(ctx, "%v Write close failed and reach max times. id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, errNet)
			return errors.New("writer close failed")
		}

		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway,
			websocket.CloseNoStatusReceived) || c.IsStopped() {
			log.Debug(ctx, "%v Websocket close error. client id: %v, ptr: %p, error: %v",
				c.typ, c.id, c, err)
		} else {
			log.Warn(ctx, "%v Writer close failed. id: %v, ptr: %p, error: %v", c.typ, c.id, c, err)
		}
		return errors.New("writer close failed")
	}
}

func (c *Connection) callback(ctx context.Context, sc SendCallback, e error) {
	if sc != nil {
		sc(ctx, c, e)
	}
}

// 消息分发器，分发器会根据消息的协议ID查找对应的Handler。
func (c *Connection) dispatch(ctx context.Context, msg *Message) error {
	handler, exist := msgHandlers.Load(msg.pMsg.ProtocolId)
	if exist && handler != nil {
		return (handler.(MsgHandler))(ctx, c, msg)
	} else {
		log.Debug(ctx, "%v No handler. CMD: %d, Body len: %v", c.typ, msg.pMsg.ProtocolId, len(msg.pMsg.Data))
		return errors.New("no handler")
	}
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
