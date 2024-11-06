package downloader

import (
	"fmt"
	"strings"
)

type DownloadStateEnum int

const (
	Done = iota + 1
	Error
	MissingURLs
)

type DownloadState struct {
	stateEnum DownloadStateEnum
	err       error
}

func NewDownloadState(state DownloadStateEnum, associatedError error) *DownloadState {
	return &DownloadState{
		stateEnum: state,
		err:       associatedError,
	}
}

func (state *DownloadState) String() string {
	switch state.stateEnum {
	case Done:
		return "Done"
	case Error:
		return fmt.Sprintf("Error: %v", state.err)
	case MissingURLs:
		return "Missing URLs"
	}
	return "Unknown DownloadState"
}

func (state *DownloadState) StringNoNewLines() string {
	switch state.stateEnum {
	case Done:
		return "Done"
	case Error:
		return strings.ReplaceAll(fmt.Sprintf("Error: %v", state.err), "\n", ", ")
	case MissingURLs:
		return "Missing URLs"
	}
	return "Unknown DownloadState"
}
