package days

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func init() {
	DaySolutions[3] = &day3Solution{}
}

type day3Solution struct {
	matches [][][]byte
}

func (s *day3Solution) HasData() bool {
	return true
}

func (s *day3Solution) ReadData(reader ioReader) (err error) {
	const regexStr = `mul\(([0-9]*),([0-9]*)\)`
	r := regexp.MustCompile(regexStr)
	// scanner := ProcessReader(reader)
	d, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	s.matches = r.FindAllSubmatch(d, -1)
	fmt.Printf("%q\n", s.matches)
	return
}

func (s *day3Solution) SolvePt1() (answer string, err error) {
	result := 0
	for _, v := range s.matches {
		var num1, num2 int
		num1, err = strconv.Atoi(string(v[1]))
		if err != nil {
			return "", err
		}
		num2, err = strconv.Atoi(string(v[2]))
		if err != nil {
			return "", err
		}
		result += num1 * num2
	}
	answer = Stringify(result)
	return
}

func (s *day3Solution) SolvePt2() (answer string, err error) {
	return
}
