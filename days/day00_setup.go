package days

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	ioReader = io.Reader
)

type Solution interface {
	HasData() bool
	ReadData(reader ioReader) (err error)
	SolvePt1() (string, error)
	SolvePt2() (string, error)
}

var (
	DaySolutions map[uint8]Solution

	ErrEOF    = io.EOF
	ErrNoData = fmt.Errorf("no data prepared!")
)

func init() {
	DaySolutions = make(map[uint8]Solution, 25)
}

func ProcessReader(reader ioReader) (scanner bufio.Reader) {
	scanner = *bufio.NewReaderSize(reader, 64)
	return
}

func ParseNumbers(value string) (result []int, err error) {
	unpacked := strings.Fields(value)
	result = make([]int, 0, len(unpacked))
	num := 0
	for _, v := range unpacked {
		num, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		result = append(result, num)
	}
	return
}

func ParseNumbers2(value string) (n1, n2 int, err error) {
	const mustLen = 2

	unpacked := strings.Fields(value)
	if len(unpacked) != mustLen {
		err = fmt.Errorf("expected len %v but actual is %v",
			mustLen, len(unpacked))
		return
	}
	s1, s2 := unpacked[0], unpacked[1]
	n1, err = strconv.Atoi(s1)
	if err != nil {
		return
	}
	n2, err = strconv.Atoi(s2)
	return
}

func Stringify(v any) string {
	return fmt.Sprint(v)
}
