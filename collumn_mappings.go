package main

import (
	"github.com/F0903/pdf_downloader_uge5/utils"
)

const (
	// We need to start at 1 because the map will return 0 for unknown values.
	IdColumn = iota + 1
	NameColumn
	PrimaryDownloadColumn
	SecondaryDownloadColumn
)

// Map Excel column names to our column "enum" values.
var ColumnMappings = map[int]int{
	utils.ColumnNameToIndex("A"):  IdColumn,
	utils.ColumnNameToIndex("C"):  NameColumn,
	utils.ColumnNameToIndex("AL"): PrimaryDownloadColumn,
	utils.ColumnNameToIndex("AM"): SecondaryDownloadColumn,
}
