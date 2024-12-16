package days

func init() {
	DaySolutions[15] = &day15Solution{}
}

type day15Solution struct {
	byteField
}

func (s *day15Solution) HasData() bool {
	return true
}

func (s *day15Solution) ReadData(reader ioReader) (err error) {
	return
}

func (s *day15Solution) SolvePt1() (answer string, err error) {
	return
}

func (s *day15Solution) SolvePt2() (answer string, err error) {
	return
}

func (s *day15Solution) tryMove(pos, dir posInt) bool {
	const (
		wall = iota
		box
	)

	type moveOp struct {
		from, to posInt
	}

	moveOps := make([]moveOp, 0, 16)

	checkPos := pos.add(dir)
	for s.rect.contains(checkPos) {
		// wall -> false
		// empty -> true
		// box -> append to move ops
		moveOps = append(moveOps, moveOp{checkPos, checkPos.add(dir)})
		checkPos = checkPos.add(dir)
	}
	return false
}
