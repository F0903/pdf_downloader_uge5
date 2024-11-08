package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/F0903/pdf_downloader_uge5/args"
	"github.com/F0903/pdf_downloader_uge5/downloader/report_downloader"
	"github.com/F0903/pdf_downloader_uge5/excel"
	"github.com/F0903/pdf_downloader_uge5/utils"
)

func run() error {
	argMap, err := args.ParseArgs()
	if err != nil {
		return fmt.Errorf("argument error: %w", err)
	}

	// Check our required args are present
	if err := args.AssertArgsPresent(argMap, []string{"input_data", "output_dir"}); err != nil {
		return fmt.Errorf("missing args: %w", err)
	}

	excelDataPath := argMap["input_data"].Value
	outputDir := argMap["output_dir"].Value

	// Create the output directory if it doesn’t exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: \nw%w", err)
	}

	startTime := time.Now()

	reports, err := excel.ReadReports(excelDataPath)
	if err != nil {
		return fmt.Errorf("failed to read Excel: \n%w", err)
	}

	// Cancel downloads on CTRL+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	reportDownloader := report_downloader.NewReportDownloader(ctx, outputDir)
	defer reportDownloader.Close()

	results := reportDownloader.DownloadReports(reports)

	// Write our metadata
	err = excel.WriteDownloadResults(results, outputDir)
	if err != nil {
		return fmt.Errorf("failed to write download result metadata!\n%w", err)
	}

	endTime := time.Since(startTime)

	fmt.Printf("Downloaded %d documents.\n", report_downloader.CountSuccesfulReportDownloads(results))
	fmt.Printf("Time taken: %s\n", endTime.Round(time.Second))

	return nil
}

func main() {
	// We wrap all the stuff in the run() func so it's easier to
	if err := run(); err != nil {
		fmt.Printf("Error:\n %v\n", err)
	} else {
		fmt.Println("Done!")
	}

	fmt.Println("Press Enter to exit...")
	utils.WaitForKey('\n')
}
