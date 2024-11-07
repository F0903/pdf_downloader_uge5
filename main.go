package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/F0903/pdf_downloader_uge5/downloading/report_downloader"
	"github.com/F0903/pdf_downloader_uge5/excel"
	"github.com/F0903/pdf_downloader_uge5/utils"
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

	//DEBUGGING: only download subset of reports
	reports = reports[:2]

	if err != nil {
		return fmt.Errorf("failed to read Excel: \n%w", err)
	}

	// Cancel downloads on CTRL+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	reportDownloader := report_downloader.NewReportDownloader(ctx, outputDir)
	defer reportDownloader.Close()

	results := reportDownloader.DownloadReports(reports)
	err = excel.WriteDownloadResults(results, outputDir)
	if err != nil {
		return fmt.Errorf("failed to write download result metadata!\n%w", err)
	}

	endTime := time.Since(startTime)
	fmt.Printf("Downloaded %d documents.\n", report_downloader.CountSuccesfulDownloads(results))
	fmt.Printf("Time taken: %s\n", endTime)

	return nil
}

func main() {
	logFile := utils.RedirectLoggingToFile("log.txt")
	defer logFile.Close()

	if err := run(); err != nil {
		fmt.Printf("Error:\n %v\n", err)
	} else {
		fmt.Println("Done!")
	}

	fmt.Println("Press Enter to exit...")
	utils.WaitForInput('\n')
}
