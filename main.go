package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

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

	// Get the "reports"
	reports, err := excel.ReadReports(dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to read Excel: \n%w", err)
	}

	//DEBUGGING: get subset of reports
	reports = reports[:10]

	results := downloader.DownloadReports(reports, outputDir)
	err = excel.WriteDownloadResults(results, outputDir)
	if err != nil {
		return fmt.Errorf("failed to write download result metadata!\n%w", err)
	}

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
