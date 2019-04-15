package rpcpool2

import (
	"context"
	"net"
	"net/rpc"
	"time"

	"github.com/liumingmin/goutils/utils"
)

type Client struct {
	*rpc.Client
	addr        string
	tcp         *net.TCPConn
	ctx         context.Context
	cnl         context.CancelFunc
	err         error
	refCnt      int
	releaseFunc func(*Client, error)
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	err := c.Client.Call(serviceMethod, args, reply)
	if err != nil {
		if err == rpc.ErrShutdown {
			c.err = err
		}
	}

	return err
}

func (c *Client) CallWithTimeout(serviceMethod string, args interface{}, reply interface{}) error {
	var err error

	for i := 0; i < 3; i++ {
		ok := utils.AsyncInvokeWithTimeout(time.Second, func() {
			err = c.Client.Call(serviceMethod, args, reply)
		})

		if ok {
			break
		}
	}

	if err == rpc.ErrShutdown {
		c.err = err
	}

	return err
}

func (c *Client) Release() {
	if c.releaseFunc != nil {
		c.releaseFunc(c, c.err)
	}
}

func (c *Client) Close() (err error) {
	if c.cnl != nil {
		c.cnl()
	}

	c.Client.Close()

	if c.tcp != nil {
		err = c.tcp.Close()
	}

	return
}

func newClient(opt *Option) (c *Client, err error) {
	//log4go.Debug("Start to connect to server. addr: %s", opt.Addr)

	c = &Client{
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

	//log4go.Debug("Success connect to server. addr: %s", c.addr)

	c.Client = rpc.NewClient(c.tcp)
	//log4go.Debug("Success NewClient . addr: %s", c.addr)
	return
}
