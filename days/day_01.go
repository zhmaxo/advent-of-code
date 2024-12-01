package days

import (
	"io"
	"sort"
)

func init() {
	DaySolutions[1] = WrapSolution(solve_day1)
}

func solve_day1(reader ioReader) (answer string, err error) {
	const bufSize = 1024
	totalDist := uint(0)
	lhs := make([]int, 0, bufSize)
	rhs := make([]int, 0, bufSize)

	size := 0
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
		lNum, rNum, err = UnpackNumbers2(string(line))
		if err != nil {
			return
		}
		lhs = append(lhs, lNum)
		rhs = append(rhs, rNum)
		size++
	}
	sort.Ints(lhs)
	sort.Ints(rhs)
	for i := 0; i < size; i++ {
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
