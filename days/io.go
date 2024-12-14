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

func fplf(w io.Writer, f string, a ...any) {
	fmt.Fprintf(w, f, a...)
	fmt.Fprintln(w)
}
