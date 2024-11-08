package utils

import (
	"bufio"
	"os"
)

func WaitForKey(keyValue byte) {
	bufio.NewReader(os.Stdin).ReadBytes(keyValue) // Wait for Enter
}
