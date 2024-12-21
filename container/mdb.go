package container

import (
	"strconv"
)

type DataRow struct {
	row    []string
	colMap map[string]int
}

func (r *DataRow) String(fieldName string) string {
	if colIdx, ok := r.colMap[fieldName]; ok {
		if colIdx >= len(r.row) {
			return ""
		}

		return r.row[colIdx]
	}
	return ""
}

func (r *DataRow) Int64(fieldName string) int64 {
	i, _ := strconv.ParseInt(r.String(fieldName), 10, 64)
	return i
}

func (r *DataRow) UInt64(fieldName string) uint64 {
	i, _ := strconv.ParseUint(r.String(fieldName), 10, 64)
	return i
}

func (r *DataRow) SetString(fieldName, fieldValue string) {
	if colIdx, ok := r.colMap[fieldName]; ok {
		if colIdx >= len(r.row) {
			return
		}

		r.row[colIdx] = fieldValue
	}
}

func (r *DataRow) Data() []string {
	return r.row
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

// meta info
func (t *DataTable) Cols() []string {
	return t.cols[:]
}

func (t *DataTable) PkCol() string {
	return t.pkCol
}

func (t *DataTable) Indexes() []string {
	return t.indexes[:]
}

// data info
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

func (t *DataTable) RowsByPredicate(predicate func(*DataRow) bool) []*DataRow {
	result := make([]*DataRow, 0)
	for _, row := range t.rows {
		if predicate(row) {
			result = append(result, row)
		}
	}
	return result
}

func (t *DataTable) RowsByIndexPredicate(indexName, indexValue string, predicate func(*DataRow) bool) []*DataRow {
	rows := t.RowsBy(indexName, indexValue)
	if len(rows) == 0 {
		return rows
	}

	result := make([]*DataRow, 0)
	for _, row := range rows {
		if predicate(row) {
			result = append(result, row)
		}
	}
	return result
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

func NewDataSet() *DataSet {
	return &DataSet{}
}

func (s *DataSet) AddTable(tName string, dt *DataTable) {
	(*s)[tName] = dt
}

func (s *DataSet) Table(tName string) *DataTable {
	return (*s)[tName]
}
