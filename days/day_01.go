package days

import (
	"io"
	"sort"
)

func init() {
	DaySolutions[1] = &day1Solution{}
}

type day1Solution struct {
	leftList, rightList []int
}

func (s *day1Solution) HasData() bool {
	return s.leftList != nil && s.rightList != nil
}

func (s *day1Solution) ReadData(reader ioReader) (err error) {
	const bufSize = 1024
	s.leftList = make([]int, 0, bufSize)
	s.rightList = make([]int, 0, bufSize)

	scanner := ProcessReader(reader)
	hasLine := true
	for hasLine {
		var line []byte
		var lNum, rNum int

		line, _, err = scanner.ReadLine()
		if err != nil {
			if err == io.EOF {
				hasLine = false
				err = nil
				break
			} else {
				return
			}
		}
		lNum, rNum, err = ParseNumbers2(string(line))
		if err != nil {
			return
		}
		s.leftList = append(s.leftList, lNum)
		s.rightList = append(s.rightList, rNum)
	}
	return
}

func (s *day1Solution) SolvePt1() (answer string, err error) {
	totalDist := uint(0)
	if !s.HasData() {
		err = ErrNoData
		return
	}

	lhs, rhs := s.leftList, s.rightList

	sort.Ints(lhs)
	sort.Ints(rhs)
	for i := 0; i < len(lhs); i++ {
		totalDist += dist(lhs[i], rhs[i])
	}
	answer = Stringify(totalDist)
	return
}

func dist(lhs, rhs int) uint {
	biggest, smallest := lhs, rhs
	if biggest < smallest {
		biggest, smallest = smallest, biggest
	}
	return uint(biggest - smallest)
}
