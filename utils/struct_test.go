package utils

import (
	"testing"
	"time"
)

type ConfItem struct {
	Id          int       `bson:"_id,omitempty" json:"id"`
	ServiceName string    `bson:"service_name" json:"service_name"`
	Body        string    `bson:"body" json:"body"`
	Version     int       `bson:"version" json:"version"`
	UpdateTime  time.Time `bson:"update_time" json:"update_time"`
}

type ConfItemVo struct {
	Id          string    `form:"id" json:"id"`
	ServiceName string    `form:"serviceName" json:"serviceName"`
	Body        string    `form:"body" json:"body"`
	Version     int       `form:"version" json:"version"`
	UpdateTime  time.Time `form:"updateTime" json:"updateTime"`
}

func TestCopyStruct(t *testing.T) {
	vo := &ConfItemVo{}
	do := ConfItem{Id: 12345, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()}
	CopyStruct(do, vo)
	t.Log(vo)
}

func TestCopyStructs(t *testing.T) {
	var vos []ConfItemVo
	var dos = []*ConfItem{
		{Id: 1234, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()},
		{Id: 1234, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()},
		{Id: 1234, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()},
		{Id: 1234, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()},
		{Id: 1234, ServiceName: "test", Body: "testBody", Version: 2, UpdateTime: time.Now()},
	}

	err := CopyStructs(&dos, &vos)
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
