package rpcpool2

import (
	"container/list"
	"sync"
	"time"

	"github.com/liumingmin/goutils/log4go"
)

type Option struct {
	Addr      string // 连接地址
	Size      int    // 连接数
	RefSize   int
	KeepAlive time.Duration
}

type Pool struct {
	*Option             // 连接选项
	idle    *list.List  // 空闲列表
	actives int         // 当前总数
	mutx    *sync.Mutex // 同步锁
	cond    *sync.Cond  // 等待信号量
}

func NewPool(opt *Option) (p *Pool, err error) {
	log4go.Debug("Start to create connect pool. option: %+v", opt)
	idle := list.New()

	p = &Pool{
		Option: opt,
		idle:   idle,
	}

	var conn *Client
	for i := 0; i < opt.Size; i++ {
		conn, err = newClient(opt)
		conn.releaseFunc = p.Put
		if err != nil {
			for e := idle.Front(); e != nil; e = e.Next() {
				e.Value.(*Client).Close()
			}
			return
		}
		idle.PushBack(conn)
	}

	p.mutx = new(sync.Mutex)
	p.cond = sync.NewCond(p.mutx)
	p.actives = idle.Len()

	return
}

func (p *Pool) Close() (err error) {
	if p.idle != nil {
		for e := p.idle.Front(); e != nil; e = e.Next() {
			e.Value.(*Client).Close()
		}
	}
	return
}

func (p *Pool) Get() (c *Client, err error) {
	p.mutx.Lock()
	for p.idle.Len() == 0 && p.actives >= p.Size {
		p.cond.Wait()
	}
	if p.idle.Len() > 0 {
		c = p.idle.Remove(p.idle.Front()).(*Client)
	} else {
		c, err = newClient(p.Option)
		c.releaseFunc = p.Put
		if err == nil {
			p.actives++
		}
	}
	p.mutx.Unlock()
	return
}

func (p *Pool) Put(c *Client, err error) {
	if err != nil {
		p.mutx.Lock()
		p.actives--
		p.cond.Signal()
		p.mutx.Unlock()
		if c != nil {
			c.Close()
		}
	} else {
		p.mutx.Lock()
		p.idle.PushBack(c)
		p.cond.Signal()
		p.mutx.Unlock()
	}
}
