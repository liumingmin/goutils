package container

import (
	"container/ring"
	"fmt"
	"hash/crc32"
	"sync"
)

type NodeHealth interface {
	Health() bool
}

type Chash struct {
	ring      *ring.Ring
	nodeRing  *ring.Ring
	nodeCount int
	lock      sync.RWMutex
}

type chashSlot struct {
	start uint32
	end   uint32
	node  NodeHealth
}

func NewChash(nodes []NodeHealth) *Chash {
	initLen := len(nodes)
	r := ring.New(initLen)
	step := ^uint32(0) / uint32(initLen)

	next := r
	for i := 0; i < initLen; i++ {
		next = next.Next()
		slot := &chashSlot{start: uint32(i) * step, end: uint32(i+1) * step, node: nodes[i]}

		if i == initLen-1 {
			slot.end = ^uint32(0)
		}

		next.Value = slot
	}

	return &Chash{
		ring:      r,
		nodeRing:  r,
		nodeCount: initLen,
	}
}

func (c *Chash) AddNode(node NodeHealth) {
	c.lock.Lock()
	defer c.lock.Unlock()

	next := c.ring.Next()

	slot := next.Value.(*chashSlot)
	newstep := (slot.end - slot.start) / 2

	r := ring.New(1)
	r.Value = &chashSlot{start: slot.start, end: slot.start + newstep, node: node}

	c.ring.Link(r)
	c.nodeCount++
	slot.start = slot.start + newstep
	c.ring = next
}

func (c *Chash) GetNode(key string) NodeHealth {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c32 := crc32.ChecksumIEEE([]byte(key))

	currSlot := c.nodeRing.Value.(*chashSlot)
	if c32 >= currSlot.start && c32 < currSlot.end {
		return currSlot.node
	}

	currStep := currSlot.end - currSlot.start
	distance := (int64(c32) - int64(currSlot.start)) / int64(currStep)

	cursor := c.nodeRing.Move(int(distance))

	moveCnt := 0
	for {
		slot := cursor.Value.(*chashSlot)
		if c32 >= slot.start && c32 < slot.end {
			for {
				if slot.node.Health() {
					return slot.node
				}

				if moveCnt > c.nodeCount {
					return nil
				}

				cursor = cursor.Next()
				slot = cursor.Value.(*chashSlot)
				moveCnt++
			}
		}

		fmt.Println("no hit,skip ", slot.node)
		if c32 < slot.start {
			cursor = cursor.Prev()
		} else if c32 >= slot.end {
			cursor = cursor.Next()
		}

		moveCnt++
	}

	return nil
}
