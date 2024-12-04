package days

import (
	"fmt"
	"strings"
)

func init() {
	DaySolutions[4] = &day4Solution{}
}

type day4Solution struct {
	charTable  [][]byte
	columnSize int
}

type seekDir = int

const (
	back seekDir = iota - 1
	none
	forth
)

func (s *day4Solution) HasData() bool {
	return true
}

func (s *day4Solution) ReadData(reader ioReader) (err error) {
	const startSize = 1024
	s.charTable = make([][]byte, 0, startSize)
	sc := ProcessReader(reader)
	for err == nil {
		var line []byte
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}
		if s.columnSize == 0 {
			s.columnSize = len(line)
			idxs := make([]string, len(line))
			for i := range idxs {
				idxs[i] = Stringify(i)
			}
			// fmt.Printf(" |%v\n", strings.Join(idxs, ""))
			// fmt.Println(strings.Repeat("-", len(line)+2))
		} else if s.columnSize != len(line) {
			err = fmt.Errorf("columns count mismatch: expected %v but actually %v (%q)",
				s.columnSize, len(line), line)
			break
		}
		// fmt.Printf("%v|%s\n", len(s.charTable), line)
		fmt.Printf("%s\n", line)
		l := make([]byte, len(line))
		copy(l, line)
		s.charTable = append(s.charTable, l)
	}
	fmt.Println("")
	printTable(s.charTable, 0, 0, len(s.charTable), s.columnSize)
	fmt.Println("")
	for _, l := range s.charTable {
		fmt.Printf("%s\n", l)
	}

	if err == ErrEOF {
		err = nil
	}
	return
}

func (s *day4Solution) SolvePt1() (answer string, err error) {
	checkBytes := []byte("XMAS")
	matchesCount := 0
	for l, line := range s.charTable {
		for c, ch := range line {
			// search for word start
			if ch != checkBytes[0] {
				continue
			}
			// seek up to 8 directions
			for vdir := back; vdir <= forth; vdir++ {
				for hdir := back; hdir <= forth; hdir++ {
					if hdir == 0 && vdir == 0 {
						continue
					}
					if hasMatch(s.charTable, checkBytes, hdir, vdir, l, c) {
						matchesCount++
					}
				}
			}
		}
	}
	answer = Stringify(matchesCount)
	return
}

func hasMatch(table [][]byte, value []byte, hdir, vdir seekDir, l, c int) bool {
	cachedL, cachedC := l, c
	for i, ch := range value {
		check := table[l][c]
		if ch != check {
			return false
		}
		// no need to make further checks
		if i == len(value)-1 {
			break
		}
		l += vdir
		if !isValidIdx(l, table) {
			return false
		}
		c += hdir
		if !isValidIdx(c, table[l]) {
			return false
		}
	}
	fmt.Printf("found %q in ({%v,%v}-{%v,%v})\n", value, cachedL, cachedC, l, c)
	return true
}

func isValidIdx[T any](idx int, slice []T) bool {
	return idx >= 0 && idx < len(slice)
}

func printTable(table [][]byte, lSince, cSince, vcount, hcount int) {
	if lSince >= len(table) {
		return
	}
	lBefore := lSince + vcount
	if lBefore >= len(table) {
		lBefore = len(table)
	}
	if cSince >= len(table[lSince]) {
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
