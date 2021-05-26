package mdb

import (
	"context"
	"strconv"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
)

type DataRow struct {
	row    []string
	colMap map[string]int
}

func (r *DataRow) String(fieldName string) string {
	if colIdx, ok := r.colMap[fieldName]; ok {
		return r.row[colIdx]
	}
	return ""
}

func (r DataRow) Int(fieldName string) int {
	i, _ := strconv.Atoi(r.String(fieldName))
	return i
}

type tableIndex map[string][]*DataRow

type DataTable struct {
	//define
	cols    []string
	pkCol   string
	indexes []string

	//row share
	colMap map[string]int

	//data
	rows []*DataRow

	//index
	pkMap      map[string]*DataRow
	indexesMap map[string]tableIndex
}

func NewDataTable(cols []string, pkCol string, indexes []string, initCap int) *DataTable {
	colMap := make(map[string]int, len(cols))
	for i, col := range cols {
		colMap[col] = i
	}

	indexesMap := make(map[string]tableIndex, len(indexes))
	if len(indexes) > 0 {
		for _, indexName := range indexes {
			indexesMap[indexName] = tableIndex{}
		}
	}

	dt := &DataTable{
		cols:    cols,
		pkCol:   pkCol,
		indexes: indexes,
		colMap:  colMap,

		rows:       make([]*DataRow, 0, initCap),
		pkMap:      make(map[string]*DataRow, initCap),
		indexesMap: indexesMap,
	}
	return dt
}

func (t *DataTable) Row(pk string) *DataRow {
	return t.pkMap[pk]
}

func (t *DataTable) Rows() (rows []*DataRow) {
	return t.rows
}

func (t *DataTable) RowsBy(indexName, indexValue string) []*DataRow {
	indexMap, ok := t.indexesMap[indexName]
	if !ok {
		return nil
	}

	return indexMap[indexValue]
}

func (t *DataTable) PkString(row *DataRow) string {
	return row.String(t.pkCol)
}

func (t *DataTable) PkInt(row *DataRow) int {
	i, _ := strconv.Atoi(t.PkString(row))
	return i
}

func (t *DataTable) Push(row []string) {
	dr := &DataRow{row: row, colMap: t.colMap}
	t.rows = append(t.rows, dr)

	t.pkMap[t.PkString(dr)] = dr
	for indexName, indexMap := range t.indexesMap {
		indexValue := dr.String(indexName)
		indexMap[indexValue] = append(indexMap[indexValue], dr)
	}
}

func (t *DataTable) PushAll(rows [][]string) {
	for _, row := range rows {
		t.Push(row)
	}
}

type DataSet map[string]*DataTable

func (s DataSet) Table(tName string) *DataTable {
	return s[tName]
}

func ReadCsvToDataTable(filePath string, comma rune, colNames []string, pkCol string, indexes []string) (dataTable *DataTable, err error) {
	keys, rowsData, err := utils.ReadCsvToData(filePath, comma, colNames)
	if err != nil {
		return
	}

	if pkCol == "" {
		pkCol = keys[0]
	}

	log.Info(context.Background(), "%s keys: %v, %d", filePath, keys, len(keys))

	dataTable = NewDataTable(keys, pkCol, indexes, len(rowsData))
	dataTable.PushAll(rowsData[1:])

	return
}
