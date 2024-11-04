package main

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Report struct {
	id                      string
	name                    string
	primary_download_link   string
	secondary_download_link string
}

func createReportFromRow(row []string) *Report {
	report := &Report{}
	for colIndex, colCell := range row {
		// Check if the column index maps to a column we want
		switch ColumnMappings[colIndex] {
		case IdColumn:
			report.id = colCell
		case NameColumn:
			report.name = colCell
		case PrimaryDownloadColumn:
			report.primary_download_link = colCell
		case SecondaryDownloadColumn:
			report.secondary_download_link = colCell
		}
	}

	return report
}

func ReadReports(path string) ([]*Report, error) {
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

	reports := make([]*Report, 0)
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get single row in spreadsheet!\n%w", err)
		}

		fmt.Println()

		report := createReportFromRow(row)
		reports = append(reports, report)
	}

	return reports, nil
}
