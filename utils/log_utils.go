package utils

import (
	"fmt"
	"log"
	"os"
)

// Returned file must remain open for the duration of log redirection
func RedirectLoggingToFile(path string) *os.File {
	logFile, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("could not create log file: %v", err))
	}
	log.SetOutput(logFile)
	return logFile
}
