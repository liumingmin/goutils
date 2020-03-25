package distlock

import (
	"time"

	"goutils/conf"
	"goutils/safego"

	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
)

var gRdsPool *redis.Pool

type RdsLuaLock struct {
	key    string
	value  string
	expire int
}

func NewRdsLuaLock(key string, expire int) (*RdsLuaLock, error) {
	return &RdsLuaLock{key: key, value: uuid.New().String(), expire: expire}, nil
}

func (l *RdsLuaLock) TryLock() bool {
	c := gRdsPool.Get()
	defer c.Close()

	val, err := redis.Int(c.Do("EVAL",
		`if(redis.call("EXISTS",KEYS[1])==1)then if(redis.call("GET",KEYS[1])==ARGV[2])then redis.call("EXPIRE",KEYS[1],ARGV[1]);return 1;else return 0;end;else redis.call("SETEX",KEYS[1],ARGV[1],ARGV[2]);return 1;end`,
		1, l.key, l.expire, l.value))
	if err != nil {
		return false
	}

	return val == 1
}

func (l *RdsLuaLock) Lock(timeout int) bool {
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
			result := l.TryLock()
			if result {
				return true
			}
		}
	}
}

func (l *RdsLuaLock) Unlock() {
	c := gRdsPool.Get()
	defer c.Close()

	c.Do("DEL", l.key)
}

func init() {
	gRdsPool = &redis.Pool{
		Wait:        true,
		MaxActive:   conf.ExtInt("redis.ext.maxActive", 1024),
		MaxIdle:     conf.ExtInt("redis.ext.maxIdle", 32),
		IdleTimeout: conf.ExtDuration("redis.ext.idleTimeout"),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.ExtString("redis.host"),
				redis.DialConnectTimeout(conf.ExtDuration("redis.ext.connectTimeout", "5s")))
			if err != nil {
				return c, err
			}

			pwd := conf.ExtString("redis.password", "")
			if pwd != "" {
				err = c.Send("AUTH")
				if err != nil {
					return c, err
				}
			}

			_, err = c.Do("SELECT", conf.ExtInt("redis.ext.distlockdb", 0))
			if err != nil {
				return c, err
			}
			if event := conf.ExtString("redis.ext.notifyKeyspaceEvents", ""); event != "" {
				_, err := c.Do("CONFIG", "SET", "notify-keyspace-events", event)
				if err != nil {
					return c, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do(conf.ExtString("redis.ext.testOnBorrow"))
			return err
		},
	}
}
