package days

import "fmt"

func init() {
	DaySolutions[6] = &day6Solution{}
}

type day6Solution struct {
	field fieldDay6
}

func (s *day6Solution) HasData() bool {
	return len(s.field.obstacles) > 0
}

func (s *day6Solution) ReadData(reader ioReader) (err error) {
	const startSize = 2048

	dirMap := map[byte]posInt{
		'^': {0, -1}, '>': {1, 0},
		'v': {0, 1}, '<': {-1, 0},
	}

	s.field.obstacles = make(map[posInt]any, startSize)
	sc := ProcessReader(reader)
	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}
		for i, ch := range line {
			x, y := i, s.field.rect.size.y
			switch ch {
			case '.':
			// just empty position
			case '#':
				// obstacle
				s.field.obstacles[posInt{x, y}] = struct{}{}
			default:
				// check actor
				if dir, ok := dirMap[ch]; ok {
					s.field.actorDir = dir
					s.field.actorPos = posInt{x, y}
				} else {
					err = fmt.Errorf("unexpected token %q", ch)
				}
			}
		}
		// TODO: check consistency
		s.field.rect.size.x = len(line)
		s.field.rect.size.y++
	}

	// fmt.Printf("field {w:%v,h:%v} actor: %v, %v; obstacles: %v\n",
	// 	s.field.width, s.field.height, s.field.actorPos, s.field.actorDir, s.field.obstacles)

	if err == ErrEOF {
		err = nil
	}
	return
}

func (s *day6Solution) SolvePt1() (answer string, err error) {
	const startSize = 512
	visited := make(map[posInt]any, startSize)

	f := s.field // copy to not change original values
	for f.rect.contains(f.actorPos) {
		visited[f.actorPos] = struct{}{}
		_ = f.tick()
	}
	answer = Stringify(len(visited))
	return
}

func (s *day6Solution) SolvePt2() (answer string, err error) {
	const startSize = 2048
	visited := make(map[posInt]map[posInt]any, startSize)
	obstaclesToLoop := make(map[posInt]any, startSize/8)

	f := s.field // copy to not change original values
	for f.rect.contains(f.actorPos) {
		willRetIfObstacle := willReturnToPrevPath(f, visited)
		if vps, ok := visited[f.actorPos]; ok {
			vps[f.actorDir] = struct{}{}
		} else {
			dirs := map[posInt]any{
				f.actorDir: struct{}{},
			}
			visited[f.actorPos] = dirs
		}
		result := f.tick()
		if willRetIfObstacle && result.action == step {
			// fmt.Printf("set obstacle in %v to loop\n", f.actorPos)
			obstaclesToLoop[f.actorPos] = struct{}{}
		}
	}
	// fmt.Printf("possible obstacles: %v", obstaclesToLoop)
	answer = Stringify(len(obstaclesToLoop))
	return
}

type fieldDay6 struct {
	obstacles map[posInt]any
	actorPos  posInt
	actorDir  posInt

	rect rect
}

type actionKind uint8

const (
	step actionKind = iota
	turn
)

type tickResult struct {
	action actionKind
}

func (f *fieldDay6) tick() (result tickResult) {
	f.actorPos, f.actorDir, result.action = f.simulateTick(
		f.actorPos, f.actorDir, f.obstacles)
	return
}

func (_ fieldDay6) simulateTick(pos, dir posInt, obstacles map[posInt]any,
) (nextPos, nextDir posInt, action actionKind) {
	nextPos, nextDir = pos, dir
	checkPos := pos.add(dir)
	if _, ok := obstacles[checkPos]; ok {
		action = turn
		// do turn
		nextDir = dir.rotateRight()
	} else {
		action = step
		// do step
		nextPos = checkPos
	}
	return
}

func wasSavedSameState(pos, dir posInt, m map[posInt]map[posInt]any) bool {
	prevDirs, isVisited := m[pos]
	if !isVisited {
		return false
	}

	if _, sameDir := prevDirs[dir]; sameDir {
		return true
	}
	return false
}

func willReturnToPrevPath(f fieldDay6, visited map[posInt]map[posInt]any) bool {
	simulatelyVisited := make(map[posInt]map[posInt]any, 256)
	pos, dir := f.actorPos, f.actorDir
	// what if we turn now
	dir = dir.rotateRight()
	var action actionKind
	for f.rect.contains(pos) {
		pos, dir, action = f.simulateTick(pos, dir, f.obstacles)
		if action != step {
			if simulatelyVisited[pos] == nil {
				simulatelyVisited[pos] = make(map[posInt]any, 1)
			}
			simulatelyVisited[pos][dir] = struct{}{}
			continue
		}
		//
		if wasSavedSameState(pos, dir, visited) {
			return true
		}
		if wasSavedSameState(pos, dir, simulatelyVisited) {
			// loop to way without returning to prev, only imagined
			return false
		}
		if simulatelyVisited[pos] == nil {
			simulatelyVisited[pos] = make(map[posInt]any, 1)
		}
		simulatelyVisited[pos][dir] = struct{}{}
	}
	return false
}
