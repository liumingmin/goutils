package csv

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/liumingmin/goutils/container"
	"github.com/liumingmin/goutils/log"
)

func SerializeDataTableToCsvFile(ctx context.Context, dataTable *container.DataTable, filePath string, comma rune) error {
	fWriter, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error(ctx, "Open file %s failed. error: %v", filePath, err)
		return err
	}
	defer fWriter.Close()

	return SerializeDataTableToCsv(ctx, dataTable, fWriter, comma)
}

func SerializeDataTableToCsv(ctx context.Context, dataTable *container.DataTable, writer io.Writer, comma rune) error {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = comma

	csvWriter.Write(dataTable.Cols())
	for _, row := range dataTable.Rows() {
		csvWriter.Write(row.Data())
	}
	csvWriter.Flush()
	return nil
}
