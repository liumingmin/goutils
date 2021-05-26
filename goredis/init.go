package goredis

import (
	"context"
	"fmt"
	"time"

	"github.com/demdxx/gocast"
	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/log"
)

var (
	goRedisClients = make(map[string]redis.UniversalClient)
)

func init() {
	gorediss := conf.Ext("goredis", []interface{}{})
	goredissArray, ok := gorediss.([]interface{})
	if !ok {
		log.Info(context.Background(), "goredis conf not exists")
		return
	}

	for _, goredisc := range goredissArray {
		gorediscMap, ok := goredisc.(map[interface{}]interface{})
		if !ok {
			continue
		}

		key, ok := gorediscMap["key"]
		if !ok {
			continue
		}

		addrsObj, ok := gorediscMap["addrs"]
		if !ok {
			continue
		}
		db, _ := gorediscMap["db"]
		poolSize, _ := gorediscMap["pool_size"]
		password, _ := gorediscMap["password"]
		readTimeout, _ := gorediscMap["read_timeout"]
		writeTimeout, _ := gorediscMap["write_timeout"]
		idleTimeout, _ := gorediscMap["idle_timeout"]

		redisOpt := &redis.UniversalOptions{
			Addrs:         gocast.ToStringSlice(addrsObj),
			PoolSize:      gocast.ToInt(poolSize),
			Password:      gocast.ToString(password),
			DB:            gocast.ToInt(db),
			ReadOnly:      true,
			RouteRandomly: true, //http://vearne.cc/archives/1113
			DialTimeout:   10 * time.Second,
			ReadTimeout:   defDurationValue(readTimeout, 10*time.Minute),
			WriteTimeout:  defDurationValue(writeTimeout, 10*time.Minute),
			IdleTimeout:   defDurationValue(idleTimeout, -1), // Default is 5 minutes. -1 disables idle timeout check.
			MaxRetries:    0,
			MaxRedirects:  -1,
		}
		client := redis.NewUniversalClient(redisOpt)

		log.Info(context.Background(), "redis client address: %v, db: %v, poolSize: %v", redisOpt.Addrs, redisOpt.DB, redisOpt.PoolSize)
		goRedisClients[gocast.ToString(key)] = client
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
	return goRedisClients[key]
}
