package mongo

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Client struct {
	*mongo.Client
}

var mgoClients = make(map[string]*Client, 0)

func MgoClient(key string) (mongo *Client, err error) {
	if v, ok := mgoClients[key]; ok {
		return v, nil
	} else {
		log.Error(context.Background(), "Client not exist. key: %v", key)
		return nil, errors.New("Client not exist")
	}
}

func InitClients() {
	dbs := conf.Conf.Databases
	if dbs == nil {
		fmt.Fprintf(os.Stderr, "No database configuration")
		return
	}

	for _, db := range dbs {
		if db.Type == "mongo" {
			client, err := initClient(db)
			if err != nil {
				continue
			}

			mgoClients[db.Key] = client
		}
	}
}

func initClient(dbconf *conf.Database) (ret *Client, err error) {
	return NewClient(newConfig(dbconf))
}

func NewClient(opt *Config) (ret *Client, err error) {
	opts := options.Client()
	opts.SetHosts(opt.Address)
	if opt.Username != "" {
		var auth options.Credential
		auth.Username = opt.Username
		auth.Password = opt.Password
		auth.PasswordSet = true
		if opt.Source != "" {
			auth.AuthSource = opt.Source
		} else {
			auth.AuthSource = opt.Database
		}
		opts.SetAuth(auth)
	}
	if len(opt.Compressors) > 0 {
		opts.SetCompressors(opt.Compressors)
	}

	if opt.Keepalive > 0 {
		var dialer net.Dialer
		dialer.KeepAlive = opt.Keepalive
		if opt.ConnectTimeout > 0 {
			dialer.Timeout = opt.ConnectTimeout
		}
		opts.SetDialer(&dialer)
	} else if opt.ConnectTimeout > 0 {
		opts.SetConnectTimeout(opt.ConnectTimeout)
	}

	opts.SetDirect(opt.Direct)

	socketTimeout := opt.ReadTimeout
	if socketTimeout < opt.WriteTimeout {
		socketTimeout = opt.WriteTimeout
	}
	if socketTimeout > 0 {
		opts.SetSocketTimeout(socketTimeout)
	}

	if opt.MaxPoolIdleTimeMS > 0 {
		opts.SetMaxConnIdleTime(time.Duration(opt.MaxPoolIdleTimeMS) * time.Millisecond)
	}
	if opt.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(uint64(opt.MaxPoolSize))
	}
	if opt.MinPoolSize > 0 {
		opts.SetMinPoolSize(uint64(opt.MinPoolSize))
	}
	if opt.MaxPoolWaitTimeMS > 0 {
		opts.SetServerSelectionTimeout(time.Duration(opt.MaxPoolWaitTimeMS) * time.Millisecond)
	}

	if opt.Mode > 0 {
		var pref *readpref.ReadPref
		if pref, err = readpref.New(opt.Mode); err != nil {
			return
		}
		opts.SetReadPreference(pref)
	}

	if opt.Safe != nil {

		if opt.Safe.RMode != "" {
			opts.SetReadConcern(readconcern.New(readconcern.Level(strings.ToLower(opt.Safe.RMode))))
		}

		var ops []writeconcern.Option
		if opt.Safe.J {
			ops = append(ops, writeconcern.J(true))
		}
		if opt.Safe.W > 0 {
			ops = append(ops, writeconcern.W(opt.Safe.W))
		}
		if strings.ToLower(opt.Safe.WMode) == "majority" {
			ops = append(ops, writeconcern.WMajority())
		}
		if opt.Safe.WTimeout > 0 {
			ops = append(ops, writeconcern.WTimeout(time.Duration(opt.Safe.WTimeout)*time.Millisecond))
		}
		opts.SetWriteConcern(writeconcern.New(ops...))
	}

	useLocalTimeZone := true
	tc := bsoncodec.NewTimeCodec(&bsonoptions.TimeCodecOptions{UseLocalTimeZone: &useLocalTimeZone})
	r := bson.NewRegistryBuilder().RegisterTypeDecoder(reflect.TypeOf(time.Time{}), tc).Build()
	opts.SetRegistry(r)

	log.Info(context.Background(), "mgo official readPreference: %s, writeConcern: %#v",
		opts.ReadPreference.Mode(), opts.WriteConcern)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return
	}
	ret = &Client{
		Client: client,
	}
	return
}

//需要mongodb4.2以上才能支持事务
func (client *Client) ExecTx(ctx context.Context, f func(context.Context)) error {
	return client.UseSession(ctx, func(sctx mongo.SessionContext) (err error) {
		err = sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		defer func() {
			if err1 := recover(); err1 != nil {
				sctx.AbortTransaction(sctx)

				var ok bool
				if err, ok = err1.(error); !ok {
					err = errors.New(fmt.Sprint(err1))
				}

				log.Error(sctx, "mgo transaction failed, err: %v", err)
				return
			} else {
				for i := 0; i < 3; i++ {
					err = sctx.CommitTransaction(sctx)
					switch e := err.(type) {
					case nil:
						return
					case mongo.CommandError:
						if e.HasErrorLabel("UnknownTransactionCommitResult") {
							log.Warn(sctx, "UnknownTransactionCommitResult, retrying %v commit operation...", i)
							continue
						}
						log.Error(sctx, "mgo transaction failed, err: %v", e)
						return
					default:
						log.Error(sctx, "mgo transaction failed, err: %v", e)
						return
					}
				}
			}
		}()

		f(sctx)
		return nil
	})

}

//write
func (client *Client) Insert(ctx context.Context, collection *mongo.Collection, data interface{},
	opts ...*options.InsertOneOptions) error {

	_, err := collection.InsertOne(ctx, data, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) BatchInsert(ctx context.Context, collection *mongo.Collection, data []interface{},
	opts ...*options.InsertManyOptions) error {

	_, err := collection.InsertMany(ctx, data, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteById(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID,
	opts ...*options.DeleteOptions) error {

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id}, opts...)
	return err
}

func (client *Client) Delete(ctx context.Context, collection *mongo.Collection, selector interface{},
	opts ...*options.DeleteOptions) error {

	_, err := collection.DeleteMany(ctx, selector, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateById(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, update interface{},
	opts ...*options.UpdateOptions) error {
	return client.Update(ctx, collection, bson.M{"_id": id}, update, opts...)
}

func (client *Client) UpdateByIds(ctx context.Context, collection *mongo.Collection, ids []primitive.ObjectID, update interface{},
	opts ...*options.UpdateOptions) error {
	return client.Update(ctx, collection, bson.M{"_id": bson.M{"$in": ids}}, update, opts...)
}

func (client *Client) Update(ctx context.Context, collection *mongo.Collection, selector interface{}, update interface{},
	opts ...*options.UpdateOptions) error {

	_, err := collection.UpdateMany(ctx, selector, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Upsert(ctx context.Context, collection *mongo.Collection, selector interface{}, update interface{},
	opts ...*options.UpdateOptions) error {
	return client.Update(ctx, collection, selector, update, append([]*options.UpdateOptions{options.Update().SetUpsert(true)}, opts...)...)
}

type BulkUpdateItem struct {
	Selector    bson.M
	Update      bson.M
	Replacement interface{}
	IsMulti     bool
}

func (client *Client) BulkUpdateItems(ctx context.Context, collection *mongo.Collection, bulkUpdateItems []*BulkUpdateItem,
	opts ...*options.BulkWriteOptions) error {
	bulkModels := make([]mongo.WriteModel, 0)
	for _, bulkUpdateItem := range bulkUpdateItems {
		if bulkUpdateItem.IsMulti {
			updateModel := mongo.NewUpdateManyModel()
			updateModel.Filter = bulkUpdateItem.Selector
			updateModel.Update = bulkUpdateItem.Update
			bulkModels = append(bulkModels, updateModel)
		} else {
			updateModel := mongo.NewUpdateOneModel()
			updateModel.Filter = bulkUpdateItem.Selector
			updateModel.Update = bulkUpdateItem.Update
			bulkModels = append(bulkModels, updateModel)
		}
	}

	return client.BulkUpdate(ctx, collection, bulkModels, opts...)
}

//The replacement parameter must be a document that will be used to replace the selected document.
//It cannot be nil and cannot contain any update operators
func (client *Client) BulkUpsertItems(ctx context.Context, collection *mongo.Collection, bulkUpdateItems []*BulkUpdateItem,
	opts ...*options.BulkWriteOptions) error {
	bulkModels := make([]mongo.WriteModel, 0)
	for _, bulkUpdateItem := range bulkUpdateItems {
		upsertModel := mongo.NewReplaceOneModel()
		upsertModel.Filter = bulkUpdateItem.Selector
		upsertModel.Replacement = bulkUpdateItem.Replacement
		bulkModels = append(bulkModels, upsertModel)
	}
	return client.BulkUpdate(ctx, collection, bulkModels, opts...)
}

func (client *Client) BulkUpdate(ctx context.Context, collection *mongo.Collection, bulkModels []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) error {
	_, err := collection.BulkWrite(ctx, bulkModels, opts...)
	return err
}

func (client *Client) FindOneAndUpdate(ctx context.Context, collection *mongo.Collection, result, selector, updateMap interface{},
	returnNew bool, upsert bool) error {
	returnDocument := options.Before
	if returnNew {
		returnDocument = options.After
	}

	options := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDocument,
	}
	mongoResult := collection.FindOneAndUpdate(ctx, selector, updateMap, &options)
	if mongoResult.Err() != nil {
		return mongoResult.Err()
	}
	return mongoResult.Decode(result)
}

//read
func (client *Client) FindById(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, result interface{},
	opts ...*options.FindOneOptions) error {
	err := collection.FindOne(ctx, bson.M{"_id": id}, opts...).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) FindByIds(ctx context.Context, collection *mongo.Collection, ids []primitive.ObjectID, result interface{},
	opts ...*options.FindOptions) error {
	query := bson.M{"_id": bson.M{"$in": ids}}

	cur, err := collection.Find(ctx, query, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, result)
}

func (client *Client) Count(ctx context.Context, collection *mongo.Collection, query interface{},
	opts ...*options.CountOptions) (int64, error) {
	number, err := collection.CountDocuments(ctx, query, opts...)
	if err != nil {
		return -1, err
	}

	return number, nil
}

type FindModel struct {
	Query   interface{}
	Fields  interface{}
	Sort    []string
	Cursor  int
	Size    int
	Results interface{}
}

func (client *Client) FindByModel(ctx context.Context, collection *mongo.Collection, model FindModel,
	opts ...*options.FindOptions) error {
	option := options.Find()
	option.SetSkip(int64(model.Cursor))
	option.SetLimit(int64(model.Size))
	option.SetProjection(model.Fields)

	if len(model.Sort) > 0 {
		option.SetSort(populateSliceToBsonD(model.Sort))
	}

	if model.Query == nil {
		model.Query = bson.M{}
	}

	options := []*options.FindOptions{option}
	if len(opts) > 0 {
		options = append(options, opts...)
	}

	cur, err := collection.Find(ctx, model.Query, options...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, model.Results)
}

func populateSliceToBsonD(sortFields []string) bson.D {
	var order bson.D
	for _, field := range sortFields {
		n := 1
		var kind string
		if field != "" {
			if field[0] == '$' {
				if c := strings.Index(field, ":"); c > 1 && c < len(field)-1 {
					kind = field[1:c]
					field = field[c+1:]
				}
			}
			switch field[0] {
			case '+':
				field = field[1:]
			case '-':
				n = -1
				field = field[1:]
			}
		}
		if field == "" {
			continue
		}
		if kind == "textScore" {
			order = append(order, bson.E{Key: field, Value: bson.M{"$meta": kind}})
		} else {
			order = append(order, bson.E{Key: field, Value: n})
		}
	}
	return order
}

func (client *Client) FindOneByModel(ctx context.Context, collection *mongo.Collection, model FindModel,
	opts ...*options.FindOneOptions) error {
	option := options.FindOne()
	option.SetSkip(int64(model.Cursor))
	option.SetProjection(model.Fields)

	if len(model.Sort) > 0 {
		option.SetSort(populateSliceToBsonD(model.Sort))
	}

	if model.Query == nil {
		model.Query = bson.M{}
	}

	options := []*options.FindOneOptions{option}
	if len(opts) > 0 {
		options = append(options, opts...)
	}

	err := collection.FindOne(ctx, model.Query, options...).Decode(model.Results)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Aggregate(ctx context.Context, collection *mongo.Collection, pipeline interface{}, result interface{},
	opts ...*options.AggregateOptions) error {
	cur, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, result)
}
