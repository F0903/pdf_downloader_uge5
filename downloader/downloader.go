package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/F0903/pdf_downloader_uge5/models"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func isPdf(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return contentType == "application/pdf"
}

func downloadResourceWithProgress(url string, fullDownloadPath string, progressBar *mpb.Bar) error {
	// Request the web resource
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http GET request failed %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code was not OK: %d", resp.StatusCode)
	}

	if !isPdf(resp) {
		return errors.New("resource is not a PDF")
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

func downloadReportWithProgress(report *models.Report, fullDownloadPath string, progressBar *mpb.Bar) DownloadResult {
	if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
		progressBar.Abort(true)
		return NewDownloadResult(report, NewDownloadState(MissingURLs, nil))
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
					NewDownloadState(
						Error,
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
	return NewDownloadResult(report, NewDownloadState(Done, nil))
}

func DownloadReports(reports []*models.Report, directory string) []DownloadResult {
	results := make([]DownloadResult, len(reports))

	var wg sync.WaitGroup

	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithAutoRefresh(),
	)

	for i, report := range reports {
		wg.Add(1)

		currentUrl := report.PrimaryDownloadLink
		fileName := report.Id
		urlExt := path.Ext(currentUrl)
		fullDownloadPath := path.Join(directory, fileName+urlExt)

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

		// Explicitly pass variables or it will get a little funky
		go func(index int, report *models.Report, progressBar *mpb.Bar) {
			defer wg.Done() // Decrease counter when we exit here

			result := downloadReportWithProgress(report, fullDownloadPath, progressBar)
			results[index] = result
		}(i, report, progressBar)
	}

	p.Wait()
	return results
}
