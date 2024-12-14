package days

import (
	"bufio"
	"fmt"
	"io"
)

type (
	ioReader = io.Reader
)

func ProcessReader(reader ioReader) (scanner bufio.Reader) {
	scanner = *bufio.NewReader(reader)
	return
}

func pf(f string, a ...any) {
	fmt.Printf(f, a...)
}

func plf(f string, a ...any) {
	fmt.Printf(f, a...)
	fmt.Println()
}
