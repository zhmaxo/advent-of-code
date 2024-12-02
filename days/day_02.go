package days

type reportTrend uint8

const (
	RepTrend_Asc reportTrend = iota
	RepTrend_Desc
	RepTrend_Eq
)

func init() {
	DaySolutions[2] = &day2Solution{}
}

type day2Solution struct {
	reports [][]int
}

func (s *day2Solution) HasData() bool {
	return s.reports != nil
}

func (s *day2Solution) ReadData(reader ioReader) (err error) {
	const bufSize = 1024
	s.reports = make([][]int, 0, bufSize)

	scanner := ProcessReader(reader)
	hasLine := true

	var line []byte
	var report []int
	for hasLine {

		line, _, err = scanner.ReadLine()
		if err != nil {
			if err == ErrEOF {
				hasLine = false
				err = nil
				break
			} else {
				return
			}
		}
		report, err = ParseNumbers(string(line))
		if err != nil {
			return
		}
		s.reports = append(s.reports, report)
	}
	return
}

func (s *day2Solution) SolvePt1() (answer string, err error) {
	safeReports := 0
	const bufSize = 10
	for _, report := range s.reports {
		if isReportSafe(report) {
			safeReports++
		}
	}
	answer = Stringify(safeReports)
	return
}

func (s *day2Solution) SolvePt2() (answer string, err error) {
	safeReports := 0
	const bufSize = 10
	for _, report := range s.reports {
		safe := isReportSafe(report)
		for i := 0; i < len(report) && !safe; i++ {
			safe = isReportSafeSkipLevel(report, i)
		}
		if safe {
			safeReports++
		}
	}
	answer = Stringify(safeReports)
	return
}

func isReportSafe(report []int) bool {
	const maxDiff = 3

	var prevTrend reportTrend
	for i := 0; i < len(report)-1; i++ {
		nextTrend := getReportTrend(report[i], report[i+1])
		if i == 0 {
			prevTrend = nextTrend
		}
		if nextTrend == RepTrend_Eq || nextTrend != prevTrend {
			return false
		}

		if dist(report[i], report[i+1]) > maxDiff {
			return false
		}
	}
	return true
}

func isReportSafeSkipLevel(report []int, skipIdx int) bool {
	const maxDiff = 3

	var prevTrend reportTrend
	initialIdx := 0
	if skipIdx == 0 {
		initialIdx = 1
	}
	lastIdx := len(report) - 1
	if lastIdx == skipIdx {
		// avoid index out of range
		lastIdx--
	}
	for i := 0; i < lastIdx; i++ {
		if skipIdx == i {
			continue
		}
		nextIdx := i + 1
		if nextIdx == skipIdx {
			nextIdx++
		}
		nextTrend := getReportTrend(report[i], report[nextIdx])
		if i == initialIdx {
			prevTrend = nextTrend
		}
		if nextTrend == RepTrend_Eq || nextTrend != prevTrend {
			return false
		}

		if dist(report[i], report[nextIdx]) > maxDiff {
			return false
		}
	}
	return true
}

func getReportTrend(num1, num2 int) reportTrend {
	switch t := num1 - num2; {
	case t < 0:
		return RepTrend_Asc
	case t > 0:
		return RepTrend_Desc
	case t == 0:
		fallthrough
	default:
		return RepTrend_Eq
	}
}
