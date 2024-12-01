package days

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	ioReader = io.Reader
	ioCloser = io.Closer

	daySolutionFunc = func(ioReader) (answer string, err error)
)

var DaySolutions map[uint8]daySolutionFunc

func init() {
	DaySolutions = make(map[uint8]daySolutionFunc, 25)
}

func StrToReader(input string) ioReader {
	return strings.NewReader(input)
}

func ReadFile(filename string) (reader ioReader, err error) {
	file, err := os.Open(filename)
	return file, err
}

func WrapSolution(solution daySolutionFunc) daySolutionFunc {
	return func(reader ioReader) (string, error) {
		defer func() {
			if closer, ok := reader.(ioCloser); ok {
				log.Printf("close ioCloser %v", reader)
				closer.Close()
			}
		}()
		return solution(reader)
	}
}

func ProcessReader(reader ioReader) (scanner bufio.Reader) {
	scanner = *bufio.NewReaderSize(reader, 10)
	return
}

func UnpackNumbers2(value string) (n1, n2 int, err error) {
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
