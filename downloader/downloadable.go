package downloader

// The interface for downloadable types
type Downloadable interface {
	GetDownloadableURLs() []string
}
