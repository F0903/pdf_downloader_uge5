package excel

const (
	// We need to start at 1 because the map will return 0 for unknown values.
	IdColumn = iota + 1
	NameColumn
	PrimaryDownloadColumn
	SecondaryDownloadColumn
)

// Map Excel column names to our column "enum" values.
var ColumnMappings = map[int]int{
	ColumnNameToIndex("A"):  IdColumn,
	ColumnNameToIndex("C"):  NameColumn,
	ColumnNameToIndex("AL"): PrimaryDownloadColumn,
	ColumnNameToIndex("AM"): SecondaryDownloadColumn,
}
