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

func (s *day11Solution) blink(times uint) (total uint64) {
	const startSize = 1024 * 1024

	ruleSet := newPlutostoneRuleset()

	srcBuf := s.numbers
	for _, n := range srcBuf {
		count := ruleSet.calculateResultSplitRecursive(uint(n), times)
		pf("%v by %v blinks splits to %v numbers\n", n, times, count)
		total += count
	}
	return
}

type ruleset struct {
	fallbackRule func(num uint) uint
	rules        []ruleProcFunc
}

type (
	ruleAppResult uint8
	ruleProcFunc  func(num uint) (res ruleAppResult, n1 uint, n2 int)
)

const (
	ruleres_na ruleAppResult = iota
	ruleres_n1               // return 1 number
	ruleres_n2               // result 2 numbers - splitted
)

func newPlutostoneRuleset() ruleset {
	return ruleset{
		rules: []ruleProcFunc{
			ruleFor0, ruleForEvenDigits,
		},
		fallbackRule: defaultRule,
	}
}

func ruleFor0(num uint) (res ruleAppResult, n1 uint, n2 int) {
	switch num {
	case 0:
		res = ruleres_n1
		n1 = 1
	default:
		res = ruleres_na
	}
	return
}

func ruleForEvenDigits(num uint) (res ruleAppResult, n1 uint, n2 int) {
	digits := int(math.Log10(float64(num))) + 1
	if digits%2 != 0 {
		return
	}
	res = ruleres_n2
	div := uint(pow(10, uint8(digits/2)))
	n1, n2 = uint(num/div), int(num%div)
	return
}

func defaultRule(num uint) uint {
	const mul = 2024
	return num * mul
}

func (r ruleset) calculateResultSplitRecursive(n, times uint) (result uint64) {
	if times == 0 {
		return 1
	}
	times--
	n1, n2 := r.applyRules(n)
	result += r.calculateResultSplitRecursive(n1, times)
	if n2 >= 0 {
		result += r.calculateResultSplitRecursive(uint(n2), times)
	}
	// pf("in %v on top of %v blinks returning %v  (%v & %v)\n",
	// 	n, times, result, n1, n2)
	return
}

func (r ruleset) applyRules(num uint) (n1 uint, n2 int) {
	var res ruleAppResult
	for i := 0; i < len(r.rules); i++ {
		res, n1, n2 = r.rules[i](num)
		switch res {
		case ruleres_na:
			continue
		case ruleres_n1:
			n2 = -1
		}
		break
	}
	if res == ruleres_na {
		n1 = r.fallbackRule(num)
		n2 = -1
	}
	return
}
