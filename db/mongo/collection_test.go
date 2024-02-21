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
	err := op.Insert(ctx, testUser{
		UserId:   "1",
		Nickname: "超级棒",
		Status:   "valid",
		Type:     "normal",
	})
	log.Info(context.Background(), err)

	var result []bson.M
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})
	if err != nil {
		log.Error(ctx, "Mgo find err: %v", err)
		return
	}
	for _, item := range result {
		t.Log(item)
	}
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

func TestFind(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)

	var result []bson.M
	err := op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})
	if err != nil {
		log.Error(ctx, "Mgo find err: %v", err)
		return
	}

	for _, item := range result {
		t.Log(item)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	err := op.Delete(ctx, bson.M{"user_id": "1"})
	log.Info(context.Background(), err)

	var result interface{}
	op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": "1"},
		Results: &result,
	})

	log.Info(ctx, "result: %v", result)
}

func TestUpert(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)
	err := op.Upsert(ctx, bson.M{"name": "tom2"}, bson.M{"$set": bson.M{"birth": "2020"}}, bson.M{"birth2": "2024"})
	t.Log(err)
}

func TestBulkUpdateItems(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)

	err := op.BulkUpdateItems(ctx, []*BulkUpdateItem{
		{Selector: bson.M{"name": "tom"}, Update: bson.M{"$set": bson.M{"birth": "1"}}},
		{Selector: bson.M{"name": "tom1"}, Update: bson.M{"$set": bson.M{"birth2": "2"}}},
	})
	t.Log(err)
}

func TestBulkUpsertItems(t *testing.T) {
	ctx := context.Background()
	InitClients()
	c, _ := MgoClient(dbKey)

	op := NewCompCollectionOp(c, dbName, collectionName)

	err := op.BulkUpsertItem(ctx, []*BulkUpsertItem{
		{Selector: bson.M{"name": "tim"}, Replacement: bson.M{"name": "tim", "birth": "3"}},
		{Selector: bson.M{"name": "tim1"}, Replacement: bson.M{"name": "tim1", "birth2": "4"}},
	})
	t.Log(err)
}
