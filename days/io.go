package days

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type (
	ioReader = io.Reader
	ioWriter = io.Writer
)

var defaultWriter = os.Stdout

func ProcessReader(reader ioReader) (scanner bufio.Reader) {
	scanner = *bufio.NewReader(reader)
	return
}

func osCreate(path string) (*os.File, error) {
	return os.Create(path)
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

func printRectFunc(rect rect, charFunc func(posInt) rune) {
	fprintRectFunc(defaultWriter, rect, charFunc)
}

func fprintRectFunc(w ioWriter, rect rect, charFunc func(posInt) rune) {
	p := rect.topLeft
	for ; p.y < rect.size.y; p.y++ {
		for ; p.x < rect.size.x; p.x++ {
			fmt.Fprintf(w, "%s", string(charFunc(p)))
		}
		fmt.Fprintln(w)
		p.x = rect.topLeft.x
	}
}

func printTable(table [][]byte, lSince, cSince, vcount, hcount int) {
	if lSince < 0 || lSince >= len(table) {
		return
	}
	lBefore := lSince + vcount
	if lBefore >= len(table) {
		lBefore = len(table)
	}
	if cSince < 0 || cSince >= len(table[lSince]) {
		return
	}
	cBefore := cSince + hcount
	if cBefore >= len(table[lSince]) {
		cBefore = len(table[lSince])
	}
	cIdxs := make([]string, cBefore-cSince)
	for i := range cIdxs {
		cIdxs[i] = Stringify(cSince + i)
	}
	fmt.Printf(" |%v\n", strings.Join(cIdxs, ""))
	fmt.Printf("-|%v\n", strings.Repeat("-", cBefore-cSince))
	for l := lSince; l < lBefore; l++ {
		fmt.Printf("%v|%s\n", l, table[l][cSince:cBefore])
	}
}
