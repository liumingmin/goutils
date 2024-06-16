package ws

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
)

var (
	msgHandlers sync.Map // map[int32]MsgHandler
	msgHeadFlag = [2]byte{0xFE, 0xEF}
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

// websocket连接封装
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
	pullChannelMap map[int]chan struct{} //pullSendNotify 拉取通知通道
	debug          bool                  //debug日志输出
	isPool         bool                  //poolObject 池对象

	snCounter uint32   //sn counter, atomic add(2), server_start=1 client_start=0
	snChanMap sync.Map //sn channel store map, map[uint32]chan IMessage

	connEstablishHandler  EventHandler
	connClosingHandler    EventHandler
	connClosedHandler     EventHandler
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

func (c *Connection) init() {
	c.id = strconv.FormatInt(time.Now().UnixNano(), 10) //default
	c.compressionLevel = 0
	c.maxMessageBytesSize = defaultMaxMessageBytesSize

	c.stopped = 0
	c.commonData = make(map[string]interface{})
	c.writeStop = make(chan interface{})
	c.writeDone = make(chan interface{})
	c.readDone = make(chan interface{})

	defaultNetParamsOption()(c)

	//client
	c.dialRetryNum = 3
	c.dialRetryInterval = time.Second
}

type ConnectionMeta struct {
	UserId   string //userId
	Typed    int    //客户端类型枚举
	DeviceId string //设备ID
	Source   string //defines where the connection comes from
	Version  int    //版本
	Charset  int    //客户端使用的字符集

	//inner set
	clientIp string
}

func (m *ConnectionMeta) BuildConnId() string {
	return fmt.Sprintf("%v-%v-%v-%v", m.UserId, m.Typed, m.DeviceId, m.Source)
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

func (c *Connection) Source() string {
	return c.meta.Source
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

func (c *Connection) reset() {
	c.id = ""
	c.meta = ConnectionMeta{}
	c.conn = nil
	c.connEstablishHandler = nil
	c.connClosingHandler = nil
	c.connClosedHandler = nil
	c.recvPingHandler = nil
	c.recvPongHandler = nil
	c.dialConnFailedHandler = nil

	c.compressionLevel = 0
	c.maxMessageBytesSize = defaultMaxMessageBytesSize

	c.stopped = 0
	c.commonData = nil
	c.writeStop = nil
	c.writeDone = nil
	c.readDone = nil

	c.maxFailureRetry = 0
	c.readWait = 0
	c.writeWait = 0
	c.temporaryWait = 0

	c.dialRetryNum = 0
	c.dialRetryInterval = 0

	c.displaced = 0
	c.displaceIp = ""

	c.sendBuffer = nil
	c.pullChannelMap = nil

	c.debug = false
	c.isPool = false

	c.snCounter = 0
	c.snChanMap = sync.Map{}

	c.upgrader = nil
	c.dialer = nil
}

func (c *Connection) GetPullChannel(pullChannelId int) (chan struct{}, bool) {
	if c.pullChannelMap == nil {
		return nil, false
	}

	v, ok := c.pullChannelMap[pullChannelId]
	return v, ok
}

func (c *Connection) IsStopped() bool {
	return atomic.LoadInt32(&c.stopped) == 1
}

func (c *Connection) setStop(ctx context.Context) {
	atomic.CompareAndSwapInt32(&c.stopped, 0, 1)

	close(c.writeStop)
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
		putPoolMessage(payload.(*Message))
		return errors.New("connect is stopped")
	}

	message := payload.(*Message)
	message.sc = sc

	c.sendBuffer <- message
	return nil
}

func (c *Connection) SendRequestMsg(ctx context.Context, reqMsg IMessage, sc SendCallback) (respMsg IMessage, err error) {
	sn := atomic.AddUint32(&c.snCounter, 2)
	if sn%2 != uint32(c.ConnType()) {
		sn += 1
	}
	(reqMsg.(*Message)).setSn(sn)

	ch := make(chan IMessage)
	c.snChanMap.Store(sn, ch)

	err = c.SendMsg(ctx, reqMsg, sc)
	if err != nil {
		c.snChanMap.Delete(sn)
		close(ch)
		return nil, err
	}

	var ok bool
	select {
	case respMsg, ok = <-ch:
		if !ok {
			err = ErrWsRpcWaitChanClosed
		}
	case <-ctx.Done():
		err = ErrWsRpcResponseTimeout
	}
	c.snChanMap.Delete(sn)

	return respMsg, err
}

func (c *Connection) SendResponseMsg(ctx context.Context, respMsg IMessage, reqSn uint32, sc SendCallback) (err error) {
	(respMsg.(*Message)).setSn(reqSn)
	return c.SendMsg(ctx, respMsg, sc)
}

func (c *Connection) SendPullNotify(ctx context.Context, pullChannelId int) (err error) {
	return c.SignalPullSend(ctx, pullChannelId)
}

// 通知指定消息通道转发消息
func (c *Connection) SignalPullSend(ctx context.Context, pullChannelId int) (err error) {
	if c.pullChannelMap == nil {
		return nil
	}

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

		ctx := log.ContextWithTraceId()

		atomic.CompareAndSwapInt32(&c.stopped, 0, 1)

		c.closePull(ctx)
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
		ctx := log.ContextWithTraceId()

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
				pingPayload = []byte(c.typ.String() + log.NewTraceId())
				log.Debug(ctx, "%v send ping. pingId: %v, ptr: %p", c.typ, string(pingPayload), c)
			}

			if err := c.conn.WriteControl(websocket.PingMessage, pingPayload, time.Now().Add(c.writeWait)); err != nil {
				if c.isNetTimeoutErr(err) {
					log.Debug(ctx, "%v send ping. timeout. id: %v, error: %v", c.typ, c.id, err)
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
		atomic.CompareAndSwapInt32(&c.stopped, 0, 1)

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
		ctx := log.ContextWithTraceId()
		if c.IsStopped() {
			log.Info(ctx, "%v connect is stopped. id: %v, ptr: %p", c.typ, c.id, c)
			return
		}

		c.conn.SetReadDeadline(time.Now().Add(c.readWait))
		messageType, data, err := c.readMessageData()

		if err != nil {
			if c.isNetTimeoutErr(err) {
				log.Debug(ctx, "%v Read failure. retryTimes: %v, id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, failedRetry, c.id, c, messageType, err)

				failedRetry++
				if failedRetry < c.maxFailureRetry {
					time.Sleep(c.temporaryWait)
					continue
				}

				log.Info(ctx, "%v Read failure and reach max times. id: %v, ptr: %p messageType: %v, error: %v",
					c.typ, c.id, c, messageType, err)
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
	return err == io.EOF
}

func (c *Connection) readMessageData() (messageType int, dataBuffer []byte, err error) {
	defer log.Recover(context.Background(), func(e interface{}) string {
		err = fmt.Errorf("%v readMessageData failed, error: %v", c.typ, e)
		return err.Error()
	})

	var reader io.Reader
	messageType, reader, err = c.conn.NextReader()
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

	dataBuffer = make([]byte, length)
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

	var message *Message
	//区分是否是rpc回包
	sn := getMsgSnFromPayload(msgData)
	if sn > 0 {
		message = &Message{}
	} else {
		message = getPoolMessage()
	}

	defer putPoolMessage(message)

	err := message.unmarshal(msgData)
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

	err := c.doSendMsgToWs(ctx, message)
	c.callback(ctx, message.sc, err)
	return err
}

func (c *Connection) doSendMsgToWs(ctx context.Context, message *Message) error {
	defer log.Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("%v doSendMsgToWs failed, error: %v", c.typ, e)
	})

	err := message.marshal()
	if err != nil {
		log.Info(ctx, "%v Marshal data message to pb failed. error: %v", c.typ, err)
		return err
	}

	c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))

	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		log.Warn(ctx, "%v Unable to get next writer of connection. error: %v", c.typ, err)
		return err
	}

	msgHeadToLEBytes := message.msgHeadToLEBytes()

	var headBytes [6]byte
	headBytes[0] = msgHeadFlag[0]
	headBytes[1] = msgHeadFlag[1]
	binary.LittleEndian.PutUint32(headBytes[2:6], uint32(len(message.data)+len(msgHeadToLEBytes)))

	_, err = w.Write(headBytes[:])
	if err != nil {
		log.Warn(ctx, "%v Write packet head to writer failed. error: %v", c.typ, err)
		return err
	}

	_, err = w.Write(msgHeadToLEBytes[:])
	if err != nil {
		log.Warn(ctx, "%v Write message head to writer failed. protocolId: %v, sn: %v, error: %v",
			c.typ, message.protocolId, message.sn, err)
		return err
	}

	if len(message.data) > 0 {
		_, err = w.Write(message.data)
		if err != nil {
			log.Warn(ctx, "%v Write data to writer failed. message: %v, error: %v", c.typ, len(message.data), err)
			return err
		}
	}

	failedRetry := 0
	for {
		if err = w.Close(); err == nil {
			log.Debug(ctx, "%v finish write message. cid: %v, message: %v", c.typ, c.id, len(message.data))
			return nil
		}

		if c.isNetTimeoutErr(err) {
			log.Debug(ctx, "%v Write close failed. retryTimes: %v, id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, err)

			failedRetry++
			if failedRetry < c.maxFailureRetry {
				time.Sleep(c.temporaryWait)
				continue
			}

			log.Info(ctx, "%v Write close failed and reach max times. id: %v, ptr: %p, error: %v",
				c.typ, failedRetry, c.id, c, err)
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
	if msg.sn > 0 {
		snChan, exist := c.snChanMap.Load(msg.sn)
		if exist && snChan != nil {
			if ch, ok := snChan.(chan IMessage); ok {
				select {
				case ch <- msg:
				default:
				}
			}
		}
	}

	handler, exist := msgHandlers.Load(msg.protocolId)
	if exist && handler != nil {
		(handler.(MsgHandler))(ctx, c, msg)
	}

	if msg.sn == 0 && !exist {
		log.Debug(ctx, "%v No handler. CMD: %d, Body len: %v", c.typ, msg.protocolId, len(msg.data))
	}
	return nil
}

// 连接数据存储结构
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

func (c *Connection) isNetTimeoutErr(err error) bool {
	var errNet net.Error
	if ok := errors.As(err, &errNet); ok && errNet.Timeout() {
		return true
	}
	return false
}
