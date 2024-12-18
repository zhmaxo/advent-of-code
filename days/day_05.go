package days

import (
	"fmt"
	"sort"
)

func init() {
	DaySolutions[5] = &day5Solution{}
}

type day5Solution struct {
	rules   []pageOrderingRule
	updates []pageUpdate
}

func (s *day5Solution) HasData() bool {
	return len(s.rules) > 0 && len(s.updates) > 0
}

func (s *day5Solution) ReadData(reader ioReader) (err error) {
	const startSize = 1024
	s.rules = make([]pageOrderingRule, 0, startSize)
	s.updates = make([]pageUpdate, 0, startSize)

	expectRule := true
	sc := ProcessReader(reader)
	for err == nil {
		var line []byte
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}
		if expectRule {
			if len(line) == 0 {
				expectRule = false
				continue
			}

			// parse rule
			rule := pageOrderingRule{}
			var numbers []int
			numbers, err = ParseNumbersSep(string(line), '|')
			if err == nil && len(numbers) != 2 {
				err = fmt.Errorf("not supported rule page nums count %v", len(numbers))
			}
			if err != nil {
				break
			}
			rule.pFirst, rule.pSecond = numbers[0], numbers[1]
			s.rules = append(s.rules, rule)
		} else {
			// parse page update
			update := pageUpdate{}
			update.pages, err = ParseNumbersSep(string(line), ',')
			if err != nil {
				break
			}
			s.updates = append(s.updates, update)
		}
	}
	if err == ErrEOF {
		err = nil
	}
	return
}

func (s *day5Solution) SolvePt1() (answer string, err error) {
	result := 0
	for _, u := range s.updates {
		isMatch := true
		for _, r := range s.rules {
			if !r.matches(u) {
				isMatch = false
				break
			}
		}
		if isMatch {
			result += u.pages[len(u.pages)/2]
		}
	}
	answer = Stringify(result)
	return
}

func (s *day5Solution) SolvePt2() (answer string, err error) {
	result := 0
	for _, u := range s.updates {
		isMatch := true
		for _, r := range s.rules {
			if !r.matches(u) {
				isMatch = false
				break
			}
		}
		if !isMatch {
			// sort pages
			sort.SliceStable(u.pages, func(i, j int) bool {
				for _, r := range s.rules {
					if r.compare(u.pages[i], u.pages[j]) > 0 {
						return false
					}
				}
				return true
			})
			result += u.pages[len(u.pages)/2]
		}
	}
	answer = Stringify(result)
	return
}

type pageOrderingRule struct {
	pFirst, pSecond int
}

func (r pageOrderingRule) matches(update pageUpdate) bool {
	pSecondOccured := false
	for _, v := range update.pages {
		switch v {
		case r.pFirst:
			if pSecondOccured {
				return false
			}
		case r.pSecond:
			pSecondOccured = true
		}
	}
	return true
}

func (r pageOrderingRule) compare(p1, p2 int) int {
	if p1 == r.pFirst && p2 == r.pSecond {
		return -1
	}
	if p1 == r.pSecond && p2 == r.pFirst {
		// that ordering breaks the rule
		return 1
	}
	return 0
}

type pageUpdate struct {
	pages []int
}
