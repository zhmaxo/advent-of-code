package days

func init() {
	DaySolutions[8] = &day8Solution{}
}

type day8Solution struct {
	antennas map[byte][]posInt
	area     rect

	totalAntennas int
}

func (s *day8Solution) HasData() bool {
	return len(s.antennas) > 0
}

func (s *day8Solution) ReadData(reader ioReader) (err error) {
	const startSize = 64
	s.antennas = make(map[byte][]posInt, startSize)

	defer func() { err = RefineError(err) }()
	sc := ProcessReader(reader)
	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}
		y := s.area.size.y
		for i, v := range line {
			if v == '.' {
				continue
			}
			x := i
			s.antennas[v] = append(s.antennas[v], posInt{x, y})
			s.totalAntennas++
		}
		s.area.size.x = len(line)
		s.area.size.y++
	}
	return
}

func (s *day8Solution) SolvePt1() (answer string, err error) {
	antinodes := make(map[posInt]any, s.totalAntennas)
	for _, ag := range s.antennas {
		for i := 0; i < len(ag); i++ {
			for j := i + 1; j < len(ag); j++ {
				an1, an2 := findAntinodes(ag[i], ag[j])
				if s.area.contains(an1) {
					antinodes[an1] = well
				}
				if s.area.contains(an2) {
					antinodes[an2] = well
				}
			}
		}
	}
	answer = Stringify(len(antinodes))
	return
}

func (s *day8Solution) SolvePt2() (answer string, err error) {
	antinodes := make(map[posInt]any, s.totalAntennas)
	for _, ag := range s.antennas {
		if len(ag) < 2 {
			continue
		}
		for i := 0; i < len(ag); i++ {
			antinodes[ag[i]] = well
			for j := i + 1; j < len(ag); j++ {
				dir1 := ag[i].sub(ag[j])
				an1 := ag[i].add(dir1)
				for s.area.contains(an1) {
					antinodes[an1] = well
					an1 = an1.add(dir1)
				}
				dir2 := ag[j].sub(ag[i])
				an2 := ag[j].add(dir2)
				for s.area.contains(an2) {
					antinodes[an2] = well
					an2 = an2.add(dir2)
				}
			}
		}
	}
	answer = Stringify(len(antinodes))
	return
}

func findAntinodes(a1, a2 posInt) (an1, an2 posInt) {
	dir1 := a1.sub(a2)
	dir2 := a2.sub(a1)
	return a1.add(dir1), a2.add(dir2)
}
