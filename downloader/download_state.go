package downloader

import (
	"fmt"
	"strings"
)

type DownloadStateEnum int

const (
	done = iota + 1
	failed
	missingURLs
)

type DownloadState struct {
	stateEnum   DownloadStateEnum
	err         error
	writtenPath string
}

func NewFailedDownloadState(err error) *DownloadState {
	return &DownloadState{
		stateEnum: failed,
		err:       err,
	}
}

func NewSuccededDownloadState(writtenPath string) *DownloadState {
	return &DownloadState{
		stateEnum:   done,
		writtenPath: writtenPath,
	}
}

func NewMissingDownloadState() *DownloadState {
	return &DownloadState{
		stateEnum: missingURLs,
	}
}

func (state *DownloadState) IsDone() bool {
	return state.stateEnum == done
}

func (state *DownloadState) SetError(err error) {
	state.stateEnum = failed
	state.err = err
}

func (state *DownloadState) String() string {
	switch state.stateEnum {
	case done:
		return "Done"
	case failed:
		return fmt.Sprintf("Error: %v", state.err)
	case missingURLs:
		return "Missing URLs"
	}
	return "Unknown DownloadState"
}

func (state *DownloadState) StringNoNewLines() string {
	switch state.stateEnum {
	case done:
		return "Done"
	case failed:
		return strings.ReplaceAll(fmt.Sprintf("Error: %v", state.err), "\n", ", ")
	case missingURLs:
		return "Missing URLs"
	}
	return "Unknown DownloadState"
}
