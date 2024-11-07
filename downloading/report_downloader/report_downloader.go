package report_downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/F0903/pdf_downloader_uge5/downloading"
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
	*downloading.ConcurrentDownloader
	outputDir string
}

func NewReportDownloader(ctx context.Context, outputDir string) *ReportDownloader {
	dl := downloading.NewConcurrentDownloaderWithCtx(ctx)
	dl.SetResponseChecker(assertResponse)
	return &ReportDownloader{
		ConcurrentDownloader: dl,
		outputDir:            outputDir,
	}
}

func (dl *ReportDownloader) handleDownload(download *downloading.ConcurrentDownloadResult, reports []*models.Report, resultBuffer []*ReportDownloadResult) {
	index := download.OrderIndex
	err := download.Err

	associatedReport := reports[index]

	if err != nil {
		var downloadResult *ReportDownloadResult
		if err == context.Canceled {
			downloadResult = NewDownloadResult(associatedReport, NewCancelledDownloadState())
		} else {
			downloadResult = NewDownloadResult(
				associatedReport,
				NewFailedDownloadState(err),
			)
		}

		resultBuffer[index] = downloadResult
		return
	}

	fullPath := path.Join(dl.outputDir, associatedReport.Id+".pdf")
	file, err := os.Create(fullPath)
	if err != nil {
		resultBuffer[index] = NewDownloadResult(
			associatedReport,
			NewFailedDownloadState(err),
		)
	}
	defer file.Close()

	reader := download.Reader
	defer reader.Close()

	_, err = utils.CancellableCopy(dl.Ctx, file, reader)
	if err != nil {
		resultBuffer[index] = NewDownloadResult(
			associatedReport,
			NewFailedDownloadState(err),
		)
	}

	resultBuffer[index] = NewDownloadResult(
		associatedReport,
		NewSuccededDownloadState(fullPath),
	)
}

func (dl *ReportDownloader) handleDownloads(downloadChannel chan *downloading.ConcurrentDownloadResult, reports []*models.Report, resultBuffer []*ReportDownloadResult) {
	for {
		select {
		case <-dl.Ctx.Done():
			return
		case download := <-downloadChannel:
			if download == nil {
				return
			}

			go dl.handleDownload(download, reports, resultBuffer)
		}
	}
}

func (dl *ReportDownloader) DownloadReports(reports []*models.Report) []*ReportDownloadResult {
	results := make([]*ReportDownloadResult, len(reports))

	p := mpb.New(
		mpb.WithAutoRefresh(),
	)

	sendChannel := make(chan *downloading.ConcurrentDownloadCommand)
	downloadChannel := make(chan *downloading.ConcurrentDownloadResult)

	dl.DownloadWithProgress(
		sendChannel,
		downloadChannel,
	)

	// Needs to be done on own thread to avoid blocking
	go dl.handleDownloads(downloadChannel, reports, results)

	for i, report := range reports {
		// Filter out reports with missing URLs
		if report.PrimaryDownloadLink == "" && report.FallbackDownloadLink == "" {
			results[i] = NewDownloadResult(report, NewMissingDownloadState())
		}

		bar := p.AddBar(0,
			mpb.PrependDecorators(
				decor.Name(report.Id, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
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

		downloadCommand := downloading.NewConcurrentDownloadCommand(
			[]string{
				report.PrimaryDownloadLink,
				report.FallbackDownloadLink,
			},
			bar,
			i,
		)
		sendChannel <- downloadCommand
	}

	close(sendChannel)

	p.Wait()
	return results
}
