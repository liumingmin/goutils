package utils

import (
	"reflect"
	"strconv"
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
	Id          string `form:"id" json:"id"`
	ServiceName string `form:"serviceName" json:"serviceName"`
	Body        string `form:"body" json:"body"`
	Version     int    `form:"version" json:"version"`
	UpdateTime  string `form:"updateTime" json:"updateTime"`
	Test        string

	MoreProp  string
	MoreProp2 int
	MoreProp3 string
	MoreProp4 float64
}

type ConfItemExtend struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	StringId  string
	MoreProp  string
	MoreProp2 int
	MoreProp3 time.Time
	MoreProp4 float64
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

func TestMergeStructs(t *testing.T) {
	var vos []*ConfItemVo
	var dos = []*ConfItem{
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
		{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111},
	}
	CopyStructs(dos, &vos, BaseConvert)

	for _, vo := range vos {
		t.Log(vo)
	}

	var extends []ConfItemExtend
	for i, do := range dos {
		extends = append(extends, ConfItemExtend{
			Id:        do.Id,
			StringId:  do.Id.Hex(),
			MoreProp:  strconv.Itoa(i),
			MoreProp2: i,
			MoreProp3: time.Now(),
			MoreProp4: float64(i) * 3.3,
		})
	}

	MergeStructs(extends, &vos, BaseConvert, "StringId:Id",
		"MoreProp:MoreProp", "MoreProp3:MoreProp3", "MoreProp4:MoreProp4",
	)

	for _, vo := range vos {
		t.Log(vo)
	}
}

func BenchmarkMergeStructs(b *testing.B) {
	b.StopTimer()
	var vos []*ConfItemVo
	var dos []*ConfItem
	for i := 0; i < 500; i++ {
		dos = append(dos, &ConfItem{Id: bson.NewObjectId(), ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now(), Test: 111})
	}

	CopyStructs(dos, &vos, BaseConvert)

	var extends []ConfItemExtend
	for i, do := range dos {
		extends = append(extends, ConfItemExtend{
			Id:        do.Id,
			StringId:  do.Id.Hex(),
			MoreProp:  strconv.Itoa(i),
			MoreProp2: i,
			MoreProp3: time.Now(),
			MoreProp4: float64(i) * 3.3,
		})
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		MergeStructs(extends, &vos, BaseConvert, "Id:Id",
			"MoreProp:MoreProp", "MoreProp1:MoreProp1", "MoreProp4:MoreProp4",
		)
	}
}
