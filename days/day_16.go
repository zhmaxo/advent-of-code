package days

import (
	"fmt"

	"github.com/beefsack/go-astar"
)

func init() {
	DaySolutions[16] = &day16Solution{}
}

type day16Solution struct {
	byteField
	start, finish posInt
	wall, empty   byte
}

func (s *day16Solution) HasData() bool {
	return true
}

func (s *day16Solution) ReadData(reader ioReader) (err error) {
	const (
		wall   = '#'
		empty  = '.'
		start  = 'S'
		finish = 'E'
	)
	s.wall, s.empty = wall, empty
	s.byteField, err = scanFieldFunc(reader, func(line []byte, lineIdx int) (row []byte, err error) {
		row = make([]byte, len(line))
		for i, v := range line {
			if v == wall {
				row[i] = s.wall
			} else {
				switch v {
				case start:
					s.start = posInt{i, lineIdx}
				case finish:
					s.finish = posInt{i, lineIdx}
				case empty:
				default:
					plf("unexpected token %q", v)
				}
				row[i] = s.empty
			}
		}
		return
	})
	return
}

func (s *day16Solution) SolvePt1() (answer string, err error) {
	startNode := d16MoveNode{parent: s, posInt: s.start}
	targetNodeH := startNode
	targetNodeH.posInt = s.finish
	targetNodeV := targetNodeH
	targetNodeV.vertical = true
	possibleTargets := [2]d16MoveNode{
		targetNodeH, targetNodeV,
	}
	result := 0
	for _, finishNode := range possibleTargets {
		path, dist, found := astar.Path(startNode, finishNode)
		if !found {
			plf("not found path %v -> %v", startNode, finishNode)
			continue
		}
		if result == 0 || result > int(dist) {
			result = int(dist)
		}
		m := make(map[posInt]rune, len(path))
		for _, p := range path {
			mpkp, ok := p.(d16MoveNode)
			if !ok {
				plf("UNEXPECTED POINTER VALUE %v", p)
				continue
			}
			r := []rune(mpkp.String())[0]
			if t, has := m[mpkp.posInt]; has {
				if t != r {
					m[mpkp.posInt] = '+'
				}
			} else {
				m[mpkp.posInt] = r
			}
		}
		plf("found path with dist %v:", dist)
		printRectFunc(s.rect, func(p posInt) rune {
			if t, ok := m[p]; ok {
				return t
			}
			return rune(s.getValueAt(p))
		})
	}
	answer = Stringify(result)
	return
}

func (s *day16Solution) SolvePt2() (answer string, err error) {
	return
}

func (s *day16Solution) posToken(p posInt) byte {
	if p == s.start {
		return 'S'
	}
	if p == s.finish {
		return 'E'
	}
	return s.getValueAt(p)
}

type d16MoveNode struct {
	parent *day16Solution
	posInt
	vertical bool
}

func (this d16MoveNode) String() string {
	sign := "-"
	if this.vertical {
		sign = "|"
	}
	return fmt.Sprintf("%v {%v}", sign, this.posInt)
}

func (this d16MoveNode) rotate() d16MoveNode {
	this.vertical = !this.vertical
	return this
}

func (this d16MoveNode) translate(toTopLeft bool) d16MoveNode {
	direction := dirRight
	if this.vertical {
		direction = dirUp
	}
	if toTopLeft {
		direction = direction.neg()
	}
	this.posInt = this.add(direction)
	return this
}

func (this d16MoveNode) isReachable() bool {
	p := this.parent
	return p != nil &&
		p.rect.contains(this.posInt) &&
		p.getValueAt(this.posInt) == p.empty
}

func (this d16MoveNode) PathNeighbors() []astar.Pather {
	parent := this.parent
	if parent == nil {
		// should've be an error but so what
		return nil
	}

	count := 1
	t1 := this.translate(true)
	t1Valid := t1.isReachable()
	if t1Valid {
		count++
	}
	t2 := this.translate(false)
	t2Valid := t2.isReachable()
	if t2Valid {
		count++
	}
	neighbors := make([]astar.Pather, count)
	idx := 0
	neighbors[idx] = this.rotate()
	idx++
	if t1Valid {
		neighbors[idx] = t1
		idx++
	}
	if t2Valid {
		neighbors[idx] = t2
	}

	plf("check neighbors for %v: %v", this, neighbors)
	// DEBUG STUFF
	printRectFunc(parent.rect, func(p posInt) rune {
		switch p {
		case this.posInt:
			return '+'
		case this.translate(false).posInt:
			return []rune(this.translate(false).String())[0]
		case this.translate(true).posInt:
			return []rune(this.translate(true).String())[0]
		default:
			return rune(parent.posToken(p))
		}
	})

	return neighbors
}

func (this d16MoveNode) PathNeighborCost(to astar.Pather) float64 {
	const (
		rotateCost = 1_000
	)
	other, ok := to.(d16MoveNode)
	if !ok {
		return rotateCost * 1_000_000_000
	}
	cost := this.manhDist(other.posInt)
	if this.vertical != other.vertical {
		cost += rotateCost
	}
	plf("cost from %v to %v is %v", this, other, cost)
	printRectFunc(this.parent.rect, func(p posInt) rune {
		switch p {
		case this.posInt:
			if this.eq(other.posInt) {
				return '+'
			}
			return []rune(this.String())[0]
		case other.posInt:
			return []rune(other.String())[0]
		default:
			return rune(this.parent.posToken(p))
		}
	})
	return float64(cost)
}

func (this d16MoveNode) PathEstimatedCost(to astar.Pather) float64 {
	// use manhattan distance
	if other, ok := to.(d16MoveNode); ok {
		dist := this.manhDist(other.posInt)
		// plf("manh dist from %v to %v is %v", this, other, dist)
		return float64(dist)
	}
	return 1_000_000_000
}
