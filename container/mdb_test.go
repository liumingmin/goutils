package container

import (
	"reflect"
	"strconv"
	"testing"
)

var testDs *DataSet
var testDt *DataTable

func TestDataTable(t *testing.T) {
	if len(testDt.Rows()) != 10 {
		t.Error(len(testDt.Rows()))
	}

	if !reflect.DeepEqual(testDt.Cols(), []string{"id", "code", "name"}) {
		t.Error(testDt.Cols())
	}

	if testDt.PkCol() != "id" {
		t.Error(testDt.PkCol())
	}

	if !reflect.DeepEqual(testDt.Indexes(), []string{"code"}) {
		t.Error(testDt.Indexes())
	}

	if testDt.PkString(testDt.Row("9")) != "9" {
		t.Error(testDt.PkString(testDt.Row("9")))
	}

	if testDt.PkInt(testDt.Row("8")) != 8 {
		t.Error(testDt.PkInt(testDt.Row("8")))
	}

	if testDt.Row("9").Int64("id") != 9 {
		t.Error(testDt.Row("9").Int64("id"))
	}

	if testDt.Row("9").UInt64("id") != 9 {
		t.Error(testDt.Row("9").Int64("id"))
	}

	if testDt.Row("9").String("code") != "C9" {
		t.Error(testDt.Row("9").String("code"))
	}

	if testDt.Row("9").String("nop") != "" {
		t.Error(testDt.Row("9").String("nop"))
	}

	if !reflect.DeepEqual(testDt.Row("2").Data(), []string{"2", "C2", "N2"}) {
		t.Error(testDt.Row("2").Data())
	}

	if !reflect.DeepEqual(testDt.RowsBy("code", "C2")[0].Data(), []string{"2", "C2", "N2"}) {
		t.Error(testDt.RowsBy("code", "C2")[0].Data())
	}

	if !reflect.DeepEqual(testDt.RowsByPredicate(func(dr *DataRow) bool { return dr.String("name") == "N4" })[0].Data(), []string{"4", "C4", "N4"}) {
		t.Error("RowsByPredicate")
	}

	testDt.Push([]string{"2", "C2", "N3"})

	if !reflect.DeepEqual(testDt.RowsByIndexPredicate("code", "C2", func(dr *DataRow) bool { return dr.String("name") == "N3" })[0].Data(), []string{"2", "C2", "N3"}) {
		t.Error("RowsByIndexPredicate")
	}

	testDt.Row("9").SetString("code", "new9")
	if testDt.Row("9").String("code") != "new9" {
		t.Error(testDt.Row("9"))
	}
}

func TestDataSet(t *testing.T) {
	if testDs.Table("testDt") != testDt {
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	testDt = NewDataTable([]string{"id", "code", "name"}, "id", []string{"code"}, 10)

	rows := make([][]string, 0)
	for i := 0; i < 10; i++ {
		row := []string{strconv.Itoa(i), "C" + strconv.Itoa(i), "N" + strconv.Itoa(i)}
		rows = append(rows, row)
	}
	testDt.PushAll(rows)

	testDs = NewDataSet()
	testDs.AddTable("testDt", testDt)

	m.Run()
}
