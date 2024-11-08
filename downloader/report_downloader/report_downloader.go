package report_downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/F0903/pdf_downloader_uge5/downloader"
	"github.com/F0903/pdf_downloader_uge5/downloader/report_downloader/report_download_state"
	"github.com/F0903/pdf_downloader_uge5/models"
	"github.com/F0903/pdf_downloader_uge5/utils"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type ReportDownloader struct {
	*downloader.Downloader
	outputDir string
}

func isPdf(contentType string) bool {
	return contentType == "application/pdf"
}

// The response asserter for the report downloader
func ReportDownloaderResponseAsserter(resp *http.Response) error {
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

func NewReportDownloader(ctx context.Context, outputDir string) *ReportDownloader {
	dl := downloader.NewDownloader(ctx)
	dl.SetResponseAsserter(ReportDownloaderResponseAsserter)
	return &ReportDownloader{
		Downloader: dl,
		outputDir:  outputDir,
	}
}

func (dl *ReportDownloader) writeResponseToFileWithProgress(data *downloader.DownloadData, fullPath string, progressBar *mpb.Bar) error {
	// Create the download file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	contentLength := data.ContentLength
	reader := data.Reader
	defer reader.Close()

	// Set the "finsihed value" for our progress bar to the content length of the response
	progressBar.SetTotal(contentLength, false)

	// Proxy reader automatically increments our progress bar
	proxyReader := progressBar.ProxyReader(reader)
	defer proxyReader.Close()

	// Read from response and write to file whilst updating the progress bar
	if _, err := utils.CancellableCopy(dl.Ctx, file, proxyReader); err != nil {
		if err == context.Canceled {
			return err
		}
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func (dl *ReportDownloader) downloadResourceWithProgress(report *models.Report, fullDownloadPath string, progressBar *mpb.Bar) error {
	data, err := dl.Download(report)
	if err != nil {
		progressBar.Abort(true)
		return fmt.Errorf("download error: %w", err)
	}

	if err := dl.writeResponseToFileWithProgress(data, fullDownloadPath, progressBar); err != nil {
		return fmt.Errorf("could not write response to file: %w", err)
	}

	return nil
}

func (dl *ReportDownloader) downloadReportWithProgress(report *models.Report, fullDownloadPath string, progressBar *mpb.Bar) *ReportDownloadResult {
	// Exit early if we are missing both URLs
	if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
		progressBar.Abort(true)
		return NewReportDownloadResult(report, report_download_state.NewMissingState())
	}

	err := dl.downloadResourceWithProgress(report, fullDownloadPath, progressBar)
	if err != nil {
		if err == context.Canceled {
			return NewReportDownloadResult(report, report_download_state.NewCancelledState())
		}

		// If the error was not that the download has been cancelled, just return a generic error state
		return NewReportDownloadResult(report, report_download_state.NewFailedState(err))
	}

	progressBar.SetTotal(progressBar.Current(), true)
	return NewReportDownloadResult(report, report_download_state.NewSuccededState(fullDownloadPath))
}

// Download all reports concurrently
func (dl *ReportDownloader) DownloadReports(reports []*models.Report) []*ReportDownloadResult {
	results := make([]*ReportDownloadResult, len(reports))

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

		// Start new thread for each download
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
