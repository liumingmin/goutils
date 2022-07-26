package ws

import (
	"context"
	"sync/atomic"

	"github.com/liumingmin/goutils/log"
)

type defaultSrvComet struct {
	pullChannelId int
	firstPullFunc func(context.Context, *Connection) // first connected exec
	pullFunc      func(context.Context, *Connection) // every times exec
	isRunning     int32
}

func (c *defaultSrvComet) PullSend(conn *Connection) {
	ok := atomic.CompareAndSwapInt32(&c.isRunning, 0, 1)
	if !ok {
		log.Debug(context.Background(), "comet is running, pullChannelId: %v", c.pullChannelId)
		return
	}
	defer atomic.StoreInt32(&c.isRunning, 0)

	pullChannel, ok := conn.GetPullChannel(c.pullChannelId)
	if !ok {
		return
	}

	if c.firstPullFunc != nil {
		c.firstPullFunc(context.Background(), conn)
	}

	for {
		ctx := context.Background()

		if conn.IsStopped() {
			log.Debug(ctx, "agent is stopped: %v", conn.Id())
			return
		}

		c.pullFunc(ctx, conn)

		if _, ok := <-pullChannel; !ok {
			log.Debug(ctx, "Connect stop pull channel. connId: %v, channelId: %v", conn.Id(), c.pullChannelId)
			return
		}
	}
}
