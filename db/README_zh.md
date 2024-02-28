**其他语言版本: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [数据库](#%E6%95%B0%E6%8D%AE%E5%BA%93)
  * [ES搜索引擎](#es%E6%90%9C%E7%B4%A2%E5%BC%95%E6%93%8E)
  * [kafka消息队列](#kafka%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  * [mongo数据库](#mongo%E6%95%B0%E6%8D%AE%E5%BA%93)
  * [go-redis](#go-redis)

<!-- tocstop -->

# 数据库
## ES搜索引擎
### ES6版本API
#### es_test.go
##### TestCreateIndexByModel
```go

InitClients()

client := GetEsClient(testUserIndexKey)

err := client.CreateIndexByModel(context.Background(), testUserIndexName, &MappingModel{
	Mappings: map[string]Mapping{
		testUserTypeName: {
			Dynamic: false,
			Properties: map[string]*elasticsearch.MappingProperty{
				"userId": {
					Type:  "text",
					Index: false,
				},
				"nickname": {
					Type:     "text",
					Analyzer: "standard",
					Fields: map[string]*elasticsearch.MappingProperty{
						"std": {
							Type:     "text",
							Analyzer: "standard",
							ExtProps: map[string]interface{}{
								"term_vector": "with_offsets",
							},
						},
						"keyword": {
							Type: "keyword",
						},
					},
				},
				"status": {
					Type: "keyword",
				},
				"pType": {
					Type: "keyword",
				},
			},
		},
	},
	Settings: elasticsearch.MappingSettings{
		SettingsIndex: elasticsearch.SettingsIndex{
			IgnoreMalformed:  true,
			NumberOfReplicas: 1,
			NumberOfShards:   3,
		},
	},
})

if err != nil {
	t.Error(err)
}
```
##### TestEsInsert
```go

InitClients()

client := GetEsClient(testUserIndexKey)

for i := 0; i < 100; i++ {
	ptype := "normal"
	if i%10 == 5 {
		ptype = "vip"
	}
	status := "valid"
	if i%30 == 2 {
		status = "invalid"
	}
	id := "000000000" + fmt.Sprint(i)
	err := client.Insert(context.Background(), testUserIndexName, testUserTypeName,
		id, testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})

	if err != nil {
		t.Error(err)
	}
}
```
##### TestEsBatchInsert
```go

InitClients()

client := GetEsClient(testUserIndexKey)

ids := make([]string, 0)
items := make([]interface{}, 0)

for i := 0; i < 100; i++ {
	ptype := "normal"
	if i%10 == 5 {
		ptype = "vip"
	}
	status := "valid"
	if i%30 == 2 {
		status = "invalid"
	}
	id := "x00000000" + fmt.Sprint(i)

	ids = append(ids, id)
	items = append(items, &testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
}

err := client.BatchInsert(context.Background(), testUserIndexName, testUserTypeName, ids, items)
if err != nil {
	t.Error(err)
}
```
##### TestEsUpdateById
```go

InitClients()
client := GetEsClient(testUserIndexKey)

id := "000000000" + fmt.Sprint(30)

err := client.UpdateById(context.Background(), testUserIndexName, testUserTypeName,
	id, map[string]interface{}{
		"status": "invalid",
	})
if err != nil {
	t.Error(err)
}

err = client.UpdateById(context.Background(), testUserIndexName, testUserTypeName,
	id, map[string]interface{}{
		"extField": "ext1234",
	})
if err != nil {
	t.Error(err)
}
```
##### TestDeleteById
```go

InitClients()
client := GetEsClient(testUserIndexKey)

id := "000000000" + fmt.Sprint(9)

err := client.DeleteById(context.Background(), testUserIndexName, testUserTypeName,
	id)
if err != nil {
	t.Error(err)
}
```
##### TestQueryEs
```go

InitClients()
client := GetEsClient(testUserIndexKey)

bq := elastic.NewBoolQuery()
bq.Must(elastic.NewMatchQuery("nickname", "超级棒"))

var users []testUser
total := int64(0)
err := client.FindByModel(context.Background(), elasticsearch.QueryModel{
	IndexName: testUserIndexName,
	TypeName:  testUserTypeName,
	Query:     bq,
	Size:      5,
	Results:   &users,
	Total:     &total,
})
if err != nil {
	t.Error(err)
}
```
##### TestQueryEsQuerySource
```go

InitClients()
client := GetEsClient(testUserIndexKey)

source := `{
	"from":0,
	"size":25,
	"query":{
		"match":{"nickname":"超级"}
	}
}`

var users []testUser
total := int64(0)
err := client.FindBySource(context.Background(), elasticsearch.SourceModel{
	IndexName: testUserIndexName,
	TypeName:  testUserTypeName,
	Source:    source,
	Results:   &users,
	Total:     &total,
})
if err != nil {
	t.Error(err)
}
```
##### TestAggregateBySource
```go

InitClients()
client := GetEsClient(testUserIndexKey)

source := `{
	"from": 0,
	"size": 0,
	"_source": {
		"includes": [
			"status",
			"pType",
			"COUNT"
		],
		"excludes": []
	},
	"stored_fields": [
		"status",
		"pType"
	],
	"aggregations": {
		"status": {
			"terms": {
				"field": "status",
				"size": 200,
				"min_doc_count": 1,
				"shard_min_doc_count": 0,
				"show_term_doc_count_error": false,
				"order": [
					{
						"_count": "desc"
					},
					{
						"_key": "asc"
					}
				]
			},
			"aggregations": {
				"pType": {
					"terms": {
						"field": "pType",
						"size": 10,
						"min_doc_count": 1,
						"shard_min_doc_count": 0,
						"show_term_doc_count_error": false,
						"order": [
							{
								"_count": "desc"
							},
							{
								"_key": "asc"
							}
						]
					},
					"aggregations": {
						"statusCnt": {
							"value_count": {
								"field": "_index"
							}
						}
					}
				}
			}
		}
	}
}`

var test AggregationTest
err := client.AggregateBySource(context.Background(), elasticsearch.AggregateModel{
	IndexName: testUserIndexName,
	TypeName:  testUserTypeName,
	Source:    source,
	AggKeys:   []string{"status"},
}, &test)
if err != nil {
	t.Error(err)
}
```
### ES7版本API
#### es_test.go
##### TestCreateIndexByModel
```go

InitClients()

client := GetEsClient(testUserIndexKey)

err := client.CreateIndexByModel(context.Background(), testUserIndexName, &MappingModel{
	Mapping: Mapping{
		Dynamic: false,
		Properties: map[string]*elasticsearch.MappingProperty{
			"userId": {
				Type:  "text",
				Index: false,
			},
			"nickname": {
				Type:     "text",
				Analyzer: "standard",
				Fields: map[string]*elasticsearch.MappingProperty{
					"std": {
						Type:     "text",
						Analyzer: "standard",
						ExtProps: map[string]interface{}{
							"term_vector": "with_offsets",
						},
					},
					"keyword": {
						Type: "keyword",
					},
				},
			},
			"status": {
				Type: "keyword",
			},
			"pType": {
				Type: "keyword",
			},
		},
	},
	Settings: elasticsearch.MappingSettings{
		SettingsIndex: elasticsearch.SettingsIndex{
			IgnoreMalformed:  true,
			NumberOfReplicas: 2,
			NumberOfShards:   3,
		},
	},
})

if err != nil {
	t.Error(err)
}
```
##### TestEsInsert
```go

InitClients()
client := GetEsClient(testUserIndexKey)

for i := 0; i < 100; i++ {
	ptype := "normal"
	if i%10 == 5 {
		ptype = "vip"
	}
	status := "valid"
	if i%30 == 2 {
		status = "invalid"
	}
	id := "000000000" + fmt.Sprint(i)
	err := client.Insert(context.Background(), testUserIndexName,
		id, testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
	if err != nil {
		t.Error(err)
	}
}
```
##### TestEsBatchInsert
```go

InitClients()
client := GetEsClient(testUserIndexKey)

ids := make([]string, 0)
items := make([]interface{}, 0)

for i := 0; i < 100; i++ {
	ptype := "normal"
	if i%10 == 5 {
		ptype = "vip"
	}
	status := "valid"
	if i%30 == 2 {
		status = "invalid"
	}
	id := "x00000000" + fmt.Sprint(i)

	ids = append(ids, id)
	items = append(items, &testUser{UserId: id, Nickname: "超级棒" + id, Status: status, Type: ptype})
}

err := client.BatchInsert(context.Background(), testUserIndexName, ids, items)
if err != nil {
	t.Error(err)
}
```
##### TestEsUpdateById
```go

InitClients()
client := GetEsClient(testUserIndexKey)

id := "000000000" + fmt.Sprint(30)

err := client.UpdateById(context.Background(), testUserIndexName,
	id, map[string]interface{}{
		"status": "invalid",
	})
if err != nil {
	t.Error(err)
}

err = client.UpdateById(context.Background(), testUserIndexName,
	id, map[string]interface{}{
		"extField": "ext1234",
	})
if err != nil {
	t.Error(err)
}
```
##### TestDeleteById
```go

InitClients()
client := GetEsClient(testUserIndexKey)

id := "000000000" + fmt.Sprint(9)

err := client.DeleteById(context.Background(), testUserIndexName, id)
if err != nil {
	t.Error(err)
}
```
##### TestQueryEs
```go

InitClients()
client := GetEsClient(testUserIndexKey)

bq := elastic.NewBoolQuery()
bq.Must(elastic.NewMatchQuery("nickname", "超级棒"))

var users []testUser
total := int64(0)
err := client.FindByModel(context.Background(), elasticsearch.QueryModel{
	IndexName: testUserIndexName,
	Query:     bq,
	Size:      5,
	Results:   &users,
	Total:     &total,
})
if err != nil {
	t.Error(err)
}
//bs, _ := json.Marshal(users)
//t.Log(len(users), total, string(bs), err)
```
##### TestQueryEsQuerySource
```go

InitClients()
client := GetEsClient(testUserIndexKey)
source := `{
	"from":0,
	"size":25,
	"query":{
		"match":{"nickname":"超级"}
	}
}`

var users []testUser
total := int64(0)
err := client.FindBySource(context.Background(), elasticsearch.SourceModel{
	IndexName: testUserIndexName,
	Source:    source,
	Results:   &users,
	Total:     &total,
})
if err != nil {
	t.Error(err)
}
//bs, _ := json.Marshal(users)
//t.Log(len(users), total, string(bs), err)
```
##### TestAggregateBySource
```go

InitClients()
client := GetEsClient(testUserIndexKey)
source := `{
	"from": 0,
	"size": 0,
	"_source": {
		"includes": [
			"status",
			"pType",
			"COUNT"
		],
		"excludes": []
	},
	"stored_fields": [
		"status",
		"pType"
	],
	"aggregations": {
		"status": {
			"terms": {
				"field": "status",
				"size": 200,
				"min_doc_count": 1,
				"shard_min_doc_count": 0,
				"show_term_doc_count_error": false,
				"order": [
					{
						"_count": "desc"
					},
					{
						"_key": "asc"
					}
				]
			},
			"aggregations": {
				"pType": {
					"terms": {
						"field": "pType",
						"size": 10,
						"min_doc_count": 1,
						"shard_min_doc_count": 0,
						"show_term_doc_count_error": false,
						"order": [
							{
								"_count": "desc"
							},
							{
								"_key": "asc"
							}
						]
					},
					"aggregations": {
						"statusCnt": {
							"value_count": {
								"field": "_index"
							}
						}
					}
				}
			}
		}
	}
}`

var test AggregationTest
err := client.AggregateBySource(context.Background(), elasticsearch.AggregateModel{
	IndexName: testUserIndexName,
	Source:    source,
	AggKeys:   []string{"status"},
}, &test)
if err != nil {
	t.Error(err)
}
```
## kafka消息队列
### kafka_test.go
#### TestKafkaProducer
```go

InitKafka()
producer := GetProducer("user_producer")
err := producer.Produce(&sarama.ProducerMessage{
	Topic: userTopic,
	Key:   sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
	Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
})

if err != nil {
	t.Error(err)
}

time.Sleep(time.Millisecond * 100)
```
#### TestKafkaConsumer
```go

InitKafka()

consumer := GetConsumer("user_consumer")
go func() {
	consumer.Consume(userTopic, func(msg *sarama.ConsumerMessage) error {
		fmt.Println(string(msg.Key), "=", string(msg.Value))
		return nil
	}, func(err error) {

	})
}()

producer := GetProducer("user_producer")
for i := 0; i < 10; i++ {
	err := producer.Produce(&sarama.ProducerMessage{
		Topic: userTopic,
		Key:   sarama.ByteEncoder(fmt.Sprint(i)),
		Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
	})
	if err != nil {
		t.Error(err)
	}
}

time.Sleep(time.Millisecond * 100)
```
## mongo数据库
### collection_test.go
#### TestInsert
```go

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
if err != nil {
	t.Error(err)
}

var result []bson.M
err = op.Find(ctx, FindModel{
	Query:   bson.M{"user_id": "1"},
	Results: &result,
})
if err != nil {
	t.Error(err)
}
for _, item := range result {
	t.Log(item)
}
```
#### TestUpdate
```go

ctx := context.Background()
InitClients()
c, _ := MgoClient(dbKey)

op := NewCompCollectionOp(c, dbName, collectionName)
err := op.Update(ctx, bson.M{"user_id": "1"}, bson.M{"$set": bson.M{"nick_name": "超级棒++"}})
if err != nil {
	t.Error(err)
}

var result interface{}
err = op.FindOne(ctx, FindModel{
	Query:   bson.M{"user_id": "1"},
	Results: &result,
})
if err != nil {
	t.Error(err)
}

log.Info(ctx, "result: %v", result)
```
#### TestFind
```go

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
```
#### TestDelete
```go

ctx := context.Background()
InitClients()
c, _ := MgoClient(dbKey)

op := NewCompCollectionOp(c, dbName, collectionName)
err := op.Delete(ctx, bson.M{"user_id": "1"})
if err != nil {
	t.Error(err)
}

var result interface{}
err = op.FindOne(ctx, FindModel{
	Query:   bson.M{"user_id": "1"},
	Results: &result,
})
if err != mongo.ErrNoDocuments {
	t.Error(err)
}
log.Info(ctx, "result: %v", result)
```
#### TestUpert
```go

ctx := context.Background()
InitClients()
c, _ := MgoClient(dbKey)

op := NewCompCollectionOp(c, dbName, collectionName)
err := op.Upsert(ctx, bson.M{"name": "tom2"}, bson.M{"$set": bson.M{"birth": "2020"}}, bson.M{"birth2": "2024"})
if err != nil {
	t.Error(err)
}
```
#### TestBulkUpdateItems
```go

ctx := context.Background()
InitClients()
c, _ := MgoClient(dbKey)

op := NewCompCollectionOp(c, dbName, collectionName)

err := op.BulkUpdateItems(ctx, []*BulkUpdateItem{
	{Selector: bson.M{"name": "tom"}, Update: bson.M{"$set": bson.M{"birth": "1"}}},
	{Selector: bson.M{"name": "tom1"}, Update: bson.M{"$set": bson.M{"birth2": "2"}}},
})
if err != nil {
	t.Error(err)
}
```
#### TestBulkUpsertItems
```go

ctx := context.Background()
InitClients()
c, _ := MgoClient(dbKey)

op := NewCompCollectionOp(c, dbName, collectionName)

err := op.BulkUpsertItem(ctx, []*BulkUpsertItem{
	{Selector: bson.M{"name": "tim"}, Replacement: bson.M{"name": "tim", "birth": "3"}},
	{Selector: bson.M{"name": "tim1"}, Replacement: bson.M{"name": "tim1", "birth2": "4"}},
})
if err != nil {
	t.Error(err)
}
```
## go-redis
### Redis List工具库
#### TestList
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()

err := ListPush(ctx, rds, "test_list", "stringvalue")
if err != nil {
	t.Error(err)
}
ListPop(rds, []string{"test_list"}, 3600, 1000, func(key, data string) {
	fmt.Println(key, data)
})

err = ListPush(ctx, rds, "test_list", "stringvalue")
if err != nil {
	t.Error(err)
}
```
### Redis 锁工具库
#### TestRdsAllowActionWithCD
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()

cd, ok := RdsAllowActionWithCD(ctx, rds, "test:action", 2)
t.Log(cd, ok)
cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
t.Log(cd, ok)
time.Sleep(time.Second * 3)

cd, ok = RdsAllowActionWithCD(ctx, rds, "test:action", 2)
t.Log(cd, ok)
```
#### TestRdsAllowActionByMTs
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()

cd, ok := RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
t.Log(cd, ok)
cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
t.Log(cd, ok)
time.Sleep(time.Second)

cd, ok = RdsAllowActionByMTs(ctx, rds, "test:action", 500, 10)
t.Log(cd, ok)
```
#### TestRdsLockResWithCD
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()

ok := RdsLockResWithCD(ctx, rds, "test:res", "res-1", 3)
t.Log(ok)
ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
t.Log(ok)
time.Sleep(time.Second * 4)

ok = RdsLockResWithCD(ctx, rds, "test:res", "res-2", 3)
t.Log(ok)
```
### Redis PubSub工具库
#### TestMqPSubscribe
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()

MqPSubscribe(ctx, rds, "testkey:*", func(channel string, data string) {
	fmt.Println(channel, data)
}, 10)

err := MqPublish(ctx, rds, "testkey:1", "id:1")
if err != nil {
	t.Error(err)
}
err = MqPublish(ctx, rds, "testkey:2", "id:2")
if err != nil {
	t.Error(err)
}
err = MqPublish(ctx, rds, "testkey:3", "id:3")
if err != nil {
	t.Error(err)
}
```
### redis_test.go
#### TestSentinel
```go

InitRedises()
rds := Get("rds-sentinel")
ctx := context.Background()

rds.Set(ctx, "test_senti", "test_value", time.Minute)

_, err := rds.Get(ctx, "test_senti").Result()
if err != nil {
	t.Error(err)
}
```
#### TestCluster
```go

InitRedises()
rds := Get("rds-cluster")
ctx := context.Background()

rds.Set(ctx, "test_cluster", "test_value", time.Minute)

_, err := rds.Get(ctx, "test_cluster").Result()
if err != nil {
	t.Error(err)
}
```
### Redis ZSet工具库
#### TestZDescartes
```go

InitRedises()
rds := Get("rdscdb")
ctx := context.Background()
dimValues := [][]string{{"dim1a", "dim1b"}, {"dim2a", "dim2b", "dim2c", "dim2d"}, {"dim3a", "dim3b", "dim3c"}}

dt, err := csv.ReadCsvToDataTable(ctx, "data.csv", ',',
	[]string{"id", "name", "createtime", "dim1", "dim2", "dim3", "member"}, "id", []string{})
if err != nil {
	t.Error(err)
}

err = ZDescartes(ctx, rds, dimValues, func(strs []string) (string, map[string]int64) {
	dimData := make(map[string]int64)
	for _, row := range dt.Rows() {
		if row.String("dim1") == strs[0] &&
			row.String("dim2") == strs[1] &&
			row.String("dim3") == strs[2] {
			dimData[row.String("member")] = row.Int64("createtime")
		}
	}
	return "rds" + strings.Join(strs, "-"), dimData
}, 1000, 30)

if err != nil {
	t.Error(err)
}
```
