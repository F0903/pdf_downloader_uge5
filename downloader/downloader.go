package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/F0903/pdf_downloader_uge5/models"
)

func downloadResource(url string, filename string, directory string) error {
	urlExt := path.Ext(url)
	fullDownloadPath := path.Join(directory, filename+urlExt)

	fmt.Printf("Downloading %s...", fullDownloadPath)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http GET request failed %w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(fullDownloadPath)
	if err != nil {
		return fmt.Errorf("could not create file %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy web response to file %w", err)
	}

	fmt.Printf("Downloaded %s", fullDownloadPath)
	return nil
}

func downloadReport(report *models.Report, directory string) *DownloadReport {
	if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
		return errors.New("no download links were supplied")
	}

	currentUrl := report.PrimaryDownloadLink
	onFallback := false
	for {
		err := downloadResource(currentUrl, report.Id, directory)
		if err != nil {
			if onFallback {
				return errors.New("all download links were broken")
			}

			currentUrl = report.FallbackDownloadLink
			onFallback = true
			continue
		}
		onFallback = false
		break
	}

	return nil
}

func DownloadReports(reports []*models.Report, directory string) {
	for _, report := range reports {
		downloadReport(report, directory)
	}
}
