package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"

	"github.com/F0903/pdf_downloader_uge5/models"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func isPdf(contentType string) bool {
	return contentType == "application/pdf"
}

func assertResponse(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code was not OK: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	// Download if content type is empty and just validate later
	if contentType == "" {
		return nil
	} else if !isPdf(contentType) {
		return errors.New("resource Content-Type is not PDF")
	}

	return nil
}

func downloadResourceWithProgress(url string, fullDownloadPath string, progressBar *mpb.Bar) error {
	// Request the web resource
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http GET request failed %w", err)
	}
	defer resp.Body.Close()

	if err := assertResponse(resp); err != nil {
		return err
	}

	// Create the download file
	file, err := os.Create(fullDownloadPath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	progressBar.SetTotal(resp.ContentLength, false)
	proxyReader := progressBar.ProxyReader(resp.Body)

	// Read from response and write to file whilst updating the progress bar
	if _, err := io.Copy(file, proxyReader); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func downloadReportWithProgress(report *models.Report, fullDownloadPath string, progressBar *mpb.Bar) *DownloadResult {
	if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
		progressBar.Abort(true)
		return NewDownloadResult(report, NewMissingDownloadState())
	}

	currentUrl := report.PrimaryDownloadLink
	onFallback := false
	var fullErr error // Start with an empty error to wrap around
	for {
		err := downloadResourceWithProgress(currentUrl, fullDownloadPath, progressBar)
		if err != nil {
			fullErr = errors.Join(fullErr, err)

			if onFallback {
				progressBar.Abort(true)
				return NewDownloadResult(
					report,
					NewFailedDownloadState(
						fmt.Errorf("all download links were broken: %w", fullErr),
					),
				)
			}

			currentUrl = report.FallbackDownloadLink
			onFallback = true
			continue
		}
		onFallback = false
		break
	}

	progressBar.SetTotal(progressBar.Current(), true)
	return NewDownloadResult(report, NewSuccededDownloadState(fullDownloadPath))
}

func DownloadReports(reports []*models.Report, directory string) []*DownloadResult {
	results := make([]*DownloadResult, len(reports))

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	var wg sync.WaitGroup

	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithAutoRefresh(),
	)

	for i, report := range reports {
		wg.Add(1)

		fileName := report.Id
		fullDownloadPath := path.Join(directory, fileName+".pdf")

		// It's important to create the progress bar here and not in the new thread or it will panic
		progressBar := p.AddBar(0,
			mpb.PrependDecorators(
				decor.Name(fileName, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.Percentage(),
			),
			mpb.BarRemoveOnComplete(),
		)

		go func() {
			defer wg.Done()
			result := downloadReportWithProgress(report, fullDownloadPath, progressBar)
			ValidateDownloadResult(result)
			results[i] = result
		}()
	}

	p.Wait()
	return results
}