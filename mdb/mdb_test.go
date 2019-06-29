package mdb

import "testing"

func TestReadCsvToDataTable(t *testing.T) {
	dt, err := ReadCsvToDataTable(`xxxxxxxxxxxxxxx`, '\t',
		[]string{"xx", "xx", "xx", "xx"}, "xxx", []string{"xxx"})
	if err != nil {
		t.Log(err)
		return
	}
	//for _, r := range dt.Rows() {
	//	t.Log(r.row)
	//}

	rs := dt.RowsBy("xxx", "869")
	for _, r := range rs {
		t.Log(r.row)
	}

	t.Log(dt.Row("17"))
}
