package rpcpool

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"unsafe"
)

var ErrClosed = errors.New("pool is closed")
var ErrEmpty = errors.New("pool is empty")
var ErrConnect = errors.New("connect failed")

type Option struct {
	RpcSize int
	RefSize int
	Wait    bool
}

type Client struct {
	*rpc.Client
	conn    net.Conn
	pool    *Pool
	failCnt int32
	refCnt  int32
}

func (c *Client) Release() {
	c.pool.release(c)
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	err := c.Client.Call(serviceMethod, args, reply)
	if err != nil {
		atomic.AddInt32(&c.failCnt, 1)
	}

	return err
}

func (c *Client) close() {
	c.Client.Close()
	c.conn.Close()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
type Pool struct {
	option  Option
	factory func() (net.Conn, error)

	mutx        *sync.Mutex
	cond        *sync.Cond
	clientIdles []*Client
}

func (p *Pool) newClient() (*Client, error) {
	con, err := p.factory()
	if err == nil && con != nil {
		rpcClient := rpc.NewClient(con)
		client := &Client{Client: rpcClient, conn: con, pool: p}
		return client, nil
	}

	return nil, ErrConnect
}

func (p *Pool) growPool() {
	if client, err := p.newClient(); err == nil {
		p.clientIdles = append(p.clientIdles, client)
	}
}

func (p *Pool) Init(opt Option, f func() (net.Conn, error)) {
	p.option = opt
	p.factory = f
	p.mutx = new(sync.Mutex)
	p.cond = sync.NewCond(p.mutx)
	p.clientIdles = make([]*Client, 0, p.option.RpcSize)

	for {
		if len(p.clientIdles) < p.option.RpcSize {
			p.growPool()
		} else {
			return
		}
	}
}

func (p *Pool) Get() (*Client, error) {
	p.mutx.Lock()
	defer p.mutx.Unlock()

	if p.option.Wait {
		for {
			for i := 0; i < len(p.clientIdles); i++ {
				client := p.clientIdles[i]
				if client.refCnt <= int32(p.option.RefSize) {
					client.refCnt++
					return client, nil
				}
			}

			p.cond.Wait()
		}
	} else {
		for i := 0; i < len(p.clientIdles); i++ {
			client := p.clientIdles[i]
			if client.refCnt <= int32(p.option.RefSize) {
				client.refCnt++
				return client, nil
			}
		}

		return nil, ErrEmpty
	}
}

func (p *Pool) release(client *Client) {
	if client == nil {
		return
	}

	p.mutx.Lock()
	defer p.mutx.Unlock()

	if client.failCnt > 3 {
		for i := 0; i < len(p.clientIdles); i++ {
			if client == p.clientIdles[i] {
				p.clientIdles = append(p.clientIdles[:i], p.clientIdles[i+1:]...)
				break
			}
		}
		client.close()

		p.growPool()
	} else {
		if client.refCnt > 0 {
			client.refCnt--
		}
	}

	fmt.Printf("%v:%v\n", unsafe.Pointer(client), client.refCnt)
	p.cond.Broadcast()
}
