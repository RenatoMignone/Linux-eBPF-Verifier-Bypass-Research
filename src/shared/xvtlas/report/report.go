package report

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"fmt"
)

type CSVRow struct {
	ID                 int
	Filename           string
	LoadParameters     string
	VulnNumber         string
	Compiled           bool
	Verified           bool
	Loaded             bool
	LoadOutput         string
	KernelVersion      string
}

func ExportCSV(rows []CSVRow, exportPath string) {
	file, _ := os.Create(filepath.Join(exportPath, "report.csv"))
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "filename", "load parameters", "vuln_number", "compiled", "verified", "loaded", "load_output"})

	for i, row := range rows {
		writer.Write([]string{
			fmt.Sprint(i),
			row.Filename,
			row.LoadParameters,
			row.VulnNumber,
			fmt.Sprint(row.Compiled),
			fmt.Sprint(row.Verified),
			fmt.Sprint(row.Loaded),
			row.LoadOutput,
		})
	}
}

