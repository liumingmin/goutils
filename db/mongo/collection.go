package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewCompCollectionOp(client *Client, dbName, collectionName string) *CompCollectionOp {
	defaultC := NewMgoCollection(client, dbName, collectionName, nil, nil)

	pOpts := []*options.CollectionOptions{options.Collection().SetReadPreference(readpref.Primary())}
	primaryC := NewMgoCollection(client, dbName, collectionName, nil, pOpts)

	sOpts := []*options.CollectionOptions{options.Collection().SetReadPreference(readpref.SecondaryPreferred())}
	slaveC := NewMgoCollection(client, dbName, collectionName, nil, sOpts)

	return &CompCollectionOp{
		CollectionOp: defaultC,
		Primary:      primaryC,
		Slave:        slaveC,
	}
}

func NewMgoCollection(client *Client, dbName, collectionName string,
	dbOpts []*options.DatabaseOptions, cOpts []*options.CollectionOptions) *CollectionOp {
	if dbOpts == nil {
		dbOpts = make([]*options.DatabaseOptions, 0)
	}

	if cOpts == nil {
		cOpts = make([]*options.CollectionOptions, 0)
	}

	mgoC := &CollectionOp{
		client:   client,
		database: client.Database(dbName, dbOpts...),
	}
	mgoC.collection = mgoC.database.Collection(collectionName, cOpts...)
	return mgoC
}

// 复合集合操作类
type CompCollectionOp struct {
	*CollectionOp               //全局配置
	Primary       *CollectionOp //强制读主库
	Slave         *CollectionOp //优先读从库
}

// 集合操作类
type CollectionOp struct {
	client     *Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (c *CollectionOp) Insert(ctx context.Context, data interface{}) error {
	return c.client.Insert(ctx, c.collection, data)
}

func (c *CollectionOp) BatchInsert(ctx context.Context, data []interface{}) error {
	return c.client.BatchInsert(ctx, c.collection, data)
}

func (c *CollectionOp) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	return c.client.DeleteById(ctx, c.collection, id)
}

func (c *CollectionOp) Delete(ctx context.Context, selector interface{}) error {
	return c.client.Delete(ctx, c.collection, selector)
}

func (c *CollectionOp) Update(ctx context.Context, selector interface{}, updateOp interface{}) error {
	return c.client.Update(ctx, c.collection, selector, updateOp)
}

func (c *CollectionOp) UpdateById(ctx context.Context, id primitive.ObjectID, updateOp interface{}) error {
	return c.client.UpdateById(ctx, c.collection, id, updateOp)
}

func (c *CollectionOp) Count(ctx context.Context, query interface{}) (int64, error) {
	return c.client.Count(ctx, c.collection, query)
}

func (c *CollectionOp) FindById(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	return c.client.FindById(ctx, c.collection, id, result)
}

func (c *CollectionOp) Find(ctx context.Context, model FindModel, opts ...*options.FindOptions) error {
	return c.client.FindByModel(ctx, c.collection, model, opts...)
}

func (c *CollectionOp) FindOne(ctx context.Context, model FindModel, opts ...*options.FindOneOptions) error {
	return c.client.FindOneByModel(ctx, c.collection, model, opts...)
}

func (c *CollectionOp) FindOneAndUpdate(ctx context.Context, result interface{}, query interface{}, updateMap interface{},
	returnNew bool, upsert bool) error {
	return c.client.FindOneAndUpdate(ctx, c.collection, result, query, updateMap, returnNew, upsert)
}

func (c *CollectionOp) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error {
	return c.client.Aggregate(ctx, c.collection, pipeline, result)
}

// setOnInsertM only write on insert
func (c *CollectionOp) Upsert(ctx context.Context, selector interface{}, updateOp bson.M, setOnInsertItem interface{}) error {
	if setOnInsertItem != nil {
		updateOp["$setOnInsert"] = setOnInsertItem
	}
	return c.client.Update(ctx, c.collection, selector, updateOp, options.Update().SetUpsert(true))
}

func (c *CollectionOp) BulkUpdateItems(ctx context.Context, bulkUpdateItems []*BulkUpdateItem,
	opts ...*options.BulkWriteOptions) error {
	return c.client.BulkUpdateItems(ctx, c.collection, bulkUpdateItems, opts...)
}

func (c *CollectionOp) BulkUpsertItem(ctx context.Context, bulkUpertItems []*BulkUpsertItem,
	opts ...*options.BulkWriteOptions) error {
	return c.client.BulkUpsertItems(ctx, c.collection, bulkUpertItems, opts...)
}
