package utils

func IntPow(base, exp int) int {
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}

func getCharValue(char rune) int {
	// The range that column names can be goes from A to Z
	// The numeric value for A is 65
	// So we can get the value starting from 0 by subtracting this
	return int(char - 65)
}

// This took way too long to get right
func ColumnNameToIndex(name string) int {
	const base = 25
	total := 0
	for i, char := range name {
		charValue := getCharValue(char) + 1

		exp := len(name) - 1 - i
		pow := IntPow(base, exp)

		totalCharValue := charValue*pow + (pow / base)
		total += totalCharValue
	}
	return total - 1
}
