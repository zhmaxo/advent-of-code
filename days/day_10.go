package days

func init() {
	DaySolutions[10] = &day10Solution{}
}

type day10Solution struct {
	pathStarters []posInt
	byteField
}

func (s *day10Solution) HasData() bool {
	return true
}

func (s *day10Solution) ReadData(reader ioReader) (err error) {
	s.byteField, err = scanFieldFunc(reader, s.convertLine)
	printRectFunc(s.rect, func(p posInt) rune {
		return rune('0' + s.getValueAt(p))
	})
	return
}

func (s *day10Solution) SolvePt1() (answer string, err error) {
	visited := make(map[posInt]any, s.rect.area())
	score := 0
	for _, p := range s.pathStarters {
		pathScore := s.calcScore(p, visited)
		score += pathScore
		plf("path %v has score %v", visited, pathScore)
		clear(visited)
	}
	answer = Stringify(score)
	return
}

func (s *day10Solution) SolvePt2() (answer string, err error) {
	return
}

func (s *day10Solution) convertLine(line []byte, lineIdx int,
) (row []byte, err error) {
	const zero = '0'
	row = make([]byte, len(line))
	for i := 0; i < len(row); i++ {
		b := line[i]
		r := b - '0'
		row[i] = r
		if r == 0 {
			s.pathStarters = append(s.pathStarters, posInt{i, lineIdx})
		}
	}
	return
}

func (s *day10Solution) canMove(from, to posInt) bool {
	h1, h2 := s.getValueAt(from), s.getValueAt(to)
	return h2-h1 == 1
}

func (s *day10Solution) calcScore(start posInt, visited map[posInt]any,
) (score int) {
	const targetAltitude = 9
	if s.getValueAt(start) == targetAltitude {
		return 1
	}
	for _, n := range start.neighbors() {
		if _, beenHere := visited[n]; beenHere {
			continue
		}
		if !s.rect.contains(n) {
			continue
		}
		if !s.canMove(start, n) {
			continue
		}
		visited[n] = s.getValueAt(n)
		score += s.calcScore(n, visited)
	}
	return
}
