package tcppool

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/liumingmin/goutils/log4go"
)

type Conn struct {
	addr   string
	tcp    *net.TCPConn
	writer *bufio.Writer
	ctx    context.Context
	cnl    context.CancelFunc
	err    error // 最后发生的错误
}

func (c *Conn) Close() (err error) {
	if c.cnl != nil {
		c.cnl()
	}

	if c.tcp != nil {
		err = c.tcp.Close()
	}

	return
}

func NewConn(opt *Option) (c *Conn, err error) {
	log4go.Debug("Start to connect to server. addr: %s", opt.Addr)

	c = &Conn{
		addr: opt.Addr,
		err:  nil,
	}

	defer func() {
		if err != nil {
			if c != nil {
				c.Close()
			}
		}
	}()

	var conn net.Conn
	if conn, err = net.DialTimeout("tcp", opt.Addr, 10*time.Second); err != nil {
		return
	} else {
		c.tcp = conn.(*net.TCPConn)
	}

	c.writer = bufio.NewWriter(c.tcp)

	if err = c.tcp.SetKeepAlive(true); err != nil {
		return
	}

	if err = c.tcp.SetKeepAlivePeriod(opt.KeepAlive); err != nil {
		return
	}

	// 设置为0时，当连接需要关闭时会被直接关闭，不向对方发送CLOSE握手协议，也不会进入CLOSE_WAIT状态
	if err = c.tcp.SetLinger(0); err != nil {
		return
	}

	c.ctx, c.cnl = context.WithCancel(context.Background())

	log4go.Debug("Success connect to server. addr: %s", c.addr)

	return
}
