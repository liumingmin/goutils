package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/log"
)

var (
	redisClients = make(map[string]redis.UniversalClient)
)

func InitRedises() {
	for _, redisItem := range conf.Conf.Redises {
		redisOpt := &redis.UniversalOptions{
			MasterName:    redisItem.MasterName,
			Addrs:         redisItem.Addrs,
			PoolSize:      redisItem.PoolSize,
			Password:      redisItem.Password,
			DB:            redisItem.Db,
			ReadOnly:      redisItem.ReadOnly,
			RouteRandomly: redisItem.RouteRandomly, //http://vearne.cc/archives/1113
			DialTimeout:   defDurationValue(redisItem.DialTimeout, 10*time.Second),
			ReadTimeout:   defDurationValue(redisItem.ReadTimeout, 10*time.Minute),
			WriteTimeout:  defDurationValue(redisItem.WriteTimeout, 10*time.Minute),
			IdleTimeout:   defDurationValue(redisItem.IdleTimeout, -1), // Default is 5 minutes. -1 disables idle timeout check.
			MaxRetries:    0,
			MaxRedirects:  -1,
		}
		client := redis.NewUniversalClient(redisOpt)

		log.Info(context.Background(), "redis client address: %v, db: %v, poolSize: %v", redisOpt.Addrs, redisOpt.DB, redisOpt.PoolSize)
		redisClients[redisItem.Key] = client
	}
}

func defDurationValue(value interface{}, defValue time.Duration) time.Duration {
	if value == nil {
		return defValue
	}

	dataValue, err := time.ParseDuration(fmt.Sprint(value))
	if err == nil {
		return dataValue
	}
	return defValue
}

func Get(key string) redis.UniversalClient {
	return redisClients[key]
}
