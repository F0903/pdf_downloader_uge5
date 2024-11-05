package downloader

import "github.com/F0903/pdf_downloader_uge5/models"

// A literal "report" of the download.
type DownloadReport struct {
	associatedReport *models.Report
	state            DownloadState
}
