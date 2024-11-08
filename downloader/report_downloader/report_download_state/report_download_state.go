package report_download_state

import (
	"fmt"
	"strings"
)

type ReportDownloadStateEnum int

const (
	done = iota + 1
	failed
	cancelled
	missingURLs
)

// This keeps track of the download state of each report,
// and based on the state, the associated error, or the path it was written to
type ReportDownloadState struct {
	stateEnum   ReportDownloadStateEnum
	err         error
	WrittenPath string
}

func NewSuccededState(writtenPath string) *ReportDownloadState {
	return &ReportDownloadState{
		stateEnum:   done,
		WrittenPath: writtenPath,
	}
}

func NewFailedState(err error) *ReportDownloadState {
	return &ReportDownloadState{
		stateEnum: failed,
		err:       err,
	}
}

func NewCancelledState() *ReportDownloadState {
	return &ReportDownloadState{
		stateEnum: cancelled,
	}
}

func NewMissingState() *ReportDownloadState {
	return &ReportDownloadState{
		stateEnum: missingURLs,
	}
}

// Has the download succeded?
func (state *ReportDownloadState) IsDone() bool {
	return state.stateEnum == done
}

// Set the error and set stateEnum to failed
func (state *ReportDownloadState) SetError(err error) {
	state.stateEnum = failed
	state.err = err
}

func (state *ReportDownloadState) String() string {
	switch state.stateEnum {
	case done:
		return "Done"
	case cancelled:
		return "Cancelled"
	case failed:
		return fmt.Sprintf("Error: %v", state.err)
	case missingURLs:
		return "Missing URLs"
	}
	return "Unknown DownloadState"
}

// Return a string representation of the download state with new lines replaced by commas
func (state *ReportDownloadState) StringNoNewLines() string {
	return strings.ReplaceAll(state.String(), "\n", ", ")
}
