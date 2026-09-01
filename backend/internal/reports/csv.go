package reports

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CSVGenerator struct {
	BaseDir string
}

func NewCSVGenerator(baseDir string) *CSVGenerator {
	return &CSVGenerator{BaseDir: baseDir}
}

func (c *CSVGenerator) Generate(report *Report, data map[string]interface{}) (string, error) {
	if err := os.MkdirAll(c.BaseDir, 0755); err != nil {
		return "", err
	}

	filename := string(report.Type) + "_" + report.ID.String() + ".csv"
	filePath := filepath.Join(c.BaseDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write metadata header
	writer.Write([]string{"Report: " + report.Name})
	writer.Write([]string{"Generated: " + report.GeneratedAt.Format(time.RFC1123)})
	writer.Write([]string{"Type: " + string(report.Type)})
	writer.Write([]string{"Format: " + strings.ToUpper(report.Format)})
	writer.Write([]string{})

	// Write tables
	if tables, ok := data["tables"].([]TableData); ok {
		for _, table := range tables {
			if title := table.Title; title != "" {
				writer.Write([]string{"=== " + title + " ==="})
				writer.Write([]string{})
			}

			if len(table.Headers) > 0 {
				if err := writer.Write(table.Headers); err != nil {
					return "", err
				}
			}

			for _, row := range table.Rows {
				record := make([]string, len(table.Headers))
				for i, val := range row {
					if i < len(table.Headers) {
						record[i] = fmt.Sprintf("%v", val)
					}
				}
				if err := writer.Write(record); err != nil {
					return "", err
				}
			}
			writer.Write([]string{})
		}
	}

	return filename, nil
}

func (c *CSVGenerator) Cleanup(filename string) error {
	path := filepath.Join(c.BaseDir, filename)
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}
