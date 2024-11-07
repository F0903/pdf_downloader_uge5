package downloading

import "github.com/vbauerster/mpb/v8"

type ConcurrentDownloadCommand struct {
	urls        []string // Multiple urls to use for backup
	progressBar *mpb.Bar
	orderIndex  int
}

func NewConcurrentDownloadCommand(urls []string, progressBar *mpb.Bar, orderIndex int) *ConcurrentDownloadCommand {
	return &ConcurrentDownloadCommand{
		urls,
		progressBar,
		orderIndex,
	}
}
