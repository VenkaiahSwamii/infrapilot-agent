package reports

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct {
	BaseDir string
}

func NewPDFGenerator(baseDir string) *PDFGenerator {
	return &PDFGenerator{BaseDir: baseDir}
}

func (p *PDFGenerator) Generate(report *Report, data map[string]interface{}) (string, error) {
	if err := os.MkdirAll(p.BaseDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s_%s.pdf", report.Type, report.ID.String())
	filePath := filepath.Join(p.BaseDir, filename)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle(report.Name, true)
	pdf.SetAuthor("InfraPilot Enterprise", true)
	pdf.SetCreator("InfraPilot Reporting Engine", true)

	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, report.Name)

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Generated: %s", report.GeneratedAt.Format(time.RFC1123)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Type: %s", strings.ReplaceAll(string(report.Type), "_", " ")))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Format: %s", strings.ToUpper(report.Format)))
	pdf.Ln(10)

	if title, ok := data["title"].(string); ok && title != "" {
		pdf.SetFont("Arial", "B", 13)
		pdf.Cell(0, 8, title)
		pdf.Ln(10)
	}

	if summary, ok := data["summary"].(string); ok && summary != "" {
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, summary, "", "", false)
		pdf.Ln(5)
	}

	if tables, ok := data["tables"].([]TableData); ok {
		for _, table := range tables {
			renderPDFTable(pdf, table)
			pdf.Ln(6)
		}
	}

	// Footer
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()))
	pdf.Ln(-1)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := pdf.Output(outFile); err != nil {
		return "", err
	}

	return filename, nil
}

func renderPDFTable(pdf *gofpdf.Fpdf, table TableData) {
	if len(table.Headers) == 0 || len(table.Rows) == 0 {
		return
	}

	colWidth := 180.0 / float64(len(table.Headers))
	if colWidth < 20 {
		colWidth = 20
	}

	pdf.SetFont("Arial", "B", 9)
	for _, header := range table.Headers {
		pdf.CellFormat(colWidth, 6, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	for _, row := range table.Rows {
		cells := 0
		for range row {
			if cells < len(table.Headers) {
				val := ""
				if cells < len(row) {
					val = fmt.Sprintf("%v", row[cells])
				}
				pdf.CellFormat(colWidth, 6, truncate(val, 32), "1", 0, "L", false, 0, "")
			}
			cells++
		}
		pdf.Ln(-1)
	}
}

type TableData struct {
	Title   string
	Headers []string
	Rows    [][]interface{}
}

func truncate(text string, max int) string {
	text = strings.ReplaceAll(text, "\n", " ")
	if len(text) <= max {
		return text
	}
	if max > 3 {
		return text[:max-3] + "..."
	}
	return text[:max]
}

func (p *PDFGenerator) Cleanup(filename string) error {
	path := filepath.Join(p.BaseDir, filename)
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}
