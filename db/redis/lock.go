package redis

import (
	"context"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/log"
)

func RdsAllowActionWithCD(ctx context.Context, rds redis.UniversalClient, actionKey string, cdSeconds int) (int, bool) {
	if cdSeconds <= 0 {
		return 0, true
	}

	leftTime, err := rds.TTL(ctx, actionKey).Result()
	if err != nil {
		return cdSeconds, false
	}

	if leftTime > 0 {
		return int(leftTime / time.Second), false
	}

	ok, err := rds.SetNX(ctx, actionKey, 1, time.Duration(cdSeconds)*time.Second).Result()
	if err != nil {
		log.Error(ctx, "RdsAllowActionWithCD failed. key: %s, err: %v", actionKey, err)
		return cdSeconds, false
	}

	log.Debug(ctx, "RdsAllowActionWithCD result: %v", ok)
	if ok {
		return cdSeconds, true
	}

	return cdSeconds, false
}

func RdsAllowActionByMTs(ctx context.Context, rds redis.UniversalClient, actionKey string, cdMillSeconds, keyTTL int) (int64, bool) {
	if cdMillSeconds <= 0 {
		return 0, true
	}

	lastMts, _ := rds.Get(ctx, actionKey).Int64()
	nowMts := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	leftMillSeconds := nowMts - lastMts
	if int64(math.Abs(float64(leftMillSeconds))) < int64(cdMillSeconds) {
		return leftMillSeconds, false
	}

	err := rds.Set(ctx, actionKey, nowMts, time.Second*time.Duration(keyTTL)).Err()
	if err != nil {
		log.Error(ctx, "RdsAllowActionByMTs failed. key: %s, err: %v", actionKey, err)
	}

	log.Debug(ctx, "RdsAllowActionByMTs result: true, %v", leftMillSeconds)
	return int64(cdMillSeconds), true
}

//RdsLockResWithCD(ctx context.Context, rds redis.UniversalClient, resKeyName, config.ServerId, runtimeSeconds*2)
func RdsLockResWithCD(ctx context.Context, rds redis.UniversalClient, resKey, resValue string, cdSeconds int) bool {
	origResValue, err := rds.Get(ctx, resKey).Result()
	if err == nil {
		if origResValue == resValue {
			rds.Expire(ctx, resKey, time.Second*time.Duration(cdSeconds)) //资源续期
			origResValue, err = rds.Get(ctx, resKey).Result()
			if err == nil && origResValue == resValue {
				return true
			}
		}
		return false
	}

	ok, err := rds.SetNX(ctx, resKey, resValue, time.Second*time.Duration(cdSeconds)).Result()
	if err != nil {
		log.Error(ctx, "RdsLockResWithCD failed. key: %s, err: %v", resKey, err)
		return false
	}

	log.Debug(ctx, "RdsLockResWithCD result: %v", ok)
	return ok
}
