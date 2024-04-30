package conv

import (
	"reflect"
	"testing"
)

func TestValueToString(t *testing.T) {
	var str string
	var err error

	str, err = ValueToString(100.2)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[float64](str); err != nil || val != 100.2 {
		t.Error(val, err)
	}

	str, err = ValueToString(200)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[int64](str); err != nil || val != 200 {
		t.Error(val, err)
	}

	str, err = ValueToString(true)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[bool](str); err != nil || val != true {
		t.Error(val, err)
	}

	st := testDataStruct{
		Field1: "f1",
		Field2: 1000,
		Field3: 2,
	}
	str, err = ValueToString(st)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[testDataStruct](str); err != nil || !reflect.DeepEqual(val, st) {
		t.Error(val, err)
	}

	if val, err := StringToValue[*testDataStruct](str); err != nil || !reflect.DeepEqual(*val, st) {
		t.Error(val, err)
	}

	mst := map[string]*testDataStruct{"one": &st}
	str, err = ValueToString(mst)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[map[string]*testDataStruct](str); err != nil || !reflect.DeepEqual(val, mst) {
		t.Error(val, err)
	}

	slice := []string{"a", "b", "c"}
	str, err = ValueToString(slice)
	if err != nil {
		t.FailNow()
	}

	if val, err := StringToValue[[]string](str); err != nil || !reflect.DeepEqual(val, slice) {
		t.Error(val, err)
	}

	// displace := &ws.P_DISPLACE{
	// 	OldIp: []byte("192.168.0.1"),
	// 	NewIp: []byte("192.168.0.100"),
	// 	Ts:    12345,
	// }
	// str, err = ValueToString(displace)
	// if err != nil {
	// 	t.FailNow()
	// }

	// if val, err := StringToValue[ws.P_DISPLACE](str); err != nil || val.Ts != displace.Ts {
	// 	t.Error(val, err)
	// }
}

type testDataStruct struct {
	Field1 string  `json:"field1" binding:"required"`
	Field2 int     `json:"field2" `
	Field3 float64 `json:"field3" `
}
