package main

import (
	"fmt"
	"os"

	"github.com/F0903/pdf_downloader_uge5/downloader"
	"github.com/F0903/pdf_downloader_uge5/excel"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("You must provide a path to the excel spreadsheet!")
		return
	}
	if len(args) < 3 {
		fmt.Println("You must provide a path to downloads output folder!")
		return
	}

	dataFilePath := args[1]
	outputPath := args[2]

	// Create the output directory if it doesn’t exist
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		fmt.Printf("Failed to create output directory: \n%v", err)
		return
	}

	// Get the "reports"
	reports, err := excel.ReadReports(dataFilePath)
	if err != nil {
		fmt.Printf("Failed to read Excel: \n%v", err)
		return
	}

	downloader.DownloadReports(reports, outputPath)
}
