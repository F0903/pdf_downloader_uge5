package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/F0903/pdf_downloader_uge5/downloader"
	"github.com/F0903/pdf_downloader_uge5/excel"
)

func assertArgs(args []string) error {
	argLen := len(args)

	// + 1 because one arg is usually always provided by Windows
	const maxArgs = 2 + 1

	if argLen < 2 {
		return errors.New("you must provide a path to the excel spreadsheet")
	} else if argLen < 3 {
		return errors.New("you must provide a path to downloads output folder")
	} else if argLen > maxArgs {
		return errors.New("too many arguments")
	}

	return nil
}

func countSuccesfulDownloads(results []*downloader.DownloadResult) int {
	counter := 0
	for _, result := range results {
		if !result.State.IsDone() {
			continue
		}
		counter += 1
	}
	return counter
}

func run() error {
	args := os.Args
	if err := assertArgs(args); err != nil {
		return fmt.Errorf("argument error: %w", err)
	}

	dataFilePath := args[1]
	outputDir := args[2]

	// Create the output directory if it doesn’t exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: \nw%w", err)
	}

	startTime := time.Now()

	reports, err := excel.ReadReports(dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to read Excel: \n%w", err)
	}
	//DEBUGGING: get subset of reports
	reports = reports[:100]

	// Cancel downloads on CTRL+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	reportDownloader := downloader.NewReportDownloader(ctx, outputDir)
	defer reportDownloader.Close()

	results := reportDownloader.DownloadReports(reports)
	err = excel.WriteDownloadResults(results, outputDir)
	if err != nil {
		return fmt.Errorf("failed to write download result metadata!\n%w", err)
	}

	endTime := time.Since(startTime)
	fmt.Printf("Downloaded %d documents.\n", countSuccesfulDownloads(results))
	fmt.Printf("Time taken: %s\n", endTime)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error:\n %v\n", err)
	} else {
		fmt.Println("Done!")
	}

	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // Wait for Enter
}
