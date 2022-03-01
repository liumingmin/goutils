package es

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/liumingmin/goutils/db/elasticsearch"

	"github.com/olivere/elastic"
)

func TestQueryEs(t *testing.T) {
	InitClients()

	bq := elastic.NewBoolQuery()
	bq.Must(elastic.NewMatchQuery("nickname", "燕双鹰"))

	var users []EsSocialUser
	total := int64(0)
	err := FindByModel(context.Background(), elasticsearch.QueryModel{
		BaseModel: elasticsearch.BaseModel{
			KeyName:   "person_cluster",
			IndexName: "socialuser00",
		},
		Query:   bq,
		Size:    5,
		Results: &users,
		Total:   &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(total, string(bs), err)
}

func TestQueryEsQuerySource(t *testing.T) {
	InitClients()

	source := `{
	"from":0,
	"size":25,
	"query":{
		"match":{"nickname":"燕双鹰"}
	}
}`

	var users []EsSocialUser
	total := int64(0)
	err := FindBySource(context.Background(), elasticsearch.SourceModel{
		BaseModel: elasticsearch.BaseModel{
			KeyName:   "person_cluster",
			IndexName: "socialuser00",
		},
		Source:  source,
		Results: &users,
		Total:   &total,
	})
	bs, _ := json.Marshal(users)
	t.Log(len(users), total, string(bs), err)
}

func TestAggregateBySource(t *testing.T) {
	InitClients()

	source := `{
    "from": 0,
    "size": 0,
    "_source": {
        "includes": [
            "status",
            "chessLadderLevel",
            "COUNT"
        ],
        "excludes": []
    },
    "stored_fields": [
        "status",
        "chessLadderLevel"
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
                "chessLadderLevel": {
                    "terms": {
                        "field": "chessLadderLevel",
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
	AggregateBySource(context.Background(), elasticsearch.AggregateModel{
		BaseModel: elasticsearch.BaseModel{
			KeyName:   "person_cluster",
			IndexName: "socialuser00",
		},
		Source:  source,
		AggKeys: []string{"status"},
	}, &test)
	t.Log(test)
}

type AggregationTest struct {
	Buckets []struct {
		ChessLadderLevel struct {
			Buckets []struct {
				DocCount  int64 `json:"doc_count"` // 21
				Key       int64 `json:"key"`       // 0
				StatusCnt struct {
					Value int64 `json:"value"` // 21
				} `json:"statusCnt"`
			} `json:"buckets"`
			DocCountErrorUpperBound int64 `json:"doc_count_error_upper_bound"` // 0
			SumOtherDocCount        int64 `json:"sum_other_doc_count"`         // 0
		} `json:"chessLadderLevel"`
		DocCount int64  `json:"doc_count"` // 2.265254e+06
		Key      string `json:"key"`       // valid
	} `json:"buckets"`
	DocCountErrorUpperBound int64 `json:"doc_count_error_upper_bound"` // 0
	SumOtherDocCount        int64 `json:"sum_other_doc_count"`         // 0
}

func TestEsUpdateById(t *testing.T) {
	InitClients()

	err := UpdateById(context.Background(), "person_cluster", "socialuser00", "_doc",
		"b28f2ab1448345f2837d0a11d88f80e2", map[string]interface{}{
			"testUpdate": "test",
		})
	t.Log(err)
}

func TestEsInsert(t *testing.T) {
	InitClients()

	err := Insert(context.Background(), "person_cluster", "socialuser00", "_doc",
		"0000000000001", EsSocialUser{PersonID: "0000000000001"})
	t.Log(err)
}

type EsSocialUser struct {
	PersonID        string    `json:"personID" binding:"required"` //自然人Id
	Nickname        string    `json:"nickname,omitempty"`          //自然人昵称
	PersonNum       int64     `json:"personNum"`                   //助手Id
	AvatarURL       string    `json:"avatarURL"`                   //自然人头像地址
	Status          string    `json:"status,omitempty"`            //自然人状态, valid/invalid
	Type            string    `json:"pType,omitempty"`             //自然人类型, OFFICIAL/GENERAL/
	ExternalId      string    `json:"externalId,omitempty"`        //注册来源ID
	ExternalType    string    `json:"externalType,omitempty"`      //注册来源类型 10: wegame, 100:指尖
	UpdateDate      time.Time `json:"updateDate,omitempty"`        //数据库更新时间
	RegisterDate    time.Time `json:"registerDate,omitempty"`      //自然人注册时间
	PassportList    []string  `json:"passports,omitempty"`         //关联通行证
	Zone            []string  `json:"zone,omitempty"`              //关联圈子
	Fans            int       `json:"fans,omitempty"`              //粉丝数量
	Articles        int       `json:"articles,omitempty"`          //作品数量
	ActivityFeeds   int       `json:"activityFeeds,omitempty"`     //动态数量
	TotalRead       int       `json:"totalRead,omitempty"`         //总阅读量
	LastSignInTs    int64     `json:"lastSignInTs,omitempty"`      //最后一次签到时间戳
	IsVIP           bool      `json:"isVIP,omitempty"`             //是否是达人
	LogonRecommend  bool      `json:"logonRecommend,omitempty"`    // 首登推荐
	HighlyRecommend bool      `json:"highlyRecommend,omitempty"`   // 重点推荐
	Introduce       string    `json:"introduce,omitempty"`         // 介绍
	Level           int       `json:"level,omitempty"`             // 等级
	TalentInfo      string    `json:"talentInfo,omitempty"`        // 达人信息[["标记","等级"],...]
	DefaultTalent   string    `json:"defaultTalent"`
}
