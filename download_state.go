package main

import "fmt"

type DownloadStateEnum int

const (
	Done = iota
	Error
	MissingURL
)

type DownloadState struct {
	stateEnum DownloadStateEnum
	err       error
}

func NewDownloadState(state DownloadStateEnum) *DownloadState {
	return &DownloadState{
		stateEnum: state,
		err:       nil,
	}
}

func (state *DownloadState) String() string {
	switch state.stateEnum {
	case Done:
		return "Done"
	case Error:
		return fmt.Sprintf("Error: %v", state.err)
	case MissingURL:
		return "Missing URL"
	}
	return "Unknown DownloadState"
}
