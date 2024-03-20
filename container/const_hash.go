package container

import (
	"bytes"
	"container/ring"
	"fmt"
	"hash/crc32"
	"math"
	"sync"
)

type CHashNode interface {
	Id() string
	Health() bool
}

type CHash interface {
	Adds([]CHashNode)
	Del(id string)
	Get(key string) CHashNode
}

type cHashSlot struct {
	node CHashNode
	slot uint32
}

type CHashRing struct {
	ring *ring.Ring
	lock sync.RWMutex
}

func (c *CHashRing) Adds(nodes []CHashNode) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.ring == nil {
		c.ring = ring.New(1)
		c.ring.Value = &cHashSlot{node: nodes[0], slot: 0}
		nodes = nodes[1:]
	}

	curr := c.ring
	next := c.ring.Next()

	for i := 0; i < len(nodes); i++ {
		cslot := curr.Value.(*cHashSlot).slot
		nslot := next.Value.(*cHashSlot).slot

		if nslot == 0 {
			nslot = ^uint32(0)
		}
		slot := (nslot-cslot)/2 + cslot

		r := ring.New(1)
		r.Value = &cHashSlot{node: nodes[i], slot: slot}

		curr.Link(r)

		curr = next
		next = next.Next()
	}
}

func (c *CHashRing) Del(nodeId string) {
	if c.ring == nil {
		return
	}

	if c.ring == c.ring.Prev() {
		return //one node cannot delete
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	slot := c.ring.Value.(*cHashSlot)
	if slot.node.Id() == nodeId {
		prev := c.ring.Prev()
		prev.Unlink(1)
		c.ring = prev
		return
	}

	for p := c.ring.Next(); p != c.ring; p = p.Next() {
		slot := p.Value.(*cHashSlot)
		if slot.node.Id() == nodeId {
			prev := p.Prev()
			prev.Unlink(1)
			return
		}
	}
}

func (c *CHashRing) GetByKey(key string, mustHealth bool) CHashNode {
	keyC32 := crc32.ChecksumIEEE([]byte(key))
	return c.GetByC32(keyC32, mustHealth)
}

func (c *CHashRing) GetByC32(keyC32 uint32, mustHealth bool) CHashNode {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var nearDistance = float64(^uint32(0))
	var nearNode CHashNode

	c.ring.Do(func(r interface{}) {
		hslot := r.(*cHashSlot)
		if mustHealth && !hslot.node.Health() {
			return
		}

		distance := math.Abs(float64(keyC32) - float64(hslot.slot))

		if distance < nearDistance {
			nearDistance = distance
			nearNode = hslot.node
		}
	})

	return nearNode
}

func (c *CHashRing) Debug() string {
	var sbuf bytes.Buffer
	c.ring.Do(func(r interface{}) {
		hslot := r.(*cHashSlot)
		sbuf.WriteString(fmt.Sprint(hslot.node, "=", hslot.slot, ","))
	})
	return sbuf.String()
}
