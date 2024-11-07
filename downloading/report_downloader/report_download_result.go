package report_downloader

import (
	"fmt"

	"github.com/F0903/pdf_downloader_uge5/models"
)

type ReportDownloadResult struct {
	AssociatedReport *models.Report
	State            *DownloadState
}

func (result *ReportDownloadResult) String() string {
	resultReport := result.AssociatedReport
	resultState := result.State

	return fmt.Sprintf("[%s | %s | (%s | %s)] = %s", resultReport.Id, resultReport.Name, resultReport.PrimaryDownloadLink, resultReport.FallbackDownloadLink, resultState.String())
}

func NewDownloadResult(associatedReport *models.Report, state *DownloadState) *ReportDownloadResult {
	return &ReportDownloadResult{
		associatedReport, state,
	}
}
