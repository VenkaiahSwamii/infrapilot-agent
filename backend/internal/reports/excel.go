package reports

import (
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelGenerator struct {
	BaseDir string
}

func NewExcelGenerator(baseDir string) *ExcelGenerator {
	return &ExcelGenerator{BaseDir: baseDir}
}

func (e *ExcelGenerator) Generate(report *Report, data map[string]interface{}) (string, error) {
	if err := os.MkdirAll(e.BaseDir, 0755); err != nil {
		return "", err
	}

	filename := string(report.Type) + "_" + report.ID.String() + ".xlsx"
	filePath := filepath.Join(e.BaseDir, filename)

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", report.Name)
	f.SetCellValue("Sheet1", "A2", "Generated: "+report.GeneratedAt.Format(time.RFC1123))
	f.SetCellValue("Sheet1", "A3", "Type: "+report.Type)
	f.SetCellValue("Sheet1", "A4", "Format: "+report.Format)

	row := 6
	if title, ok := data["title"].(string); ok && title != "" {
		f.SetCellValue("Sheet1", "A"+numberToLetters(row), title)
		row++
	}

	if tables, ok := data["tables"].([]TableData); ok {
		for _, table := range tables {
			if len(table.Headers) > 0 {
				for c, header := range table.Headers {
					cell, _ := excelize.CoordinatesToCellName(c+1, row)
					f.SetCellValue("Sheet1", cell, header)
				}
				row++
			}

			for _, rowData := range table.Rows {
				for c, val := range rowData {
					if c < len(table.Headers) {
						cell, _ := excelize.CoordinatesToCellName(c+1, row)
						f.SetCellValue("Sheet1", cell, val)
					}
				}
				row++
			}
			row++
		}
	}

	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filename, nil
}

func numberToLetters(n int) string {
	result := ""
	for n > 0 {
		n--
		result = string(rune('A'+n%26)) + result
		n /= 26
	}
	return result
}

func (e *ExcelGenerator) Cleanup(filename string) error {
	path := filepath.Join(e.BaseDir, filename)
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}
