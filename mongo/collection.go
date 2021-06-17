package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMgoCollection(client *Client, dbName, collectionName string,
	dbOpts []*options.DatabaseOptions, cOpts []*options.CollectionOptions) *MgoCollection {
	if dbOpts == nil {
		dbOpts = make([]*options.DatabaseOptions, 0)
	}

	if cOpts == nil {
		cOpts = make([]*options.CollectionOptions, 0)
	}

	mgoC := &MgoCollection{
		client:   client,
		database: client.Client.Database(dbName, dbOpts...),
	}
	mgoC.collection = mgoC.database.Collection(collectionName, cOpts...)
	return mgoC
}

func NewWRMgoCollection(client *Client, dbName, collectionName string) *MgoCollection {
	opts := []*options.CollectionOptions{options.Collection().SetReadPreference(readpref.Primary())}
	return NewMgoCollection(client, dbName, collectionName, nil, opts)
}

func NewRDMgoCollection(client *Client, dbName, collectionName string) *MgoCollection {
	opts := []*options.CollectionOptions{options.Collection().SetReadPreference(readpref.SecondaryPreferred())}
	return NewMgoCollection(client, dbName, collectionName, nil, opts)
}

type MgoCollection struct {
	client     *Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (c *MgoCollection) Insert(ctx context.Context, data interface{}) error {
	return c.client.Insert(ctx, c.collection, data)
}

func (c *MgoCollection) BatchInsert(ctx context.Context, data []interface{}) error {
	return c.client.BatchInsert(ctx, c.collection, data)
}

func (c *MgoCollection) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	return c.client.DeleteById(ctx, c.collection, id)
}

func (c *MgoCollection) Delete(ctx context.Context, selector interface{}) error {
	return c.client.Delete(ctx, c.collection, selector)
}

func (c *MgoCollection) Update(ctx context.Context, selector interface{}, updateOp interface{}) error {
	return c.client.Update(ctx, c.collection, selector, updateOp)
}

func (c *MgoCollection) UpdateById(ctx context.Context, id primitive.ObjectID, updateOp interface{}) error {
	return c.client.UpdateById(ctx, c.collection, id, updateOp)
}

func (c *MgoCollection) Count(ctx context.Context, query interface{}) (int64, error) {
	return c.client.Count(ctx, c.collection, query)
}

func (c *MgoCollection) FindById(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	return c.client.FindById(ctx, c.collection, id, result)
}

func (c *MgoCollection) Find(ctx context.Context, results, query interface{}, sort []string, fields bson.M, cursor, size int) error {
	findModel := FindModel{
		Query:   query,
		Fields:  fields,
		Sort:    sort,
		Cursor:  cursor,
		Size:    size,
		Results: results,
	}
	return c.client.FindByModel(ctx, c.collection, findModel)
}

func (c *MgoCollection) FindOneAndUpdate(ctx context.Context, result interface{}, selector interface{}, updateMap interface{}, returnNew bool, upsert bool) error {
	return c.client.FindOneAndUpdate(ctx, c.collection, result, selector, updateMap, returnNew, upsert)
}

func (c *MgoCollection) FindOne(ctx context.Context, result, query interface{}, sort []string, fields bson.M, cursor int) error {
	findModel := FindModel{
		Query:   query,
		Fields:  fields,
		Sort:    sort,
		Cursor:  cursor,
		Results: result,
	}
	return c.client.FindOneByModel(ctx, c.collection, findModel)
}

func (c *MgoCollection) Aggregate(ctx context.Context, pipeline interface{}, result interface{}) error {
	return c.client.Aggregate(ctx, c.collection, pipeline, result)
}

func (c *MgoCollection) Upsert(ctx context.Context, selector interface{}, updateOp interface{}) error {
	return c.client.Update(ctx, c.collection, selector, updateOp, options.Update().SetUpsert(true))
}
