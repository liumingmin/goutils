package distlock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/liumingmin/goutils/redis"
	"github.com/liumingmin/goutils/safego"
)

type RdsLuaLock struct {
	dbKey  string
	key    string
	value  string
	expire int
}

func NewRdsLuaLock(dbKey, key string, expire int) (*RdsLuaLock, error) {
	return &RdsLuaLock{dbKey: dbKey, key: key, value: uuid.New().String(), expire: expire}, nil
}

func (l *RdsLuaLock) TryLock(ctx context.Context) bool {
	rds := redis.Get(l.dbKey)

	val, err := rds.Eval(ctx, `if(redis.call("EXISTS",KEYS[1])==1)then if(redis.call("GET",KEYS[1])==ARGV[2])then redis.call("EXPIRE",KEYS[1],ARGV[1]);return 1;else return 0;end;else redis.call("SETEX",KEYS[1],ARGV[1],ARGV[2]);return 1;end`,
		[]string{l.key}, l.expire, l.value).Int()
	if err != nil {
		return false
	}

	return val == 1
}

func (l *RdsLuaLock) Lock(ctx context.Context, timeout int) bool {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	stopChan := make(chan struct{})
	safego.Go(func() {
		time.Sleep(time.Second * time.Duration(timeout))
		stopChan <- struct{}{}
	})

	for {
		select {
		case <-stopChan:
			return false
		case <-t.C:
			result := l.TryLock(ctx)
			if result {
				return true
			}
		}
	}
}

func (l *RdsLuaLock) Unlock(ctx context.Context) {
	rds := redis.Get(l.dbKey)

	rds.Del(ctx, l.key)
}
