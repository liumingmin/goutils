package distlock

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/liumingmin/goutils/log"
)

var gConsul *api.Client

type ConsulLock struct {
	key  string
	lock *api.Lock
}

func NewConsulLock(key string, expire int) (*ConsulLock, error) {
	opts := &api.LockOptions{
		Key:        key,
		SessionTTL: "20s", //10s ~ 24h
	}

	l, err := gConsul.LockOpts(opts)
	if err != nil {
		return nil, err
	}

	clock := &ConsulLock{key: key, lock: l}
	//safego.Go(func() {
	//	time.Sleep(time.Second * time.Duration(expire))
	//	clock.Unlock()
	//})
	return clock, nil
}

func (l *ConsulLock) Lock(ctx context.Context, timeout int) bool {
	stopChan := make(chan struct{})
	go func() {
		defer log.Recover(ctx, func(e interface{}) string {
			return fmt.Sprintf("ConsulLock Lock err: %v", e)
		})

		time.Sleep(time.Second * time.Duration(timeout))
		stopChan <- struct{}{}
	}()

	ldChan, err := l.lock.Lock(stopChan)
	//fmt.Println(time.Now())
	return ldChan != nil && err == nil
}

func (l *ConsulLock) Unlock() {
	l.lock.Unlock()
}

func InitConsul(centerAddr string) {
	if centerAddr == "" {
		centerAddr = "127.0.0.1:8500"
	}

	var e error
	config := api.DefaultConfig()
	config.Address = centerAddr
	gConsul, e = api.NewClient(config)
	if e != nil {
		log.Error(context.Background(), "Create consul client failed. error: %v", e)
		panic(e)
	}
}
