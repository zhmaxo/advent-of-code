package days

import (
	"fmt"
	"regexp"
)

func init() {
	DaySolutions[14] = &day14Solution{}
}

type day14Solution struct {
	moveList []robotMoveInfo

	rect rect
}

type robotMoveInfo struct {
	p, v posInt
}

func (s *day14Solution) HasData() bool {
	return true
}

func (s *day14Solution) ReadData(reader ioReader) (err error) {
	const (
		startSize         = 1024
		entryPattern      = `p=(-?[0-9]{1,}),(-?[0-9]{1,}) v=(-?[0-9]{1,}),(-?[0-9]{1,})`
		parsedNumsMustLen = 4
		matchSureLen      = parsedNumsMustLen + 1
		fieldHeight       = 103
		fieldWidth        = 101
	)
	s.moveList = make([]robotMoveInfo, 0, startSize)
	s.rect.size = posInt{fieldWidth, fieldHeight}
	r := regexp.MustCompile(entryPattern)
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()

	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		e := robotMoveInfo{}
		match := r.FindSubmatch(line)
		if len(match) != matchSureLen {
			err = fmt.Errorf("regexp match result len mismatch: %q", match)
			return
		}
		var nums []int
		nums, err = ParseNumbersEachBytes(match[1:]...)
		if len(nums) != parsedNumsMustLen {
			err = fmt.Errorf("parsed numbers len mismatch: %v", match)
			return
		}
		e.p.x, e.p.y, e.v.x, e.v.y = nums[0], nums[1], nums[2], nums[3]
		s.moveList = append(s.moveList, e)
	}
	return
}

func (s *day14Solution) ModifyForTest() {
	s.rect.size.x, s.rect.size.y = 11, 7
}

func (s *day14Solution) SolvePt1() (answer string, err error) {
	const ticksToSim = 100
	// TODO: copy, not change directly
	moveList := make([]robotMoveInfo, len(s.moveList))
	copy(moveList, s.moveList)
	for range ticksToSim {
		s.simulateTick(moveList, s.rect)
	}

	halfSize := s.rect.size.div(2)
	offset := s.rect.size.sub(halfSize)
	var quadrants [4]rect
	for i := range quadrants {
		const mask = 0b01
		// (0,0), (1,0), (0,1), (1,1)
		mul := posInt{
			x: i & mask,
			y: (i >> 1) & mask,
		}
		quadrants[i] = rect{
			// quadrants differs only by offset
			topLeft: offset.mul2D(mul),
			size:    halfSize,
		}
		pf("q%v: %v\n", i, quadrants[i])
	}

	var securityScore [4]uint
	poss := make([]posInt, 0, len(s.moveList))
	for _, v := range s.moveList {
		poss = append(poss, v.p)
		isInsideOfQ := false
		for i, q := range quadrants {
			if q.contains(v.p) {
				securityScore[i]++
				isInsideOfQ = true
				break
			}
		}
		if !isInsideOfQ {
			plf("not found q for %v", v.p)
		}
	}

	result := uint(1)
	for _, score := range securityScore {
		result *= score
	}
	answer = Stringify(result)
	return
}

func (s *day14Solution) SolvePt2() (answer string, err error) {
	return
}

func (_ *day14Solution) simulateTick(moveList []robotMoveInfo, rect rect) {
	for i := range moveList {
		ri := moveList[i]
		// move
		p := ri.p.add(ri.v)
		ri.p = rect.loopInside(p) // tp if out of bounds
		moveList[i] = ri
	}
}
