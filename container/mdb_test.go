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
		t.FailNow()
	}

	if testDt.PkString(testDt.Row("9")) != "9" {
		t.FailNow()
	}

	if testDt.PkInt(testDt.Row("8")) != 8 {
		t.FailNow()
	}

	if reflect.DeepEqual(testDt.Row("2"), []string{"2", "C2", "N2"}) {
		t.FailNow()
	}

	if reflect.DeepEqual(testDt.RowsBy("code", "C2")[0], []string{"2", "C2", "N2"}) {
		t.FailNow()
	}

	if reflect.DeepEqual(testDt.RowsByPredicate(func(dr *DataRow) bool { return dr.String("name") == "N4" })[0], []string{"4", "C4", "N4"}) {
		t.FailNow()
	}

	testDt.Push([]string{"2", "C2", "N3"})

	if reflect.DeepEqual(testDt.RowsByIndexPredicate("code", "C2", func(dr *DataRow) bool { return dr.String("name") == "N3" })[0], []string{"2", "C2", "N2"}) {
		t.FailNow()
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
