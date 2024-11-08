package excel

import (
	"fmt"
	"path"
	"strconv"

	"github.com/F0903/pdf_downloader_uge5/downloader/report_downloader"
	"github.com/xuri/excelize/v2"
)

const sheetName = "Metadata"

func setMainSheetWidths(f *excelize.File) error {
	// Set Name column width
	err := f.SetColWidth(sheetName, "B", "B", 50)
	if err != nil {
		return fmt.Errorf("could not set sheet B column width: %w", err)
	}

	// Set download urls column width
	err = f.SetColWidth(sheetName, "C", "D", 150)
	if err != nil {
		return fmt.Errorf("could not set sheet C-D column width: %w", err)
	}

	// Set DownloadState column width
	err = f.SetColWidth(sheetName, "E", "E", 200)
	if err != nil {
		return fmt.Errorf("could not set sheet E column widths: %w", err)
	}

	return nil
}

func writeHeader(f *excelize.File) error {
	err := f.SetSheetRow(sheetName, "A1", &[]interface{}{"ID", "Name", "PrimaryDownloadURL", "FallbackDownloadURL", "DownloadState"})
	if err != nil {
		return fmt.Errorf("could not set sheet row: %w", err)
	}

	boldTextStyle, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		return fmt.Errorf("could not create bold text style: %w", err)
	}

	err = f.SetRowStyle(sheetName, 1, 1, boldTextStyle)
	if err != nil {
		return fmt.Errorf("could not set header row style: %w", err)
	}

	return nil
}

func writeResultsToRows(f *excelize.File, results []*report_downloader.ReportDownloadResult) {
	for i, result := range results {
		// We add 2 because Excel starts counting at 1, and our header is already at A1
		index := "A" + strconv.Itoa(i+2)
		report := result.AssociatedReport
		downloadState := result.State
		err := f.SetSheetRow(
			sheetName,
			index,
			&[]interface{}{
				report.Id,
				report.Name,
				report.PrimaryDownloadLink,
				report.FallbackDownloadLink,
				downloadState.StringNoNewLines(),
			},
		)
		if err != nil {
			f.SetCellValue(sheetName, index, fmt.Sprintf("Error when writing row: %v", err))
		}
	}
}

func WriteDownloadResults(results []*report_downloader.ReportDownloadResult, directory string) error {
	fullOutputPath := path.Join(directory, "metadata.xlsx")
	fmt.Printf("Writing download result metadata to '%s'...\n", fullOutputPath)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Could not close download results metadata file!\n%v", err)
		}
	}()

	// Set our main sheet
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return fmt.Errorf("could not rename sheet on metadata spreadsheet: %w", err)
	}

	if err := writeHeader(f); err != nil {
		return fmt.Errorf("could not write header: %w", err)
	}

	if err := setMainSheetWidths(f); err != nil {
		return fmt.Errorf("could not set column widths: %w", err)
	}

	writeResultsToRows(f, results)

	if err := f.SaveAs(fullOutputPath); err != nil {
		return fmt.Errorf("could not save download result metadata spreadsheet: %w", err)
	}

	return nil
}
