package downloader

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func validateDownloadResult(filePath string) error {
	return api.ValidateFile(filePath, model.NewDefaultConfiguration())
}

func ValidateDownloadResults(reports []*DownloadResult) error {

}
