package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
)

func MqLsPush(ctx context.Context, dbKey string, key string, data ...interface{}) {
	client := Get(dbKey)
	client.LPush(ctx, key, data...)
}

//list结构主要用来做异步操作,相当于一对一发送接收
func MqLsPop(rds redis.UniversalClient, keys []string, timeout int, handler func(key, data string), goPoolSize int) {
	goPool := make(chan struct{}, goPoolSize)

	safego.Go(func() {
		//应对redis重启的情况
		for {
			ctx := context.Background()
			mqLsPopLoop(ctx, rds, keys, timeout, goPool, handler)
			time.Sleep(time.Second * 2)
		}
	})
}

func mqLsPopLoop(ctx context.Context, rds redis.UniversalClient, keys []string, timeout int, goPool chan struct{},
	handler func(key, data string)) {
	defer log.Recover(ctx, func(e interface{}) string {
		err, _ := e.(error)
		return fmt.Sprintf("mqLsDoPop failed. error: %v", err)
	})

	timerTick := time.Second * 5
	timer := time.NewTimer(timerTick)
	defer timer.Stop()

	log.Debug(ctx, "BRPop: %v", keys)
	for {
		kvReply, err := rds.BRPop(ctx, time.Second*time.Duration(timeout), keys...).Result()
		if err != nil && err != redis.Nil {
			log.Error(ctx, "mqLsDoPop err: %v", err)
			return
		}

		if len(kvReply) < 2 || len(kvReply[1]) == 0 {
			continue
		}

		timer.Reset(timerTick)

		select {
		case goPool <- struct{}{}:
		case <-timer.C:
		}

		safego.Go(func() {
			defer func() {
				select {
				case <-goPool:
				default:
				}
			}()
			log.Debug(ctx, "BRPop: %s", kvReply)
			handler(kvReply[0], kvReply[1]) //handler(listKey,value)
		})
	}
}
