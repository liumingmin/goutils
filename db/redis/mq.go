package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
)

func MqPSubscribe(c context.Context, rds redis.UniversalClient, pChannel string, handler func(channel string, data string),
	goPoolSize int) {
	coroutineChannel := make(chan struct{}, goPoolSize)

	safego.Go(func() {
		//应对redis重启的情况
		for {
			ctx := context.Background()
			func() {
				defer log.Recover(ctx, func(e interface{}) string {
					err, _ := e.(error)
					return fmt.Sprintf("Subscribe failed. error: %v", err)
				})

				sub := rds.PSubscribe(ctx, pChannel)
				defer sub.PUnsubscribe(ctx, pChannel)

				log.Info(ctx, "Subscribe: %s", pChannel)

				timerTick := time.Second * 5
				timer := time.NewTimer(timerTick)
				defer timer.Stop()

				for {
					msg, err := sub.ReceiveMessage(ctx)
					if err != nil {
						log.Error(ctx, "ReceiveMessage err: %v", err)
						return
					}

					timer.Reset(timerTick)

					select {
					case coroutineChannel <- struct{}{}:
					case <-timer.C:
					}

					safego.Go(func() {
						defer func() {
							select {
							case <-coroutineChannel:
							default:
							}
						}()
						handler(msg.Channel, msg.Payload)
					})
				}
			}()

			time.Sleep(time.Second * 2)
		}
	})
}

func MqPublish(ctx context.Context, rds redis.UniversalClient, channel string, message interface{}) error {
	return rds.Publish(ctx, channel, message).Err()
}
