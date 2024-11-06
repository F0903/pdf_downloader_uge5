package downloader

import (
	"fmt"

	"github.com/F0903/pdf_downloader_uge5/models"
)

// A literal "report" of the download.
type DownloadResult struct {
	AssociatedReport *models.Report
	State            *DownloadState
}

func (result *DownloadResult) String() string {
	resultReport := result.AssociatedReport
	resultState := result.State

	return fmt.Sprintf("[%s | %s | (%s | %s)] = %s", resultReport.Id, resultReport.Name, resultReport.PrimaryDownloadLink, resultReport.FallbackDownloadLink, resultState.String())
}

func NewDownloadResult(associatedReport *models.Report, state *DownloadState) *DownloadResult {
	return &DownloadResult{
		associatedReport, state,
	}
}
