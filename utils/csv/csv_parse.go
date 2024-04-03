package csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/liumingmin/goutils/container"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
)

var (
	ErrCsvIsEmpty = errors.New("csv is empty")
)

func ReadCsvFileToDataTable(ctx context.Context, filePath string, comma rune, colNames []string, pkCol string,
	indexes []string) (*container.DataTable, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		log.Error(ctx, "Open file %s failed. error: %v", filePath, err)
		return nil, err
	}
	defer reader.Close()

	return ReadCsvToDataTable(ctx, reader, comma, colNames, pkCol, indexes)
}

func ReadCsvToDataTable(ctx context.Context, reader io.Reader, comma rune, colNames []string, pkCol string,
	indexes []string) (*container.DataTable, error) {
	keys, rowsData, err := ReadCsvToData(ctx, reader, comma, colNames)
	if err != nil {
		return nil, err
	}

	if pkCol == "" {
		pkCol = keys[0]
	}

	log.Debug(ctx, "csv data keys: %v, %d", keys, len(keys))

	dataTable := container.NewDataTable(keys, pkCol, indexes, len(rowsData))
	dataTable.PushAll(rowsData)

	return dataTable, nil
}

func ReadCsvToData(ctx context.Context, reader io.Reader, comma rune, colNames []string) ([]string, [][]string, error) {
	rowsData, err := ParseCsv(ctx, reader, comma)
	if err != nil {
		log.Error(ctx, "read bytes failed. error: %v", err)
		return nil, nil, err
	}

	if len(rowsData) == 0 {
		return nil, nil, ErrCsvIsEmpty
	}

	log.Debug(ctx, "raw data len rowsData: %v", len(rowsData))

	header := rowsData[0]
	if len(colNames) == 0 {
		return header, rowsData[1:], nil
	}

	fieldNameMap := make(map[string]int)
	for i, fieldName := range header {
		fieldNameMap[fieldName] = i
	}

	resultData := make([][]string, 0, len(rowsData)-1)

	bodyData := rowsData[1:]
	for _, row := range bodyData {
		rowLen := len(row)
		newRow := make([]string, 0, len(colNames))
		for _, colName := range colNames {
			if idx, ok := fieldNameMap[colName]; ok && idx < rowLen {
				newRow = append(newRow, row[idx])
			}
		}
		resultData = append(resultData, newRow)
	}
	return colNames, resultData, nil
}

func ParseCsv(ctx context.Context, reader io.Reader, comma rune) (records [][]string, err error) {
	bs, err := io.ReadAll(reader)
	if err != nil {
		log.Error(ctx, "Read bytes failed, error: %v", err)
		return
	}

	var fileContent []byte
	if utils.IsUtf8(bs) {
		fileContent = bs
	} else if utils.IsGBK(bs) {
		fileContent, err = utils.GBK2UTF8(bs)
		if err != nil {
			log.Error(ctx, "GBK2UTF8 failed, error: %v", err)
			return
		}
	} else {
		err = errors.New("unsupport encoding")
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(fileContent))
	csvReader.Comma = comma
	csvReader.LazyQuotes = true
	records, err = csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		records = ParseCsvRaw(ctx, string(fileContent))
		err = nil
		log.Error(ctx, "Read bytes failed. error: %v, try parse raw: %v", err, len(records))
	}
	return
}

func ParseCsvRaw(ctx context.Context, content string) (records [][]string) {
	rowStrs := strings.Split(content, "\n")
	for _, rowStr := range rowStrs {
		row := strings.Split(rowStr, "\t")
		records = append(records, row)
	}
	return
}
