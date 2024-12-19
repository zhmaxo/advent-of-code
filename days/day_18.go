package days

import (
	"fmt"

	"github.com/beefsack/go-astar"
)

func init() {
	DaySolutions[18] = &day18Solution{}
}

type day18Solution struct {
	fallingBytes  []posInt
	simulateBytes []posInt
	walkableField
}

func (s *day18Solution) HasData() bool {
	return true
}

func (s *day18Solution) ReadData(reader ioReader) (err error) {
	const (
		size     = 71
		simTicks = 1024
	)
	// r := regexp.MustCompile("([0-9]{1,}),([0-9]{1,})")

	s.fallingBytes = make([]posInt, 0, simTicks)
	s.byteField = newField(rect{size: posInt{1, 1}.mul(size)})
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()
	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		var nums []int
		nums, err = ParseNumbersSep(string(line), ',')
		if err != nil {
			break
		}
		if len(nums) != 2 {
			continue
		}
		p := posInt{nums[0], nums[1]}
		s.fallingBytes = append(s.fallingBytes, p)
	}
	simCount := simTicks
	if simCount > len(s.fallingBytes) {
		simCount = len(s.fallingBytes)
	}
	s.simulateBytes = s.fallingBytes[:simCount]
	return
}

func (s *day18Solution) SolvePt1() (answer string, err error) {
	const (
		available = 0
		corrupted = 1
	)
	s.passable, s.obstacle = available, corrupted
	for i := range s.simulateBytes {
		p := s.simulateBytes[i]
		s.setValueAt(p, corrupted)
	}

	start, ok := s.getNode(posInt{0, 0})
	if !ok {
		err = fmt.Errorf("unexpected impassable start %v", start)
		return
	}
	finish, ok := s.getNode(s.rect.size.sub(posInt{1, 1}))
	if !ok {
		err = fmt.Errorf("unexpected impassable finish %v", finish)
		return
	}
	path, dist, found := astar.Path(d18Node{start}, d18Node{finish})
	if !found {
		err = fmt.Errorf("not found path %v -> %v", start, finish)
		return
	}

	s.printPathInfo(path)
	answer = Stringify(dist)
	return
}

func (s *day18Solution) SolvePt2() (answer string, err error) {
	const (
		passable = 0 // in that part obstacle is > 0
		size     = 71
	)
	cleanField := [size][size]byte{}
	checkField := cleanField

	s.data = make([][]byte, s.rect.size.y)
	for i := range s.rect.size.y {
		// now field resets on every checkField = cleanField
		// and no alloc here!
		s.data[i] = checkField[i][:s.rect.size.x]
	}

	start, ok := s.getNode(posInt{0, 0})
	if !ok {
		err = fmt.Errorf("unexpected impassable start %v", start)
		return
	}
	finish, ok := s.getNode(s.rect.size.sub(posInt{1, 1}))
	if !ok {
		err = fmt.Errorf("unexpected impassable finish %v", finish)
		return
	}

	from, to := d18Node{start}, d18Node{finish}

	path, _, found := astar.Path(from, to)
	if !found {
		printRectFunc(s.rect, func(p posInt) rune { return rune(s.getValueAt(p)) })
		return "", fmt.Errorf("invalid data: pathfinding failure before first byte fall")
	}
	for i, p := range s.fallingBytes {
		plf("check byte at %v", p)
		lastPath := path
		v := s.getValueAt(p) + 1
		s.setValueAt(p, v)
		path, _, found = astar.Path(from, to)

		if !found {
			plf("falling byte at %v tick", i)
			s.printPathInfo(lastPath)
			answer = fmt.Sprintf("%v,%v", p.x, p.y)
			break
		}
	}
	return
}

func (s *day18Solution) simulateTick(since, before int) error {
	step := 1
	if since > before {
		step = -1
	}

	for i := since; i != before; i += step {
		p := s.fallingBytes[i]
		v := int8(s.getValueAt(p))
		v += int8(step)
		if v < 0 {
			return fmt.Errorf("unexpected change value at %v: %v -> %v",
				p, s.getValueAt(p), v)
		}
		s.setValueAt(p, byte(v))
	}
	return nil
}

func (s *day18Solution) ModifyForTest() {
	const (
		size = 7
		sims = 12
	)
	s.byteField = newField(rect{size: posInt{1, 1}.mul(size)})
	s.simulateBytes = s.fallingBytes[:sims]
}

func (s day18Solution) printPathInfo(path []astar.Pather) {
	m := make(map[posInt]rune, len(path))
	for _, n := range path {
		m[n.(d18Node).posInt] = 'O'
	}
	printRectFunc(s.rect, func(p posInt) rune {
		if r, ok := m[p]; ok {
			return r
		}
		v := s.getValueAt(p)
		if v == s.obstacle {
			return '#'
		}
		return '.'
	})
}

type d18Node struct {
	walkableNode
}

func (this d18Node) PathNeighbors() []astar.Pather {
	n, p, c := this.passableNeighbors()
	result := make([]astar.Pather, 0, c)
	for i := range p {
		if p[i] {
			result = append(result, d18Node{this.gotoNode(n[i])})
		}
	}
	return result
}

func (this d18Node) PathNeighborCost(to astar.Pather) float64 {
	return float64(this.manhDist(to.(d18Node).posInt))
}

func (this d18Node) PathEstimatedCost(to astar.Pather) float64 {
	return float64(this.manhDist(to.(d18Node).posInt))
}
