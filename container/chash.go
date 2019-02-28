package container

import (
	"container/ring"
	"fmt"
	"hash/crc32"
)

type Chash struct {
	ring     *ring.Ring
	nodeRing *ring.Ring
	initLen  uint32
	initStep uint32
}

type chashSlot struct {
	start uint32
	end   uint32
	node  interface{}
}

func NewChash(nodes []interface{}) *Chash {
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
		nodeRing: r,
		ring:     r,
		initLen:  uint32(initLen),
		initStep: step,
	}
}

func (c *Chash) AddNode(node interface{}) {
	next := c.nodeRing.Next()

	slot := next.Value.(*chashSlot)
	newstep := (slot.end - slot.start) / 2

	r := ring.New(1)
	r.Value = &chashSlot{start: slot.start, end: slot.start + newstep, node: node}

	c.nodeRing.Link(r)
	slot.start = slot.start + newstep

	c.nodeRing = next
}

func (c *Chash) GetNode(key string) interface{} {
	c32 := crc32.ChecksumIEEE([]byte(key))
	distance := c32 / c.initStep
	cursor := c.ring.Move(int(distance + 1))

	for {
		slot := cursor.Value.(*chashSlot)
		if c32 >= slot.start && c32 < slot.end {
			return slot.node
		}

		fmt.Println("no hit,skip ", slot.node)
		cursor = cursor.Next()
	}

	return nil
}
