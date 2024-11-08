package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrorEmptyURL = errors.New("empty url")

type ResponseAsserter = func(*http.Response) error

type Downloader struct {
	httpClient       *http.Client
	Ctx              context.Context
	responseAsserter ResponseAsserter
}

type DownloadData struct {
	Reader        io.ReadCloser
	ContentLength int64
}

// The default respose asserter, that just checks for status 200
func DefaultDownloaderResponseAsserter(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code was not OK: %d", resp.StatusCode)
	}
	return nil
}

func NewDownloader(ctx context.Context) *Downloader {
	httpClient := &http.Client{}
	return &Downloader{
		httpClient,
		ctx,
		DefaultDownloaderResponseAsserter,
	}
}

// Sets the function responsible for asserting HTTP response is valid.
func (dl *Downloader) SetResponseAsserter(asserter ResponseAsserter) {
	dl.responseAsserter = asserter
}

func (dl *Downloader) Close() {
	dl.httpClient.CloseIdleConnections()
}

func (dl *Downloader) downloadUrl(url string) (*DownloadData, error) {
	req, err := http.NewRequestWithContext(dl.Ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP GET request %w", err)
	}

	resp, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := dl.responseAsserter(resp); err != nil {
		return nil, err
	}

	return &DownloadData{resp.Body, resp.ContentLength}, nil
}

// Caller must handle closing of reader
// Both return valus will be nil on empty download strings
// Returns a data stream and the length of the stream
func (dl *Downloader) Download(downloadable Downloadable) (*DownloadData, error) {
	urls := downloadable.GetDownloadableURLs()

	// We use an empty error like this to join the later ones onto
	var combinedErr error
	for _, url := range urls {
		if url == "" {
			combinedErr = errors.Join(combinedErr, ErrorEmptyURL)
			continue
		}

		data, err := dl.downloadUrl(url)
		if err != nil {
			combinedErr = errors.Join(combinedErr, err)
			continue
		}
		return data, nil
	}

	return nil, combinedErr
}
