package rpcpool2

import (
	"container/heap"
	"sync"

	"github.com/liumingmin/goutils/log4go"
)

type ClientsHeap struct {
	clients []*Client
}

func (h *ClientsHeap) Len() int {
	return len(h.clients)
}

func (h *ClientsHeap) Less(i, j int) bool {
	t1, t2 := h.clients[i].refCnt, h.clients[j].refCnt

	return t1 < t2
}

func (h *ClientsHeap) Swap(i, j int) {
	var tmp *Client
	tmp = h.clients[i]
	h.clients[i] = h.clients[j]
	h.clients[j] = tmp
}

func (h *ClientsHeap) Push(x interface{}) {
	h.clients = append(h.clients, x.(*Client))
}

func (h *ClientsHeap) Pop() (ret interface{}) {
	l := len(h.clients)
	h.clients, ret = h.clients[:l-1], h.clients[l-1]
	return
}

type HeapPool struct {
	*Option              // 连接选项
	idle    *ClientsHeap // 空闲列表
	actives int          // 当前总数
	mutx    *sync.Mutex  // 同步锁
	cond    *sync.Cond   // 等待信号量
}

func NewHeapPool(opt *Option) (p *HeapPool, err error) {
	log4go.Debug("Start to create connect pool. option: %+v", opt)
	idle := &ClientsHeap{}

	heap.Init(idle)

	p = &HeapPool{
		Option: opt,
		idle:   idle,
	}

	var conn *Client
	for i := 0; i < opt.Size; i++ {
		conn, err = newClient(opt)
		conn.releaseFunc = p.Put
		if err != nil {
			for _, client := range idle.clients {
				client.Close()
			}
			return
		}
		heap.Push(idle, conn)
	}

	p.mutx = new(sync.Mutex)
	p.cond = sync.NewCond(p.mutx)
	p.actives = idle.Len()

	return
}

func (p *HeapPool) Close() (err error) {
	if p.idle != nil {
		for _, client := range p.idle.clients {
			client.Close()
		}
	}
	return
}

func (p *HeapPool) Get() (c *Client, err error) {
	p.mutx.Lock()
	for {
		if p.actives < p.Size {
			con, err := newClient(p.Option)
			con.releaseFunc = p.Put
			if err == nil {
				heap.Push(p.idle, con)
				p.actives++

				c = con
				c.refCnt++
				break
			} else {
				continue
			}
		}

		con := heap.Pop(p.idle).(*Client)
		if con.err != nil {
			con.Close()
			p.actives--
			continue
		} else {
			heap.Push(p.idle, con)
		}

		if con.refCnt > p.Option.RefSize {
			p.cond.Wait()
		} else {
			c = con
			c.refCnt++
			break
		}
	}

	p.mutx.Unlock()

	return
}

func (p *HeapPool) Put(c *Client, err error) {
	//if c.refCnt >= p.RefSize {
	//	fmt.Println(c, "reach max ref")
	//}
	p.mutx.Lock()
	c.refCnt--
	p.cond.Signal()
	p.mutx.Unlock()
}
