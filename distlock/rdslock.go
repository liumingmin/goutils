package distlock

import (
	"net"
	"time"

	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/liumingmin/goutils/conf"
)

var gRdsPool *redis.Pool

func AquireLock(res string, timeout int) bool {
	c := gRdsPool.Get()
	defer c.Close()

	val, err := redis.Int(c.Do("EVAL",
		`if(redis.call("EXISTS",KEYS[1])==1)then if(redis.call("GET",KEYS[1])==ARGV[2])then redis.call("EXPIRE",KEYS[1],ARGV[1]);return 1;else return 0;end;else redis.call("SETEX",KEYS[1],ARGV[1],ARGV[2]);return 1;end`,
		1, res, timeout, gLocakKey))
	if err != nil {
		return false
	}

	return val == 1
}

var gLocakKey = func() string {
	localKey := os.Args[0]

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localKey
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localKey += "~" + ipnet.IP.String()
			}
		}
	}
	return localKey
}()

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
