package rpcpool

import (
	"errors"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"

	"fmt"

	"time"

	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
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
	if c.refCnt > 0 {
		result := atomic.AddInt32(&c.refCnt, -1)
		if result < 0 {
			atomic.StoreInt32(&c.refCnt, 0)
		}
	}

	c.pool.cond.Broadcast()
	//fmt.Printf("%v:%v\n", unsafe.Pointer(client), client.refCnt)
}

func (c *Client) borrow(refSize int32) bool {
	if c.refCnt < refSize {
		result := atomic.AddInt32(&c.refCnt, 1)
		if result <= refSize {
			return true
		}
	}

	return false
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	err := c.Client.Call(serviceMethod, args, reply)
	if err != nil {
		if err == rpc.ErrShutdown {
			safego.Go(func() {
				c.pool.swapBadClient(c)
			})
		}
	}

	return err
}

func (c *Client) CallWithTimeout(serviceMethod string, args interface{}, reply interface{}) error {
	var err error
	var failCnt int32

	for i := 0; i < 3; i++ {
		ok := utils.AsyncInvokeWithTimeout(time.Second, func() {
			err = c.Client.Call(serviceMethod, args, reply)
		})

		if ok {
			break
		}

		failCnt = atomic.AddInt32(&c.failCnt, 1)
	}

	if err == rpc.ErrShutdown {
		safego.Go(func() {
			c.pool.swapBadClient(c)
		})
	}

	if failCnt > 6 {
		safego.Go(func() {
			c.pool.swapBadClient(c)
		})
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
	rollIdx     int
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

	clientNum := len(p.clientIdles)

	if p.option.Wait {
		for {
			circle := 0
			for {
				p.rollIdx++
				circle++

				idx := p.rollIdx % clientNum
				client := p.clientIdles[idx]

				if client.borrow(int32(p.option.RefSize)) {
					return client, nil
				}

				if circle >= clientNum {
					break
				}
			}

			p.cond.Wait()
		}
	} else {
		circle := 0
		for {
			p.rollIdx++
			circle++

			idx := p.rollIdx % clientNum
			client := p.clientIdles[idx]

			if client.borrow(int32(p.option.RefSize)) {
				return client, nil
			}

			if circle >= clientNum {
				break
			}
		}

		return nil, ErrEmpty
	}
}

func (p *Pool) swapBadClient(client *Client) {
	p.mutx.Lock()
	defer p.mutx.Unlock()

	for i := 0; i < len(p.clientIdles); i++ {
		if client == p.clientIdles[i] {
			p.clientIdles = append(p.clientIdles[:i], p.clientIdles[i+1:]...)
			break
		}
	}

	client.close()

	p.growPool()

	fmt.Println("repair a client")
}
