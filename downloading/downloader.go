package downloading

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"

	"github.com/vbauerster/mpb/v8"
)

type ConcurrentDownloadResult struct {
	OrderIndex int
	Reader     io.ReadCloser
	Err        error
}

type ConcurrentDownloader struct {
	httpClient      *http.Client
	Ctx             context.Context
	responseChecker func(resp *http.Response) error
	activeWorkers   atomic.Int32
}

func NewConcurrentDownloaderWithCtx(ctx context.Context) *ConcurrentDownloader {
	return &ConcurrentDownloader{
		httpClient:      http.DefaultClient,
		Ctx:             ctx,
		responseChecker: DefaultResponseChecker,
	}
}

func DefaultResponseChecker(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code was not OK: %d", resp.StatusCode)
	}
	return nil
}

func (dl *ConcurrentDownloader) Close() {
	dl.httpClient.CloseIdleConnections()
}

func (dl *ConcurrentDownloader) SetResponseChecker(responseChecker func(resp *http.Response) error) {
	dl.responseChecker = responseChecker
}

func (dl *ConcurrentDownloader) getResourceWithProgress(urls []string, progressBar *mpb.Bar) (io.ReadCloser, error) {
	combinedErrors := errors.New("")

	for _, url := range urls {
		if url == "" {
			continue
		}

		req, err := http.NewRequestWithContext(dl.Ctx, "GET", url, nil)
		if err != nil {
			progressBar.Abort(true)
			combinedErrors = errors.Join(combinedErrors, fmt.Errorf("could not create HTTP GET request %w", err))
			continue
		}

		resp, err := dl.httpClient.Do(req)
		if err != nil {
			progressBar.Abort(true)
			combinedErrors = errors.Join(combinedErrors, fmt.Errorf("HTTP GET request failed %w", err))
			continue
		}

		if err := dl.responseChecker(resp); err != nil {
			progressBar.Abort(true)
			combinedErrors = errors.Join(combinedErrors, err)
			continue
		}

		progressBar.SetTotal(resp.ContentLength, false)
		proxyReader := progressBar.ProxyReader(resp.Body)
		return proxyReader, nil
	}

	return nil, errors.Join(combinedErrors, errors.New("all URLs have failed"))
}

func (dl *ConcurrentDownloader) waitForInput(sendChannel chan *ConcurrentDownloadCommand, outputChannel chan *ConcurrentDownloadResult) {
	for {
		select {
		case <-dl.Ctx.Done():
			return

		case input := <-sendChannel:
			if input == nil {
				return
			}

			dl.activeWorkers.Add(1)
			go func() {
				reader, err := dl.getResourceWithProgress(input.urls, input.progressBar)
				result := &ConcurrentDownloadResult{
					input.orderIndex, reader, err,
				}

				outputChannel <- result
				activeWorkers := dl.activeWorkers.Add(-1)
				if activeWorkers == 0 {
					close(outputChannel)
				}
			}()
		}
	}
}

func (dl *ConcurrentDownloader) DownloadWithProgress(sendChannel chan *ConcurrentDownloadCommand, outputChannel chan *ConcurrentDownloadResult) {
	// Needs to be done on own thread to avoid blocking
	go dl.waitForInput(sendChannel, outputChannel)
}
