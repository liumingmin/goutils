package csv

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/liumingmin/goutils/container"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_csv")
var testCsvFilePath = "goutils.csv"

func TestReadCsvToDataTable(t *testing.T) {
	dt, err := ReadCsvFileToDataTable(context.Background(), filepath.Join(testTempDirPath, testCsvFilePath), ',',
		[]string{"id", "name", "age", "remark"}, "id", []string{"name"})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(dt.Row("10").Data(), []string{"10", "name10", "10", "remark10"}) {
		t.FailNow()
	}

	if !reflect.DeepEqual(dt.RowsBy("name", "name10")[0].Data(), []string{"10", "name10", "10", "remark10"}) {
		t.FailNow()
	}
}

func TestParseCsvRaw(t *testing.T) {
	records := ParseCsvRaw(context.Background(),
		`id	name	age	remark
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
19	name19	19	remark19`)

	dt := container.NewDataTable(records[0], "id", []string{"name"}, 20)
	dt.PushAll(records[1:])

	if !reflect.DeepEqual(dt.Row("10").Data(), []string{"10", "name10", "10", "remark10"}) {
		t.FailNow()
	}

	if !reflect.DeepEqual(dt.RowsBy("name", "name10")[0].Data(), []string{"10", "name10", "10", "remark10"}) {
		t.FailNow()
	}
}

func TestParseCsvGBK(t *testing.T) {
	u8 := `id	name	age	remark
0	name0	0	描述0
1	name1	1	描述1
2	name2	2	描述2
3	name3	3	描述3
4	name4	4	描述4
5	name5	5	描述5
6	name6	6	描述6
7	name7	7	描述7
8	name8	8	描述8
9	name9	9	描述9
10	name10	10	描述10
11	name11	11	描述11
12	name12	12	描述12
13	name13	13	描述13
14	name14	14	描述14
15	name15	15	描述15
16	name16	16	描述16
17	name17	17	描述17
18	name18	18	描述18
19	name19	19	描述19`

	reader := transform.NewReader(bytes.NewReader([]byte(u8)), simplifiedchinese.GB18030.NewEncoder())

	dt, err := ReadCsvToDataTable(context.Background(), reader, '\t',
		[]string{"id", "name", "remark"}, "", []string{"name"}) //pk default cols[0]
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(dt.Row("10").Data(), []string{"10", "name10", "描述10"}) {

		t.Error(dt.Row("10"))
	}

	if !reflect.DeepEqual(dt.RowsBy("name", "name10")[0].Data(), []string{"10", "name10", "描述10"}) {
		t.FailNow()
	}
}

func TestParseShortCsv(t *testing.T) {
	data := `id	name	age	remark`
	dt, err := ReadCsvToDataTable(context.Background(), strings.NewReader(data), '\t',
		[]string{"id", "name", "remark"}, "", []string{"name"}) //pk default cols[0]
	if err != nil {
		t.Error(err)
	}

	if dt == nil {
		t.FailNow()
	}

	if len(dt.Rows()) != 0 {
		t.Error(dt.Rows())
	}

	dt, err = ReadCsvToDataTable(context.Background(), strings.NewReader(""), '\t',
		[]string{"id", "name", "remark"}, "", []string{"name"}) //pk default cols[0]
	if err != ErrCsvIsEmpty {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	os.MkdirAll(testTempDirPath, 0666)

	csvFilePath := filepath.Join(testTempDirPath, testCsvFilePath)
	file, err := os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	file.WriteString("id,name,age,remark\n")
	for l1 := 0; l1 < 20; l1++ {
		file.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", l1, "name"+fmt.Sprint(l1), l1, "remark"+fmt.Sprint(l1)))
	}

	file.Close()

	m.Run()

	os.RemoveAll(testTempDirPath)
}
