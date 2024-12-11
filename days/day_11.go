package days

import "math"

func init() {
	DaySolutions[11] = &day11Solution{}
}

type day11Solution struct {
	numbers []int
}

func (s *day11Solution) HasData() bool {
	return true
}

func (s *day11Solution) ReadData(reader ioReader) (err error) {
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()
	var line []byte
	line, _, err = sc.ReadLine()
	if err != nil && err != ErrEOF {
		pf("err: %v\n", err)
		return
	}
	s.numbers, err = ParseNumbers(string(line))
	pf("parsed: %v\n", s.numbers)
	return
}

func (s *day11Solution) SolvePt1() (answer string, err error) {
	count := s.blink(25)
	answer = Stringify(count)
	return
}

func (s *day11Solution) SolvePt2() (answer string, err error) {
	return
}

func (s *day11Solution) blink(times int) (total uint64) {
	const startSize = 1024 * 1024

	var defaultResBuf [2]int

	ruleSet := ruleset{
		rules: []func(num int) (result []int, isApplicable bool){
			func(num int) (result []int, isApplicable bool) {
				isApplicable = num == 0
				if isApplicable {
					defaultResBuf[0] = 1
					result = defaultResBuf[:1]
				}
				return
			},
			func(num int) (result []int, isApplicable bool) {
				digits := int(math.Log10(float64(num))) + 1
				isApplicable = digits%2 == 0
				if isApplicable {
					pow10 := func(pow int) int {
						res := 1
						for range pow {
							res *= 10
						}
						return res
					}
					splitter := pow10(digits / 2)
					defaultResBuf[0] = num / splitter // left part
					defaultResBuf[1] = num % splitter // right part
					result = defaultResBuf[:]
				}
				return
			},
		},
		fallbackRule: func(num int) []int {
			defaultResBuf[0] = num * 2024
			return defaultResBuf[:1]
		},
	}

	srcBuf := s.numbers
	for _, n := range srcBuf {
		count := ruleSet.calculateResultSplitRecursive(n, times)
		pf("%v by %v blinks splits to %v numbers\n", n, times, count)
		total += count
	}
	return
}

type ruleset struct {
	fallbackRule func(num int) []int
	rules        []func(num int) (result []int, isApplicable bool)
}

func (r ruleset) calculateResultSplitRecursive(n, times int) (result uint64) {
	if times == 0 {
		return 1
	}
	times--
	changed := r.applyRules(n)
	switch len(changed) {
	case 1:
		newNum := changed[0]
		// pf("%v -> %v\n", n, newNum)
		result += r.calculateResultSplitRecursive(newNum, times)
	case 2:
		n1, n2 := changed[0], changed[1]
		// pf("%v -> %v | %v\n", n, n1, n2)
		result += r.calculateResultSplitRecursive(n1, times)
		result += r.calculateResultSplitRecursive(n2, times)
	default:
		pf("unexpected case\n")
	}
	// pf("in %v on top of %v blinks returning %v\n", n, times, result)
	return
}

func (r ruleset) applyRules(num int) (result []int) {
	for i := 0; i < len(r.rules) && result == nil; i++ {
		var isApplicable bool
		result, isApplicable = r.rules[i](num)
		if !isApplicable {
			result = nil
		} else {
			// pf("applied rule %v\n", i)
		}
	}
	if result == nil {
		// pf("applied fallback\n")
		result = r.fallbackRule(num)
	}
	return
}
