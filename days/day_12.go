package days

func init() {
	DaySolutions[12] = &day12Solution{}
}

type day12Solution struct {
	byteField
}

func (s *day12Solution) HasData() bool {
	return true
}

func (s *day12Solution) ReadData(reader ioReader) (err error) {
	s.byteField, err = scanFieldBytesAsIs(reader)
	return
}

func (s *day12Solution) SolvePt1() (answer string, err error) {
	result := uint64(0)
	visited := make(map[posInt]any, s.rect.size.x*s.rect.size.y)
	for y, row := range s.data {
		for x := range row {
			p := posInt{x, y}
			if _, hasBeen := visited[p]; hasBeen {
				// skip visited plots
				continue
			}
			area, perimeter := s.collectRegionInfo(p, visited)
			result += uint64(area * perimeter)
		}
	}
	answer = Stringify(result)
	return
}

func (s *day12Solution) collectRegionInfo(p posInt, visited map[posInt]any,
) (area, perimeter uint) {
	visited[p] = well // mark as visited
	area = 1          // this plot area
	plant := s.data[p.y][p.x]
	for _, n := range p.neighbors() {
		if !s.rect.contains(n) {
			// out of field
			perimeter++
			continue
		}
		if s.data[n.y][n.x] != plant {
			// another plant
			perimeter++
			continue
		}
		if _, hasBeen := visited[n]; hasBeen {
			// skip visited plots
			continue
		}
		sub_area, sub_perimeter := s.collectRegionInfo(n, visited)
		// summarize
		area += sub_area
		perimeter += sub_perimeter
	}
	return
}

func (s *day12Solution) SolvePt2() (answer string, err error) {
	return
}
