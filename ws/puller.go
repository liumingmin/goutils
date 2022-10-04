package ws

import (
	"context"
	"sync/atomic"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
)

type Puller interface {
	PullSend()
}

func NewDefaultPuller(conn IConnection, pullChannelId int, firstPullFunc, pullFunc func(context.Context, IConnection)) Puller {
	return &defaultPuller{
		conn:          conn,
		pullChannelId: pullChannelId,
		firstPullFunc: firstPullFunc,
		pullFunc:      pullFunc,
	}
}

type defaultPuller struct {
	conn          IConnection
	pullChannelId int
	firstPullFunc func(context.Context, IConnection) // first connected exec
	pullFunc      func(context.Context, IConnection) // every times exec
	isRunning     int32
}

func (c *defaultPuller) PullSend() {
	ctx := utils.ContextWithTsTrace()

	ok := atomic.CompareAndSwapInt32(&c.isRunning, 0, 1)
	if !ok {
		log.Debug(ctx, "comet is running, pullChannelId: %v", c.pullChannelId)
		return
	}
	defer atomic.StoreInt32(&c.isRunning, 0)

	pullChannel, ok := c.conn.GetPullChannel(c.pullChannelId)
	if !ok {
		return
	}

	if c.firstPullFunc != nil {
		c.firstPullFunc(ctx, c.conn)
	}

	for {
		if c.conn.IsStopped() {
			log.Info(ctx, "agent is stopped: %v", c.conn.Id())
			return
		}

		c.pullFunc(ctx, c.conn)

		if _, ok := <-pullChannel; !ok {
			log.Debug(ctx, "Connect stop pull channel. connId: %v, channelId: %v", c.conn.Id(), c.pullChannelId)
			return
		}

		ctx = utils.ContextWithTsTrace()
	}
}
