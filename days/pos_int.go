package days

type posInt struct {
	x, y int
}

var dirNames = map[posInt]string{
	{1, 0}:  "right",
	{0, 1}:  "down",
	{-1, 0}: "left",
	{0, -1}: "up",
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

func (p posInt) neighbors() (result [4]posInt) {
	i := 0
	for k := range dirNames {
		result[i] = p.add(k)
		i++
	}
	return
}

func (p posInt) rotateRight() posInt {
	// 1,0 -> 0,1 -> -1,0 -> 0,-1 -> 1,0
	p.x, p.y = -p.y, p.x
	// not pointer, so original struct is unchanged
	return p
}
