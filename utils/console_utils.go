package utils

import (
	"bufio"
	"os"
)

func WaitForInput(inputToWaitFor byte) {
	bufio.NewReader(os.Stdin).ReadBytes(inputToWaitFor)
}
