package ws

import (
	"context"
	"sync/atomic"

	"github.com/liumingmin/goutils/log"
)

type defaultComet struct {
	conn          *Connection
	pullChannelId int
	firstPullFunc func(context.Context, *Connection) // first connected exec
	pullFunc      func(context.Context, *Connection) // every times exec
	isRunning     int32
}

func (c *defaultComet) Pull() {
	ok := atomic.CompareAndSwapInt32(&c.isRunning, 0, 1)
	if !ok {
		log.Debug(context.Background(), "comet is running, pullChannelId: %v", c.pullChannelId)
		return
	}
	defer atomic.StoreInt32(&c.isRunning, 0)

	pullChannel, ok := c.conn.GetPullChannel(c.pullChannelId)
	if !ok {
		return
	}

	if c.firstPullFunc != nil {
		c.firstPullFunc(context.Background(), c.conn)
	}

	for {
		ctx := context.Background()

		if c.conn.IsStopped() {
			log.Debug(ctx, "agent is stopped: %v", c.conn.Id())
			return
		}

		c.pullFunc(ctx, c.conn)

		if _, ok := <-pullChannel; !ok {
			log.Debug(ctx, "Connect stop pull channel. connId: %v, channelId: %v", c.conn.Id(), c.pullChannelId)
			return
		}
	}
}
