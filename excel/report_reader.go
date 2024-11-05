package excel

import (
	"errors"
	"fmt"

	"github.com/F0903/pdf_downloader_uge5/models"
	"github.com/xuri/excelize/v2"
)

func createReportFromRow(row []string) *models.Report {
	report := &models.Report{}
	for colIndex, colCell := range row {
		// Check if the column index maps to a column we want
		switch ColumnMappings[colIndex] {
		case IdColumn:
			report.Id = colCell
		case NameColumn:
			report.Name = colCell
		case PrimaryDownloadColumn:
			report.PrimaryDownloadLink = colCell
		case SecondaryDownloadColumn:
			report.FallbackDownloadLink = colCell
		}
	}

	return report
}

func ReadReports(path string) ([]*models.Report, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel spreadsheet!\n%w", err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// We just need the first sheet
	firstSheetName := f.GetSheetName(0)

	rows, err := f.Rows(firstSheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows in spreadsheet!\n%w", err)
	}
	defer rows.Close()

	// We start by skipping the header row.
	if !rows.Next() {
		return nil, errors.New("empty spreadsheet")
	}

	reports := make([]*models.Report, 0)
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get single row in spreadsheet!\n%w", err)
		}

		report := createReportFromRow(row)
		reports = append(reports, report)
	}

	return reports, nil
}
