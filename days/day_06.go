package days

import "fmt"

func init() {
	DaySolutions[6] = &day6Solution{}
}

type day6Solution struct {
	field field
}

func (s *day6Solution) HasData() bool {
	return true
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
			x, y := i, s.field.height
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
		s.field.width = len(line)
		s.field.height++
	}

	fmt.Printf("field {w:%v,h:%v} actor: %v, %v; obstacles: %v\n",
		s.field.width, s.field.height, s.field.actorPos, s.field.actorDir, s.field.obstacles)

	if err == ErrEOF {
		err = nil
	}
	return
}

func (s *day6Solution) SolvePt1() (answer string, err error) {
	const startSize = 512
	visited := make(map[posInt]any, startSize)

	f := s.field // copy to not change original values
	for f.actorPos.within(f.width, f.height) {
		visited[f.actorPos] = struct{}{}
		_ = f.tick()
	}
	answer = Stringify(len(visited))
	return
}

func (s *day6Solution) SolvePt2() (answer string, err error) {
	return
}

type field struct {
	obstacles map[posInt]any
	actorPos  posInt
	actorDir  posInt

	height, width int
}

type posInt struct {
	x, y int
}

var dirNames = map[posInt]string{
	{1, 0}:  "right",
	{0, 1}:  "down",
	{-1, 0}: "left",
	{0, -1}: "up",
}

func (p posInt) add(other posInt) posInt {
	return posInt{x: p.x + other.x, y: p.y + other.y}
}

func (p posInt) rotateRight() posInt {
	// 1,0 -> 0,1 -> -1,0 -> 0,-1 -> 1,0
	p.x, p.y = -p.y, p.x
	// not pointer, so original struct is unchanged
	return p
}

// 0,0 is top-left; for more flexible logic it may be rect struct as area
func (p posInt) within(width, height int) bool {
	return p.x >= 0 && p.y >= 0 && p.x < width && p.y < height
}

type actionKind uint8

const (
	step actionKind = iota
	turn
)

type tickResult struct {
	action actionKind
}

func (f *field) tick() (result tickResult) {
	nextPos := f.actorPos.add(f.actorDir)
	if _, ok := f.obstacles[nextPos]; ok {
		result.action = turn
		// do turn
		f.actorDir = f.actorDir.rotateRight()
	} else {
		result.action = step
		// do step
		f.actorPos = nextPos
	}
	return
}
