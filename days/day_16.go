package days

import (
	"fmt"

	"github.com/beefsack/go-astar"
)

func init() {
	DaySolutions[16] = &day16Solution{}
}

type day16Solution struct {
	cachedTransforms map[posInt]map[posInt]*movePositionKeepPrev
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
	s.cachedTransforms = make(map[posInt]map[posInt]*movePositionKeepPrev, s.rect.area())
	return
}

func (s *day16Solution) SolvePt1() (answer string, err error) {
	startDir := posInt{1, 0} // right
	startNode := movePositionKeepPrev{
		parent: s,
		pos:    s.start,
		prev:   s.start.sub(startDir),
	}
	for _, n := range s.finish.neighbors() {
		if s.getValueAt(n) == s.wall {
			continue
		}
		finishNode := movePositionKeepPrev{
			parent: s,
			pos:    s.finish,
			prev:   n,
		}
		path, dist, found := astar.Path(&startNode, &finishNode)
		if !found {
			plf("not found path %v -> %v", startNode, finishNode)
			continue
		}
		m := make(map[posInt]byte, len(path))
		for _, p := range path {
			mpkp, ok := p.(*movePositionKeepPrev)
			if !ok {
				plf("UNEXPECTED POINTER VALUE %v", p)
				continue
			}
			m[mpkp.pos] = mpkp.dirToken()
		}
		plf("found path with dist %v:", dist)
		printRectFunc(s.rect, func(p posInt) rune {
			if t, ok := m[p]; ok {
				return rune(t)
			}
			return rune(s.getValueAt(p))
		})
	}
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

type movePositionKeepPrev struct {
	parent *day16Solution

	cachedPathNeighbors []astar.Pather

	pos, prev posInt
}

func (mpkp movePositionKeepPrev) dirToken() byte {
	if t, ok := dirTokens[mpkp.pos.sub(mpkp.prev)]; ok {
		return t
	}
	return '?'
}

func (mpkp movePositionKeepPrev) String() string {
	return fmt.Sprintf("%q {pos: %v come from %v}", mpkp.dirToken(), mpkp.pos, mpkp.prev)
}

func (mpkp movePositionKeepPrev) neighborPositions() (neighbors [4]posInt) {
	// no need to include prev position
	// I want to keep specific order so won't use neighbors()
	// positions := mpkp.pos.neighbors()
	mainDir := mpkp.pos.sub(mpkp.prev)
	neighbors = [4]posInt{
		mpkp.pos.add(mainDir),
		mpkp.pos.add(mainDir.rotateRight()),
		mpkp.pos.add(mainDir.rotateLeft()),
		mpkp.pos.sub(mainDir),
	}
	return
}

func (mpkp movePositionKeepPrev) costCoefficient(p posInt) int {
	// no need to include prev position
	// I want to keep specific order so won't use neighbors()
	// positions := mpkp.pos.neighbors()
	mainDir := mpkp.pos.sub(mpkp.prev)
	switch p {
	case mpkp.pos.add(mainDir):
		return 0
	case mpkp.pos.add(mainDir.rotateRight()):
		return 1
	case mpkp.pos.add(mainDir.rotateLeft()):
		return 2
	case mpkp.pos.sub(mainDir):
		return 0
	default:
		return -1
	}
}

func (mpkp *movePositionKeepPrev) PathNeighbors() []astar.Pather {
	plf("check neighbors for %v", *mpkp)
	parent := mpkp.parent
	if parent == nil {
		// should've be an error but so what
		return nil
	}
	if mpkp.cachedPathNeighbors != nil {
		// DEBUG STUFF
		printRectFunc(parent.rect, func(p posInt) rune {
			if p == mpkp.pos {
				return rune(dirTokens[mpkp.pos.sub(mpkp.prev)])
			}
			if p.manhDist(mpkp.pos) == 1 && parent.getValueAt(p) == parent.empty {
				// exactly free neighbor position
				for _, v := range mpkp.cachedPathNeighbors {
					other := v.(*movePositionKeepPrev)
					if other.pos == p {
						return rune(other.dirToken())
					}
				}
				return '?'
			}
			return rune(parent.posToken(p))
		})

		return mpkp.cachedPathNeighbors
	}

	neighbors := mpkp.neighborPositions()
	mpkp.cachedPathNeighbors = make([]astar.Pather, 0, len(neighbors))
	// {4,1} -> {5,1} : pos.sub(prev) = {1,0} :
	for _, n := range neighbors {
		if !parent.rect.contains(n) {
			plf("unexpected neighbor for %v (%q) -> %v is out of %v",
				mpkp.pos, parent.getValueAt(mpkp.pos), n, parent.rect)
			continue
		}
		if parent.getValueAt(n) == parent.wall {
			// we shouldn't count walls as pather neighbors
			continue
		}
		mapL1 := parent.cachedTransforms[n]
		var neighbor *movePositionKeepPrev
		if mapL1 == nil {
			mapL1 = make(map[posInt]*movePositionKeepPrev, 4)
			parent.cachedTransforms[n] = mapL1
		}
		neighbor = mapL1[n.sub(mpkp.pos)]
		if neighbor == nil {
			neighbor = &movePositionKeepPrev{
				pos:    n,
				prev:   mpkp.pos,
				parent: parent,
			}
			mapL1[n.sub(mpkp.pos)] = neighbor
		}
		mpkp.cachedPathNeighbors = append(mpkp.cachedPathNeighbors, neighbor)
	}

	// DEBUG STUFF
	printRectFunc(parent.rect, func(p posInt) rune {
		if p == mpkp.pos {
			return rune(dirTokens[mpkp.pos.sub(mpkp.prev)])
		}
		if p.manhDist(mpkp.pos) == 1 && parent.getValueAt(p) == parent.empty {
			// exactly free neighbor position
			for _, v := range mpkp.cachedPathNeighbors {
				other := v.(*movePositionKeepPrev)
				if other.pos == p {
					return rune(other.dirToken())
				}
			}
			return '?'
		}
		return rune(parent.posToken(p))
	})

	return mpkp.cachedPathNeighbors
}

func (mpkp movePositionKeepPrev) PathNeighborCost(to astar.Pather) float64 {
	const (
		rotateCost = 1_000
	)
	if mpkp.cachedPathNeighbors == nil {
		return 1_000_000_000
	}
	for _, n := range mpkp.cachedPathNeighbors {
		if n == to {
			other := n.(*movePositionKeepPrev)
			// first neighbor has 0, others has 1
			i := mpkp.costCoefficient(other.pos)
			if i < 0 {
				return 1_000_000_000
			}
			coeff := (i + 1) / 2
			res := float64(rotateCost*coeff + 1)
			plf("cost from %v to %v is %v", mpkp, n, res)
			printRectFunc(mpkp.parent.rect, func(p posInt) rune {
				switch p {
				case mpkp.pos:
					return rune(mpkp.dirToken())
				case other.pos:
					return rune(other.dirToken())
				default:
					return rune(mpkp.parent.posToken(p))
				}
			})
			return res
		}
	}
	plf("cant find neighbor like %v for %v", to, mpkp)
	return 1_000_000_000
}

func (mpkp movePositionKeepPrev) PathEstimatedCost(to astar.Pather) float64 {
	// use manhattan distance
	if other, ok := to.(*movePositionKeepPrev); ok {
		if other == nil {
			plf("UNEXPECTED NIL POINTER %v", other)
			return 1_000_000_000
		}
		dist := mpkp.pos.manhDist(other.pos)
		plf("manh dist from %v to %v is %v", mpkp, other, dist)
		return float64(dist)
	}
	return 1_000_000_000
}
