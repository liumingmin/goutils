package utils

import (
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"

	"goutils/log"
)

func ReadCsvToData(filePath string, comma rune, colNames []string) (keys []string, resultData [][]string, err error) {
	rowsData, err := ParseCsv(filePath, comma)
	if err != nil {
		log.Error(context.Background(), "read file %s failed. error: %v", filePath, err)
		return
	}

	if len(rowsData) < 2 {
		log.Error(context.Background(), "read file %s is empty data file", filePath)
		return
	}

	log.Info(context.Background(), "%s len rowsData: %v", filePath, len(rowsData))

	header := rowsData[0]
	if len(colNames) == 0 {
		keys = header[:]
		resultData = rowsData[1:]
		return
	}

	fieldNameMap := make(map[string]int)
	for i, fieldName := range header {
		fieldNameMap[fieldName] = i
	}

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
	keys = colNames
	return
}

func ParseCsv(filePath string, comma rune) (records [][]string, err error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Error(context.Background(), "Open file %s failed. error: %v", filePath, err)
		return
	}
	defer csvFile.Close()

	bs, err := ioutil.ReadAll(csvFile)
	if err != nil {
		log.Error(context.Background(), "Read file %s failed. error: %v", filePath, err)
		return
	}

	fileContent, err := GBK2UTF8(bs)
	contentStr := string(fileContent)

	csvReader := csv.NewReader(strings.NewReader(contentStr))
	csvReader.Comma = comma
	csvReader.LazyQuotes = true
	records, err = csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		records = ParseCsvRaw(contentStr)
		err = nil
		log.Error(context.Background(), "Read file %s failed. error: %v, try parse raw: %v", filePath, err, len(records))
	}
	return
}

func ParseCsvRaw(content string) (records [][]string) {
	rowStrs := strings.Split(content, "\n")
	for _, rowStr := range rowStrs {
		row := strings.Split(rowStr, "\t")
		records = append(records, row)
	}
	return
}
