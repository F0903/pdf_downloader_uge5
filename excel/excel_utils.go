package excel

import "github.com/F0903/pdf_downloader_uge5/utils"

func getCharValue(char rune) int {
	// The range that column names can be goes from A to Z
	// The numeric value for A is 65
	// So we can get the value starting from 0 by subtracting this
	return int(char - 65)
}

// This took way too long to get right
// and I just found out the excel library
// has this function
func ColumnNameToIndex(name string) int {
	const base = 25
	total := 0
	for i, char := range name {
		charValue := getCharValue(char) + 1

		exp := len(name) - 1 - i
		pow := utils.IntPow(base, exp)

		totalCharValue := charValue*pow + (pow / base)
		total += totalCharValue
	}
	return total - 1
}
