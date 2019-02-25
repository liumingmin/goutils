package tcppool

import (
	"container/list"
	"sync"
	"time"

	"github.com/liumingmin/goutils/log4go"
)

type Option struct {
	Addr        string        // 连接地址
	Size        int           // 连接数
	ReadTimeout time.Duration // 读超时秒数
	KeepAlive   time.Duration
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
	var conn *Conn
	for i := 0; i < opt.Size; i++ {
		conn, err = NewConn(opt)
		if err != nil {
			for e := idle.Front(); e != nil; e = e.Next() {
				e.Value.(*Conn).Close()
			}
			return
		}
		idle.PushBack(conn)
	}

	mutx := new(sync.Mutex)
	cond := sync.NewCond(mutx)
	p = &Pool{
		Option:  opt,
		idle:    idle,
		actives: idle.Len(),
		mutx:    mutx,
		cond:    cond,
	}

	return
}

func (p *Pool) Close() (err error) {
	if p.idle != nil {
		for e := p.idle.Front(); e != nil; e = e.Next() {
			e.Value.(*Conn).Close()
		}
	}
	return
}

func (p *Pool) Get() (c *Conn, err error) {
	p.mutx.Lock()
	for p.idle.Len() == 0 && p.actives >= p.Size {
		p.cond.Wait()
	}
	if p.idle.Len() > 0 {
		c = p.idle.Remove(p.idle.Front()).(*Conn)
	} else {
		c, err = NewConn(p.Option)
		if err == nil {
			p.actives++
		}
	}
	p.mutx.Unlock()
	return
}

func (p *Pool) Put(c *Conn, err error) {
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
