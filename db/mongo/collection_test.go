package mongo

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const testDbKey = "testDbKey"
const testDbName = "testDb"
const testCollectionName = "testUser"

type testUser struct {
	UserId   string `bson:"user_id"`
	Nickname string `bson:"nick_name"`
	Status   string `bson:"status"`
	Type     string `bson:"p_type"`
}

var testDbClient *Client
var testMgoDoc = testUser{
	UserId:   "1",
	Nickname: "超级棒",
	Status:   "valid",
	Type:     "normal",
}

func TestInsert(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)
	err := op.Insert(ctx, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	var result []testUser
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &result,
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("not found row")
	}

	if !reflect.DeepEqual(result[0], testMgoDoc) {
		t.Error("row data not equal")
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)

	err := op.Insert(ctx, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	err = op.Update(ctx, bson.M{"user_id": testMgoDoc.UserId}, bson.M{"$set": bson.M{"nick_name": "超级棒++"}})
	if err != nil {
		t.Error(err)
	}

	testMgoDoc.Nickname = "超级棒++"

	var result []testUser
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &result,
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("not found row")
	}

	if !reflect.DeepEqual(result[0], testMgoDoc) {
		t.Error("row data not equal")
	}

	var oneTestUser testUser
	err = op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &oneTestUser,
	})
	if err == mongo.ErrNoDocuments {
		t.Error(err)
	}

	if !reflect.DeepEqual(oneTestUser, testMgoDoc) {
		t.Error("row data not equal")
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)

	testMgoDoc.UserId = "1234567"
	err := op.Insert(ctx, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	err = op.Delete(ctx, bson.M{"user_id": testMgoDoc.UserId})
	if err != nil {
		t.Error(err)
	}

	var result []bson.M
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &result,
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) != 0 {
		t.Error("not delete success")
	}
}

func TestUpert(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)

	testMgoDoc.UserId = "abcdefg"
	err := op.Upsert(ctx, bson.M{"user_id": testMgoDoc.UserId}, bson.M{"$set": bson.M{"birth": "2018"}}, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	var result []testUser
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &result,
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) == 0 {
		t.Error("not found row")
	}

	if !reflect.DeepEqual(result[0], testMgoDoc) {
		t.Error("row data not equal")
	}

	err = op.Upsert(ctx, bson.M{"user_id": testMgoDoc.UserId}, bson.M{"$set": bson.M{"birth": "2018"}}, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	var oneTestUser bson.M
	err = op.FindOne(ctx, FindModel{
		Query:   bson.M{"user_id": testMgoDoc.UserId},
		Results: &oneTestUser,
	})
	if err == mongo.ErrNoDocuments {
		t.Error(err)
	}

	if oneTestUser["birth"] != "2018" {
		t.Error("field birth is not equal")
	}
}

func TestBulkUpdateItems(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)

	u1 := "000001"
	u2 := "000002"

	f1 := "1"
	f2 := "2"

	testMgoDoc.UserId = u1
	err := op.Insert(ctx, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	testMgoDoc.UserId = u2
	err = op.Insert(ctx, testMgoDoc)
	if err != nil {
		t.Error(err)
	}

	err = op.BulkUpdateItems(ctx, []*BulkUpdateItem{
		{Selector: bson.M{"user_id": u1}, Update: bson.M{"$set": bson.M{"birth": f1}}},
		{Selector: bson.M{"user_id": u2}, Update: bson.M{"$set": bson.M{"birth2": f2}}},
	})
	if err != nil {
		t.Error(err)
	}

	var result []bson.M
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": bson.M{"$in": []string{u1, u2}}},
		Results: &result,
		Sort:    []string{"+user_id"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) < 2 {
		t.Error("row less 2")
	}

	if !reflect.DeepEqual(result[0]["birth"], f1) {
		t.Error("row1 data not equal")
	}

	if !reflect.DeepEqual(result[1]["birth2"], f2) {
		t.Error("row2 data not equal")
	}
}

func TestBulkUpsertItems(t *testing.T) {
	ctx := context.Background()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)

	u1 := "tim1"
	u2 := "tim2"

	f1 := "1"
	f2 := "2"
	var err error

	err = op.BulkUpsertItem(ctx, []*BulkUpsertItem{
		{Selector: bson.M{"user_id": u1}, Replacement: bson.M{"user_id": u1, "birth": f1}},
		{Selector: bson.M{"user_id": u2}, Replacement: bson.M{"user_id": u2, "birth2": f2}},
	})
	if err != nil {
		t.Error(err)
	}

	var result []bson.M
	err = op.Find(ctx, FindModel{
		Query:   bson.M{"user_id": bson.M{"$in": []string{u1, u2}}},
		Results: &result,
		Sort:    []string{"+user_id"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(result) < 2 {
		t.Error("row less 2")
	}

	if !reflect.DeepEqual(result[0]["birth"], f1) {
		t.Error("row1 data not equal")
	}

	if !reflect.DeepEqual(result[1]["birth2"], f2) {
		t.Error("row2 data not equal")
	}
}

func TestMain(m *testing.M) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:27017", time.Second*2)
	if err != nil {
		fmt.Println("Please install mongo on local and start at port: 27017, then run test.")
		return
	}
	conn.Close()

	InitClients()

	testDbClient, err = MgoClient(testDbKey)
	if err != nil {
		fmt.Println("Please config mongo yml section testDbKey to local address and port: '127.0.0.1:27017', then run test.")
		return
	}

	m.Run()

	op := NewCompCollectionOp(testDbClient, testDbName, testCollectionName)
	op.collection.Drop(context.Background())
}
