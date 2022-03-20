package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liumingmin/goutils/algorithm"
)

type ZDescartesFilter func([]string) (string, map[string]int64)

func ZDescartes(ctx context.Context, rds redis.UniversalClient, dimValues [][]string, filter ZDescartesFilter,
	ttl, batchSize int) error {

	combinations := algorithm.DescartesCombine(dimValues)
	for _, combination := range combinations {
		keyName, items := filter(combination)

		members := make([]*redis.Z, 0, len(items))
		for itemKey, itemValue := range items {
			members = append(members, &redis.Z{
				Score:  float64(itemValue),
				Member: itemKey,
			})
		}
		err := ZBatchAdd(ctx, rds, keyName, members, ttl, batchSize)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZBatchAdd(ctx context.Context, rds redis.UniversalClient, keyName string, members []*redis.Z, ttl, batchSize int) error {
	defer func() {
		rds.Expire(ctx, keyName, time.Duration(ttl)*time.Second)
	}()

	batchMembers := make([]*redis.Z, 0, batchSize)
	for i := 0; i < len(members); i++ {
		batchMembers = append(batchMembers, members[i])

		if len(batchMembers) == batchSize || i == len(members)-1 {
			err := rds.ZAdd(ctx, keyName, batchMembers...).Err()
			if err != nil {
				return err
			}
			batchMembers = batchMembers[0:0]
		}
	}
	return nil
}
