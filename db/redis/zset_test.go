package redis

import (
	"context"
	"strings"
	"testing"

	"github.com/liumingmin/goutils/utils/csv"
)

func TestZDescartes(t *testing.T) {
	InitRedises()
	rds := Get("rdscdb")
	ctx := context.Background()
	dimValues := [][]string{{"dim1a", "dim1b"}, {"dim2a", "dim2b", "dim2c", "dim2d"}, {"dim3a", "dim3b", "dim3c"}}

	dt, err := csv.ReadCsvToDataTable(ctx, "data.csv", ',',
		[]string{"id", "name", "createtime", "dim1", "dim2", "dim3", "member"}, "id", []string{})
	if err != nil {
		t.Error(err)
	}

	err = ZDescartes(ctx, rds, dimValues, func(strs []string) (string, map[string]int64) {
		dimData := make(map[string]int64)
		for _, row := range dt.Rows() {
			if row.String("dim1") == strs[0] &&
				row.String("dim2") == strs[1] &&
				row.String("dim3") == strs[2] {
				dimData[row.String("member")] = row.Int64("createtime")
			}
		}
		return "rds" + strings.Join(strs, "-"), dimData
	}, 1000, 30)

	if err != nil {
		t.Error(err)
	}
}
