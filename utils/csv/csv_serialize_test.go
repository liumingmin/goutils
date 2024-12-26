package csv

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestSerializeDataTableToCsvFile(t *testing.T) {
	data := `id	name	age	remark
0	name0	0	remark0
1	name1	1	remark1
2	name2	2	remark2
3	name3	3	remark3
4	name4	4	remark4
5	name5	5	remark5
6	name6	6	remark6
7	name7	7	remark7
8	name8	8	remark8
9	name9	9	remark9
10	name10	10	remark10
11	name11	11	remark11
12	name12	12	remark12
13	name13	13	remark13
14	name14	14	remark14
15	name15	15	remark15
16	name16	16	remark16
17	name17	17	remark17
18	name18	18	remark18
19	name19	19	remark19`

	dt, err := ReadCsvToDataTable(context.Background(), strings.NewReader(data), '\t', []string{}, "id", []string{})
	if err != nil {
		t.Error(err)
	}

	var buffer bytes.Buffer
	SerializeDataTableToCsv(context.Background(), dt, &buffer, '\t')
	t.Log(buffer.String())
}
