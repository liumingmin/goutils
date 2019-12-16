package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

type ConfItem struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ServiceName string        `bson:"service_name" json:"service_name"`
	Body        string        `bson:"body" json:"body"`
	Version     int           `bson:"version" json:"version"`
	UpdateTime  time.Time     `bson:"update_time" json:"update_time"`
	Test        int
}

type ConfItemVo struct {
	Id          string    `form:"id" json:"id"`
	ServiceName string    `form:"serviceName" json:"serviceName"`
	Body        string    `form:"body" json:"body"`
	Version     int       `form:"version" json:"version"`
	UpdateTime  time.Time `form:"updateTime" json:"updateTime"`
	Test        string
}

func TestCopyStruct(t *testing.T) {
	vo := &ConfItemVo{}
	do := ConfItem{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()}
	CopyStructDefault(do, vo)
	t.Log(vo)
}

func TestCopyStructs(t *testing.T) {
	var vos []ConfItemVo
	var dos = []*ConfItem{
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
	}

	err := CopyStructs(&dos, &vos, func(src interface{}, dstType reflect.Type) interface{} {
		if bid, ok := src.(bson.ObjectId); ok && dstType.Kind() == reflect.String {
			return bid.Hex()
		}
		return nil
	})
	t.Log(err)
	t.Log(vos)
}

func TestConvertFieldStyle(t *testing.T) {
	t.Log(ConvertFieldStyle("", TAG_STYLE_UNDERLINE))
	t.Log(ConvertFieldStyle("a", TAG_STYLE_UNDERLINE))
	t.Log(ConvertFieldStyle("aB", TAG_STYLE_UNDERLINE))
	t.Log(ConvertFieldStyle("AB", TAG_STYLE_UNDERLINE))

	t.Log("------------------------------------------")
	t.Log(ConvertFieldStyle("", TAG_STYLE_SNAKE))
	t.Log(ConvertFieldStyle("a", TAG_STYLE_SNAKE))
	t.Log(ConvertFieldStyle("aB", TAG_STYLE_SNAKE))
	t.Log(ConvertFieldStyle("AB", TAG_STYLE_SNAKE))

	t.Log("------------------------------------------")
	t.Log(ConvertFieldStyle("TAestConvertFieldStyleAddZd$sT好z", TAG_STYLE_UNDERLINE))
	t.Log(ConvertFieldStyle("TAestConvertFieldStyleAddZd$sT好z", TAG_STYLE_SNAKE))
}

func TestDoToVo(t *testing.T) {
	t.Log(AutoGenTags(ConfItemVo{}, map[string]TAG_STYLE{"json": TAG_STYLE_SNAKE, "form": TAG_STYLE_SNAKE,
		"bson": TAG_STYLE_UNDERLINE}))
}
