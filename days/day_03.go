package days

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func init() {
	DaySolutions[3] = &day3Solution{}
}

type day3Solution struct {
	muls map[exprGroup][]parsedMul
}

type parsedMul struct {
	lhs, rhs int
}

type (
	exprGroup   uint8
	groupedMuls = map[exprGroup][]parsedMul
)

const (
	EG_Enabled exprGroup = iota
	EG_Disabled
)

func (s *day3Solution) HasData() bool {
	return true
}

func (s *day3Solution) ReadData(reader ioReader) (err error) {
	s.muls, err = d3_read(reader)
	return
}

func pf(f string, a ...any) {
	fmt.Printf(f, a...)
}

func d3_read(reader ioReader) (muls groupedMuls, err error) {
	const (
		startCacheSize = 512
		openToken      = "mul("
		closeToken     = ')'
		enableInstr    = "do()"
		disableInstr   = "don't()"
	)

	openBytes := []byte(openToken)
	enableBytes := []byte(enableInstr)
	disableBytes := []byte(disableInstr)

	muls = make(map[exprGroup][]parsedMul, 2)
	muls[EG_Enabled] = make([]parsedMul, 0, startCacheSize)
	muls[EG_Disabled] = make([]parsedMul, 0, startCacheSize)
	sc := ProcessReader(reader)

	// prevSkipped := false
	// lastSkipReason := ""
	currentGroup := EG_Enabled
	for err == nil {
		// if prevSkipped {
		// 	pf("skip")
		// 	switch lastSkipReason {
		// 	case "":
		// 		pf("\n")
		// 	default:
		// 		pf(" (%v)\n", lastSkipReason)
		// 	}
		// }
		// prevSkipped = true
		// lastSkipReason = ""
		var line []byte
		var groupChanged bool
		line, err = sc.ReadSlice(')')
		// pf("check line %q: ", line)
		// change currentGroup only if opposite instruction found
		switch currentGroup {
		case EG_Enabled:
			if endsWith(line, disableBytes) {
				currentGroup = EG_Disabled
				groupChanged = true
			}
		case EG_Disabled:
			if endsWith(line, enableBytes) {
				currentGroup = EG_Enabled
				groupChanged = true
			}
		}
		if groupChanged {
			// prevSkipped = false
			// pf("[group changed to %v]\n", currentGroup)
			continue
		}
		openSubstrIdx := strings.LastIndex(string(line), string(openBytes))
		if openSubstrIdx < 0 || openSubstrIdx > len(line)-2 {
			// lastSkipReason = fmt.Sprintf("len not match (osi=%v)", openSubstrIdx)
			continue
		}
		args := line[openSubstrIdx+len(openBytes) : len(line)-1]
		sepIdx := slices.Index(args, ',')
		if sepIdx < 1 {
			// lastSkipReason = fmt.Sprintf("sep not found in %q", args)
			continue
		}
		instr := parsedMul{}
		var parseErr error // should be ignored by loop
		instr.lhs, parseErr = strconv.Atoi(string(args[:sepIdx]))
		if parseErr != nil {
			// lastSkipReason = fmt.Sprintf("can't parse %q", args[:sepIdx])
			continue
		}
		instr.rhs, parseErr = strconv.Atoi(string(args[sepIdx+1:]))
		if parseErr != nil {
			// lastSkipReason = fmt.Sprintf("can't parse %q", args[sepIdx+1:])
			continue
		}
		// prevSkipped = false
		// pf("parsed %v\n", instr)
		muls[currentGroup] = append(muls[currentGroup], instr)
	}

	if err == ErrEOF {
		err = nil
	}
	return
}

func endsWith(line, ending []byte) bool {
	return len(line) >= len(ending) &&
		bytes.Equal(line[len(line)-len(ending):], ending)
}

func substrIdx(line, substr []byte) (idx int) {
	// no substr len check because of trust to myself)
	idx = slices.Index(line, substr[0])
	if idx < 0 {
		// pf("not found %q in %q", substr[0], line)
		return
	}
	endIdx := idx + len(substr)
	if endIdx >= len(line) {
		// pf("not match line stats: len(%q)=%v, idx=%v, len(%q)=%v",
		// 	line, len(line), idx, substr, len(substr))
		idx = -1
		return
	}
	if !bytes.Equal(line[idx:endIdx], substr) {
		// pf("%q != %q", line[idx:endIdx], substr)
		idx = -1
	}
	return
}

func d3_readWithRegexp(reader ioReader) (muls map[exprGroup][]parsedMul, err error) {
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

	muls = make(map[exprGroup][]parsedMul, 2)
	muls[EG_Enabled] = make([]parsedMul, 0, startCacheSize)
	muls[EG_Disabled] = make([]parsedMul, 0, startCacheSize)
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
				err = nil
				continue
			}
			instr.rhs, err = strconv.Atoi(string(mulMatch[2]))
			if err != nil {
				// ignore invalid values
				err = nil
				continue
			}
			muls[currentGroup] = append(muls[currentGroup], instr)
		}
	}

	if err == ErrEOF {
		err = nil
	}

	return
}

func (s *day3Solution) SolvePt1() (answer string, err error) {
	result := 0
	for _, g := range s.muls {
		for _, v := range g {
			result += v.lhs * v.rhs
		}
	}
	answer = Stringify(result)
	return
}

func (s *day3Solution) SolvePt2() (answer string, err error) {
	result := 0
	enabledMuls := s.muls[EG_Enabled]
	for _, v := range enabledMuls {
		result += v.lhs * v.rhs
	}
	answer = Stringify(result)
	return
}
