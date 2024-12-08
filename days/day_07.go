package days

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func init() {
	DaySolutions[7] = &day7Solution{}
}

type day7Solution struct {
	calibrations []calibrationEntry
}

type calibrationEntry struct {
	numbers   []int
	testValue uint64
}

func (s *day7Solution) HasData() bool {
	return true
}

func (s *day7Solution) ReadData(reader ioReader) (err error) {
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()
	var line []byte
	line, _, err = sc.ReadLine()
	for err == nil {
		colonIdx := slices.Index(line, ':')
		if colonIdx < 1 {
			err = fmt.Errorf("expected colon sign after test value at ")
		}
		entry := calibrationEntry{}
		entry.testValue, err = (strconv.ParseUint(string(line[:colonIdx]), 10, 64))
		if err != nil {
			break
		}

		entry.numbers, err = ParseNumbers(string(line[colonIdx+1:]))
		if err != nil {
			break
		}
		s.calibrations = append(s.calibrations, entry)
		line, _, err = sc.ReadLine()
	}
	return
}

func (s *day7Solution) SolvePt1() (answer string, err error) {
	result := 0
	for _, c := range s.calibrations {
		// fmt.Printf("check equation: %v = %v\n", c.testValue, c.numbers)
		if foundMatch, comb := couldBeTrue(c.testValue, c.numbers); foundMatch {
			result += int(c.testValue)
			fmt.Printf("found combination: %v\n", stringifyOperators(c.numbers, comb))
		}
	}
	answer = Stringify(result)
	return
}

func (s *day7Solution) SolvePt2() (answer string, err error) {
	result := 0
	for _, c := range s.calibrations {
		// fmt.Printf("check equation: %v = %v\n", c.testValue, c.numbers)
		if foundMatch, comb := couldBeTrue_3ops(c.testValue, c.numbers); foundMatch {
			result += int(c.testValue)
			// fmt.Printf("found combination: %q\n", comb)
			if slices.Index(comb, '|') >= 0 {
				fmt.Printf("+ %v ", c.testValue)
			}
		}
	}
	answer = Stringify(result)
	return
}

func couldBeTrue(ctlNumber uint64, nums []int) (result bool, combCode int) {
	const opMask = 0b1
	totalCombinations := 1
	for range nums {
		totalCombinations *= 2
	}
	totalCombinations /= 2
	for i := 0; i <= totalCombinations; i++ {
		// use bits as sequence of +/* cells
		combCode = i
		checkResult := uint64(nums[0])
		for j, n := range nums[1:] {
			opCode := (combCode >> j) & opMask
			switch opCode {
			case 0:
				// sum
				checkResult += uint64(n)
			case 1:
				// multiply
				checkResult *= uint64(n)
			}
			if checkResult == ctlNumber {
				result = true
				return
			}
			if checkResult > ctlNumber {
				// no sense to check that combination further
				break
			}
		}
	}
	return
}

func couldBeTrue_3ops(testNum uint64, nums []int) (result bool, combination []byte) {
	const opMask = 0b1
	totalCombinations := 1
	for range nums {
		totalCombinations *= 3
	}
	totalCombinations /= 3
	combination = make([]byte, len(nums)-1)
	for i := range combination {
		combination[i] = '+'
	}
	for i := 0; i <= totalCombinations; i++ {
		// use bits as sequence of +/* cells
		checkResult := uint64(nums[0])
		for j, n := range nums[1:] {
			opCode := combination[j]
			switch opCode {
			case 0:
				fallthrough
			case '+':
				// sum
				checkResult += uint64(n)
			case '*':
				// multiply
				checkResult *= uint64(n)
			case '|':
				conc := fmt.Sprintf("%v%v", checkResult, n)
				newRes, err := strconv.ParseUint(conc, 10, 64)
				if err != nil {
					panic("impossible error: " + err.Error())
				}
				checkResult = newRes
			}
			if checkResult == testNum && j == int(len(combination)-1) {
				result = true
				return
			}
			if checkResult > testNum {
				// no sense to check that combination further
				break
			}
		}
		increment(combination)
	}
	return
}

func increment(ops []byte) {
	for i, op := range ops {
		shouldIncNext := false
		switch op {
		case 0:
			fallthrough
		case '+':
			ops[i] = '*'
		case '*':
			ops[i] = '|'
		case '|':
			ops[i] = '+'
			// only when overflow
			shouldIncNext = true
		}
		if !shouldIncNext {
			break
		}
	}
}

func stringifyOperators(nums []int, combination int) string {
	resultOps := make([]rune, len(nums)-1)
	for i := 0; i < int(len(resultOps)); i++ {
		opCode := (combination >> i) & 0b1
		switch opCode {
		case 0:
			resultOps[i] = '+'
		case 1:
			resultOps[i] = '*'
		}
	}
	strWithNums := make([]string, len(nums))
	for i, n := range nums {
		op := ' '
		if i < int(len(resultOps)) {
			op = resultOps[i]
		}
		strWithNums[i] = fmt.Sprintf("%v%v", n, string(op))
	}
	return strings.Join(strWithNums, "")
}
