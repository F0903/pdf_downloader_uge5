package downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/F0903/pdf_downloader_uge5/models"
	"github.com/F0903/pdf_downloader_uge5/utils"
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

type ReportDownloader struct {
	httpClient *http.Client
	ctx        context.Context
	outputDir  string
}

func NewReportDownloader(ctx context.Context, outputDir string) *ReportDownloader {
	httpClient := &http.Client{}
	return &ReportDownloader{
		httpClient,
		ctx,
		outputDir,
	}
}

func (dl *ReportDownloader) Close() {
	dl.httpClient.CloseIdleConnections()
}

func (dl *ReportDownloader) downloadResourceWithProgress(url string, fullDownloadPath string, progressBar *mpb.Bar) error {
	req, err := http.NewRequestWithContext(dl.ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("could not create HTTP GET request %w", err)
	}

	resp, err := dl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed %w", err)
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
	if _, err := utils.CancellableCopy(dl.ctx, file, proxyReader); err != nil {
		if err == context.Canceled {
			return err
		}
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func (dl *ReportDownloader) downloadReportWithProgress(report *models.Report, fullDownloadPath string, progressBar *mpb.Bar) *DownloadResult {
	if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
		progressBar.Abort(true)
		return NewDownloadResult(report, NewMissingDownloadState())
	}

	currentUrl := report.PrimaryDownloadLink
	onFallback := false
	var fullErr error // Start with an empty error to wrap around
	for {
		err := dl.downloadResourceWithProgress(currentUrl, fullDownloadPath, progressBar)
		if err == nil {
			break
		}

		if err == context.Canceled {
			progressBar.Abort(true)
			return NewDownloadResult(report, NewCancelledDownloadState())
		}

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
	}

	progressBar.SetTotal(progressBar.Current(), true)
	return NewDownloadResult(report, NewSuccededDownloadState(fullDownloadPath))
}

func (dl *ReportDownloader) DownloadReports(reports []*models.Report) []*DownloadResult {
	results := make([]*DownloadResult, len(reports))

	var wg sync.WaitGroup
	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithAutoRefresh(),
	)

	for i, report := range reports {
		wg.Add(1)

		fileName := report.Id
		fullDownloadPath := path.Join(dl.outputDir, fileName+".pdf")

		// It's important to create the progress bar here and not in the new thread or it will panic
		progressBar := p.AddBar(0,
			mpb.PrependDecorators(
				decor.Name(fileName, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			),
			mpb.AppendDecorators(
				decor.OnAbort(decor.AverageETA(decor.ET_STYLE_GO, decor.WC{C: decor.DindentRight | decor.DextraSpace}), ""),
				decor.OnAbort(
					decor.Percentage(),
					"stopping...",
				),
			),
			mpb.BarRemoveOnComplete(),
		)

		go func() {
			defer wg.Done()
			result := dl.downloadReportWithProgress(report, fullDownloadPath, progressBar)
			ValidateDownloadResult(result)
			// Since each thread has a unique index this is thread safe, and also preserves the order.
			results[i] = result
		}()
	}

	p.Wait()
	return results
}
