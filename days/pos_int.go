package days

type posInt struct {
	x, y int
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
