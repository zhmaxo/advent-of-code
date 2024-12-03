package days

import (
	"regexp"
	"strconv"
)

func init() {
	DaySolutions[3] = &day3Solution{}
}

type day3Solution struct {
	groupedMuls map[ExprGroup][]parsedMul
}

type parsedMul struct {
	lhs, rhs int
}

type ExprGroup uint8

const (
	EG_Enabled ExprGroup = iota
	EG_Disabled
)

func (s *day3Solution) HasData() bool {
	return true
}

func (s *day3Solution) ReadData(reader ioReader) (err error) {
	s.groupedMuls, err = d3_readWithRegexp(reader)
	return
}

func d3_readWithRegexp(reader ioReader) (groupedMuls map[ExprGroup][]parsedMul, err error) {
	const (
		startCacheSize = 512
		closeToken     = ')'
		mulPattern     = `mul\(([0-9]*),([0-9]*)\)`
		enablePattern  = `do\(\)`
		disablePattern = `don't\(\)`
	)

	mulRe := regexp.MustCompile(mulPattern)
	enableRe := regexp.MustCompile(enablePattern)
	disnableRe := regexp.MustCompile(disablePattern)

	groupedMuls = make(map[ExprGroup][]parsedMul, 2)
	groupedMuls[EG_Enabled] = make([]parsedMul, 0, startCacheSize)
	groupedMuls[EG_Disabled] = make([]parsedMul, 0, startCacheSize)
	sc := ProcessReader(reader)

	currentGroup := EG_Enabled
	for err == nil {
		var line []byte
		line, err = sc.ReadSlice(')')
		// change currentGroup only if opposite instruction found
		switch currentGroup {
		case EG_Enabled:
			if disnableRe.FindIndex(line) != nil {
				currentGroup = EG_Disabled
			}
		case EG_Disabled:
			if enableRe.FindIndex(line) != nil {
				currentGroup = EG_Enabled
			}
		}
		mulMatch := mulRe.FindSubmatch(line)
		if mulMatch != nil {
			instr := parsedMul{}
			instr.lhs, err = strconv.Atoi(string(mulMatch[1]))
			if err != nil {
				// ignore invalid values
				continue
			}
			instr.rhs, err = strconv.Atoi(string(mulMatch[2]))
			if err != nil {
				// ignore invalid values
				continue
			}
			groupedMuls[currentGroup] = append(groupedMuls[currentGroup], instr)
		}
	}

	if err == ErrEOF {
		err = nil
	}

	return
}

func (s *day3Solution) SolvePt1() (answer string, err error) {
	result := 0
	for _, g := range s.groupedMuls {
		for _, v := range g {
			result += v.lhs * v.rhs
		}
	}
	answer = Stringify(result)
	return
}

func (s *day3Solution) SolvePt2() (answer string, err error) {
	result := 0
	enabledMuls := s.groupedMuls[EG_Enabled]
	for _, v := range enabledMuls {
		result += v.lhs * v.rhs
	}
	answer = Stringify(result)
	return
}
