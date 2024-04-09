package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
)

func ListPush(ctx context.Context, rds redis.UniversalClient, key string, data ...interface{}) error {
	return rds.LPush(ctx, key, data...).Err()
}

// list结构主要用来做异步操作,相当于一对一发送接收
func ListPop(rds redis.UniversalClient, keys []string, timeout, goPoolSize int, handler func(key, data string)) {
	goPool := make(chan struct{}, goPoolSize)

	safego.Go(func() {
		//应对redis重启的情况
		for {
			ctx := context.Background()
			listPopLoop(ctx, rds, keys, timeout, goPool, handler)
			time.Sleep(time.Second * 2)
		}
	})
}

func listPopLoop(ctx context.Context, rds redis.UniversalClient, keys []string, timeout int, goPool chan struct{},
	handler func(key, data string)) error {
	defer log.Recover(ctx, func(e interface{}) string {
		err, _ := e.(error)
		return fmt.Sprintf("listPopLoop failed. error: %v", err)
	})

	timerTick := time.Second * 5
	timer := time.NewTimer(timerTick)
	defer timer.Stop()

	log.Debug(ctx, "listPopLoop: %v", keys)
	for {
		kvReply, err := rds.BRPop(ctx, time.Second*time.Duration(timeout), keys...).Result()
		if err != nil && err != redis.Nil {
			log.Error(ctx, "listPopLoop err: %v", err)
			return err
		}

		if len(kvReply) < 2 || len(kvReply[1]) == 0 {
			continue
		}

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
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
			handler(kvReply[0], kvReply[1]) //handler(listKey,value)
		})
	}
}
