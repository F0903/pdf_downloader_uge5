package report_downloader

import (
	"fmt"

	"github.com/F0903/pdf_downloader_uge5/downloader/report_downloader/report_download_state"
	"github.com/F0903/pdf_downloader_uge5/models"
)

// A literal "report" of the download.
type ReportDownloadResult struct {
	AssociatedReport *models.Report
	State            *report_download_state.ReportDownloadState
}

func (result *ReportDownloadResult) String() string {
	resultReport := result.AssociatedReport
	resultState := result.State

	return fmt.Sprintf("[%s | %s | (%s | %s)] = %s", resultReport.Id, resultReport.Name, resultReport.PrimaryDownloadLink, resultReport.FallbackDownloadLink, resultState.String())
}

func NewReportDownloadResult(associatedReport *models.Report, state *report_download_state.ReportDownloadState) *ReportDownloadResult {
	return &ReportDownloadResult{
		associatedReport, state,
	}
}
