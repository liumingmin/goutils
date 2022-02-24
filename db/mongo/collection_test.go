package mongo

import (
	"context"
	"testing"

	"github.com/liumingmin/goutils/log"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCompCollection(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, err := MgoClient("goutils")
	log.Error(ctx, "err :%v", err)

	op := NewCompCollectionOp(c, "confmgr", "test")
	op.Insert(ctx, bson.M{"name": "test"})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
}
