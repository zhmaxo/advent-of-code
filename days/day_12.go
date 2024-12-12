package days

import "fmt"

func init() {
	DaySolutions[12] = &day12Solution{}
}

type day12Solution struct {
	regCache map[byte][]posInt

	field [][]byte
	rect  rect
}

func (s *day12Solution) HasData() bool {
	return true
}

func (s *day12Solution) ReadData(reader ioReader) (err error) {
	s.regCache = make(map[byte][]posInt, 26) // alphabet len
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()

	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		// consistency assurance
		if s.rect.size.x == 0 {
			s.rect.size.x = len(line)
		} else if len(line) != s.rect.size.x {
			return fmt.Errorf("inconsistent line lens: %v but expected %v",
				len(line), s.rect.size.x)
		}

		row := make([]byte, len(line))
		for i, b := range line {
			row[i] = b
			s.regCache[b] = append(s.regCache[b], posInt{i, s.rect.size.y})
		}
		s.field = append(s.field, row)
		s.rect.size.y++
	}
	return
}

func (s *day12Solution) SolvePt1() (answer string, err error) {
	result := uint64(0)
	visited := make(map[posInt]any, s.rect.size.x*s.rect.size.y)
	for y, row := range s.field {
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
	plant := s.field[p.y][p.x]
	for _, n := range p.neighbors() {
		if !s.rect.contains(n) {
			// out of field
			perimeter++
			continue
		}
		if s.field[n.y][n.x] != plant {
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

func dep(s day12Solution) {
	// 0 is no region
	type regId byte
	// how many other plants around this one
	type plotPerim = byte
	// grouped plants with cached perimeter
	type region struct {
		// count local not touching with other reg plant sides
		plants map[posInt]plotPerim
	}

	regions := make(map[regId]region, 256)
	thisPlantRegions := make(map[regId]region, 128)
	regionAffinity := make(map[regId][]regId, 128)
	visited := make(map[posInt]regId, 1024)

	for plant, positions := range s.regCache {
		for _, pos := range positions {
			var rid regId
			var per plotPerim
			for _, n := range pos.neighbors() {
				// if out of field - just increase perimeter
				if !s.rect.contains(n) {
					per++
					continue
				}

				// other plant, skip
				if s.field[pos.y][pos.x] != plant {
					per++
					continue
				}
				// check visited
				if reg, ok := visited[n]; ok {
					// check existed
					if rid != reg && rid > 0 {
						// add affinity
						from, to := rid, reg
						// from lesser to bigger
						if reg < rid {
							from, to = to, from
						}
						regionAffinity[from] = append(regionAffinity[from], to)
					} else {
						// add to region
						rid = reg
					}
				}
			}
			if rid == 0 {
				// create reg
				rid = regId(len(thisPlantRegions) + 1)
			}
		}
		// TODO: add regions to cache
		for k, v := range thisPlantRegions {
			// TODO: rewrite
			regions[k] = v
		}
		// clear temporary caches, prepare for new cycle
		clear(thisPlantRegions)
		clear(regionAffinity)
	}
}
