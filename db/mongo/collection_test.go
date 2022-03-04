package mongo

import (
	"context"
	"testing"

	"github.com/liumingmin/goutils/log"
	"go.mongodb.org/mongo-driver/bson"
)

const dbKey = "testDbKey"
const dbName = "testDb"
const collectionName = "testUser"

type testUser struct {
	UserId   string `bson:"user_id"`
	Nickname string `bson:"nick_name"`
	Status   string `bson:"status"`
	Type     string `bson:"p_type"`
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Insert(ctx, testUser{
		UserId:   "1",
		Nickname: "超级棒",
		Status:   "valid",
		Type:     "normal",
	})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Update(ctx, bson.M{"user_id": "1"}, bson.M{"$set": bson.M{"nick_name": "超级棒++"}})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	op.Delete(ctx, bson.M{"user_id": "1"})

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
}
