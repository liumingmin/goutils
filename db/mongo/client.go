package mongo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/net/proxy"

	"go.mongodb.org/mongo-driver/bson"
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
	dbs := conf.Conf.Mongos
	if len(dbs) == 0 {
		return
	}

	for _, database := range dbs {
		client, err := initClient(database)
		if err != nil {
			continue
		}

		mgoClients[database.Key] = client
	}
}

func initClient(dbconf *conf.Mongo) (ret *Client, err error) {
	opts := options.Client()
	opts.SetHosts(dbconf.Addrs)
	if dbconf.User != "" {
		opts.SetAuth(genAuthFromConf(dbconf))
	}

	if len(dbconf.Compressors) > 0 {
		opts.SetCompressors(dbconf.Compressors)
	}

	dialer, err := genDialerFromConf(dbconf)
	if err == nil {
		opts.SetDialer(dialer)
	}

	if dbconf.ConnectTimeout > 0 {
		opts.SetConnectTimeout(dbconf.ConnectTimeout)
	}

	opts.SetDirect(dbconf.Direct)

	socketTimeout := dbconf.ReadTimeout
	if socketTimeout < dbconf.WriteTimeout {
		socketTimeout = dbconf.WriteTimeout
	}
	if socketTimeout > 0 {
		opts.SetSocketTimeout(socketTimeout)
	}

	if dbconf.MaxPoolIdleTimeMS > 0 {
		opts.SetMaxConnIdleTime(time.Duration(dbconf.MaxPoolIdleTimeMS) * time.Millisecond)
	}
	if dbconf.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(uint64(dbconf.MaxPoolSize))
	}
	if dbconf.MinPoolSize > 0 {
		opts.SetMinPoolSize(uint64(dbconf.MinPoolSize))
	}
	if dbconf.MaxPoolWaitTimeMS > 0 {
		opts.SetServerSelectionTimeout(time.Duration(dbconf.MaxPoolWaitTimeMS) * time.Millisecond)
	}

	//mode
	mode := getMode(dbconf.Mode)
	if mode > 0 {
		var pref *readpref.ReadPref
		if pref, err = readpref.New(mode); err != nil {
			return
		}
		opts.SetReadPreference(pref)
	}

	//read write safe
	if dbconf.Safe == nil {
		dbconf.Safe = &conf.MongoSafe{W: 1}
	}

	if dbconf.Safe.RMode != "" {
		opts.SetReadConcern(&readconcern.ReadConcern{Level: strings.ToLower(dbconf.Safe.RMode)})
	}
	opts.SetWriteConcern(genWriteConcernFromConf(dbconf))

	opts.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: true,
	})

	//finished
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

func genAuthFromConf(dbconf *conf.Mongo) options.Credential {
	var auth options.Credential
	auth.Username = dbconf.User
	auth.Password = dbconf.Password
	auth.PasswordSet = true
	if dbconf.AuthSource != "" {
		auth.AuthSource = dbconf.AuthSource
	} else {
		auth.AuthSource = dbconf.DBName
	}
	return auth
}

func genDialerFromConf(dbconf *conf.Mongo) (options.ContextDialer, error) {
	if dbconf.Ssh != nil && dbconf.Ssh.On {
		dialer, err := genSshFromConf(dbconf)
		if err != nil {
			return nil, err
		}
		return dialer, nil
	} else {
		var dialer net.Dialer
		if dbconf.Keepalive > 0 {
			dialer.KeepAlive = dbconf.Keepalive
		}

		if dbconf.ConnectTimeout > 0 {
			dialer.Timeout = dbconf.ConnectTimeout
		}
		return &dialer, nil
	}
}

func genSshFromConf(dbconf *conf.Mongo) (options.ContextDialer, error) {
	var sshKey []byte
	var err error
	if dbconf.Ssh.PriKey != "" {
		sshKey, err = base64.StdEncoding.DecodeString(dbconf.Ssh.PriKey)
		if err != nil {
			log.Error(context.Background(), "Decode sshKeyB64 failed: %v", err)
			return nil, err
		}
	}

	sshConfig, err := proxy.NewSshClient(dbconf.Ssh.Address, dbconf.Ssh.User, sshKey, dbconf.Ssh.KeyPass)
	if err != nil {
		log.Error(context.Background(), "NewSshClient err: %v", err)
		return nil, err
	}

	return sshConfig, nil
}

func genWriteConcernFromConf(dbconf *conf.Mongo) *writeconcern.WriteConcern {
	var writeConcern *writeconcern.WriteConcern
	if strings.ToLower(dbconf.Safe.WMode) == "majority" {
		writeConcern = writeconcern.Majority()
	} else {
		writeConcern = &writeconcern.WriteConcern{}
	}

	if dbconf.Safe.J {
		writeConcern.Journal = &dbconf.Safe.J
	}
	if dbconf.Safe.W > 0 {
		writeConcern.W = dbconf.Safe.W
	}
	if dbconf.Safe.WTimeout > 0 {
		writeConcern.WTimeout = time.Duration(dbconf.Safe.WTimeout) * time.Millisecond
	}
	return writeConcern
}

func getMode(val interface{}) readpref.Mode {
	switch val := val.(type) {
	case string:
		if val == "" {
			return readpref.PrimaryMode
		}

		ret, err := readpref.ModeFromString(val)
		if err != nil {
			return readpref.PrimaryMode
		}
		return ret
	case int:
		return readpref.Mode(val)
	case int64:
		return readpref.Mode(val)
	}
	log.Error(context.Background(), "unsupport mode type: "+fmt.Sprint(val))
	return readpref.PrimaryMode
}

// 需要mongodb4.2以上才能支持事务
func (client *Client) ExecTx(ctx context.Context, f func(context.Context)) error {
	return client.UseSession(ctx, func(sctx mongo.SessionContext) (err error) {
		err = sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.Majority()),
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

// write
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
	Selector bson.M
	Update   bson.M
	IsMulti  bool
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

type BulkUpsertItem struct {
	Selector    bson.M
	Replacement interface{}
}

// The replacement parameter must be a document that will be used to replace the selected document.
// It cannot be nil and cannot contain any update operators
func (client *Client) BulkUpsertItems(ctx context.Context, collection *mongo.Collection, bulkUpsertItems []*BulkUpsertItem,
	opts ...*options.BulkWriteOptions) error {
	bulkModels := make([]mongo.WriteModel, 0)
	upsert := true
	for _, bulkUpsertItem := range bulkUpsertItems {
		upsertModel := mongo.NewReplaceOneModel()
		upsertModel.Filter = bulkUpsertItem.Selector
		upsertModel.Replacement = bulkUpsertItem.Replacement
		upsertModel.Upsert = &upsert
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

// read
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
