package days

type posInt struct {
	x, y int
}

func (p posInt) eq(other posInt) bool {
	return p.x == other.x && p.y == other.y
}

func (p posInt) add(other posInt) posInt {
	return posInt{x: p.x + other.x, y: p.y + other.y}
}

func (p posInt) sub(other posInt) posInt {
	return p.add(other.neg())
}

func (p posInt) neg() posInt {
	return posInt{-p.x, -p.y}
}

func (p posInt) rotateRight() posInt {
	// 1,0 -> 0,1 -> -1,0 -> 0,-1 -> 1,0
	p.x, p.y = -p.y, p.x
	// not pointer, so original struct is unchanged
	return p
}
