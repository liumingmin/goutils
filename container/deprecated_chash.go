package container

import (
	"container/ring"
	"fmt"
	"sync"

	"github.com/liumingmin/goutils/utils"
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
	start uint16
	end   uint16
	node  NodeHealth
}

func NewChash(nodes []NodeHealth) *Chash {
	initLen := len(nodes)
	r := ring.New(initLen)
	step := ^uint16(0) / uint16(initLen)

	next := r
	for i := 0; i < initLen; i++ {
		next = next.Next()
		slot := &chashSlot{start: uint16(i) * step, end: uint16(i+1) * step, node: nodes[i]}

		if i == initLen-1 {
			slot.end = ^uint16(0)
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
	cursor := c.nodeRing

	c16 := utils.Crc16([]byte(key))

	currSlot := cursor.Value.(*chashSlot)
	if c16 >= currSlot.start && c16 < currSlot.end {
		return currSlot.node
	}

	currStep := currSlot.end - currSlot.start
	distance := (int64(c16) - int64(currSlot.start)) / int64(currStep)

	cursor = cursor.Move(int(distance))

	moveCnt := 0
	for {
		slot := cursor.Value.(*chashSlot)
		if c16 >= slot.start && c16 < slot.end {
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
		if c16 < slot.start {
			cursor = cursor.Prev()
		} else if c16 >= slot.end {
			cursor = cursor.Next()
		}

		moveCnt++
	}

	return nil
}
