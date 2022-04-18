package es7

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/liumingmin/goutils/db/elasticsearch"
	"github.com/olivere/elastic/v7"
)

const testUserIndexKey = "testUser"
const testUserIndexName = "test_user2"

type testUser struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Status   string `json:"status"`
	Type     string `json:"pType"`
}

func TestCreateIndexByModel(t *testing.T) {
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
		Settings: Settings{
			IndexMappingIgnoreMalformed: true,
			NumberOfReplicas:            1,
			NumberOfShards:              3,
		},
	})

	t.Log(err)
}

func TestEsInsert(t *testing.T) {
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
		t.Log(err)
	}
}

func TestEsBatchInsert(t *testing.T) {
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
	t.Log(err)
}

func TestEsUpdateById(t *testing.T) {
	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(30)

	err := client.UpdateById(context.Background(), testUserIndexName,
		id, map[string]interface{}{
			"status": "invalid",
		})
	t.Log(err)

	err = client.UpdateById(context.Background(), testUserIndexName,
		id, map[string]interface{}{
			"extField": "ext1234",
		})
	t.Log(err)
}

func TestDeleteById(t *testing.T) {
	InitClients()
	client := GetEsClient(testUserIndexKey)

	id := "000000000" + fmt.Sprint(9)

	err := client.DeleteById(context.Background(), testUserIndexName, id)
	t.Log(err)
}

func TestQueryEs(t *testing.T) {
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
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
}

func TestQueryEsQuerySource(t *testing.T) {
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
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
}

func TestAggregateBySource(t *testing.T) {
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
	client.AggregateBySource(context.Background(), elasticsearch.AggregateModel{
		IndexName: testUserIndexName,
		Source:    source,
		AggKeys:   []string{"status"},
	}, &test)
	t.Log(test)
}

//  github.com/bashtian/jsonutils   json to struct
//  m, err := jsonutils.FromBytes([]byte(jsonStr), "GoStruct")
//  sb := &strings.Builder{}
//	m.Format = true
//	m.Writer = sb
//	m.WriteGo()
type AggregationTest struct {
	Buckets []struct {
		DocCount int64  `json:"doc_count"` // 94
		Key      string `json:"key"`       // valid
		PType    struct {
			Buckets []struct {
				DocCount  int64  `json:"doc_count"` // 84
				Key       string `json:"key"`       // normal
				StatusCnt struct {
					Value int64 `json:"value"` // 84
				} `json:"statusCnt"`
			} `json:"buckets"`
			DocCountErrorUpperBound int64 `json:"doc_count_error_upper_bound"` // 0
			SumOtherDocCount        int64 `json:"sum_other_doc_count"`         // 0
		} `json:"pType"`
	} `json:"buckets"`
	DocCountErrorUpperBound int64 `json:"doc_count_error_upper_bound"` // 0
	SumOtherDocCount        int64 `json:"sum_other_doc_count"`         // 0
}
