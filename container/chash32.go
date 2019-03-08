package container

import (
	"container/ring"
	"fmt"
	"hash/crc32"
	"sync"
)

type Chash32 struct {
	ring      *ring.Ring
	nodeRing  *ring.Ring
	nodeCount int
	lock      sync.RWMutex
}

type chashSlot32 struct {
	start uint32
	end   uint32
	node  NodeHealth
}

func NewChash32(nodes []NodeHealth) *Chash32 {
	initLen := len(nodes)
	r := ring.New(initLen)
	step := ^uint32(0) / uint32(initLen)

	next := r
	for i := 0; i < initLen; i++ {
		next = next.Next()
		slot := &chashSlot32{start: uint32(i) * step, end: uint32(i+1) * step, node: nodes[i]}

		if i == initLen-1 {
			slot.end = ^uint32(0)
		}

		next.Value = slot
	}

	return &Chash32{
		ring:      r,
		nodeRing:  r,
		nodeCount: initLen,
	}
}

func (c *Chash32) AddNode(node NodeHealth) {
	c.lock.Lock()
	defer c.lock.Unlock()

	next := c.ring.Next()

	slot := next.Value.(*chashSlot32)
	newstep := (slot.end - slot.start) / 2

	r := ring.New(1)
	r.Value = &chashSlot32{start: slot.start, end: slot.start + newstep, node: node}

	c.ring.Link(r)
	c.nodeCount++
	slot.start = slot.start + newstep
	c.ring = next
}

func (c *Chash32) GetNode(key string) NodeHealth {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cursor := c.nodeRing

	c32 := crc32.ChecksumIEEE([]byte(key))

	currSlot := cursor.Value.(*chashSlot32)
	if c32 >= currSlot.start && c32 < currSlot.end {
		return currSlot.node
	}

	currStep := currSlot.end - currSlot.start
	distance := (int64(c32) - int64(currSlot.start)) / int64(currStep)

	cursor = cursor.Move(int(distance))

	moveCnt := 0
	for {
		slot := cursor.Value.(*chashSlot32)
		if c32 >= slot.start && c32 < slot.end {
			for {
				if slot.node.Health() {
					return slot.node
				}

				if moveCnt > c.nodeCount {
					return nil
				}

				cursor = cursor.Next()
				slot = cursor.Value.(*chashSlot32)
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
