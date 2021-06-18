package utils

import (
	"context"
	"testing"
)

func TestReadCsvToDataTable(t *testing.T) {
	dt, err := ReadCsvToDataTable(context.Background(), `goutils.log`, '\t',
		[]string{"xx", "xx", "xx", "xx"}, "xxx", []string{"xxx"})
	if err != nil {
		t.Log(err)
		return
	}
	for _, r := range dt.Rows() {
		t.Log(r.Data())
	}

	rs := dt.RowsBy("xxx", "869")
	for _, r := range rs {
		t.Log(r.Data())
	}

	t.Log(dt.Row("17"))
}
