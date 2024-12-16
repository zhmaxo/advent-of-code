package days

type posInt struct {
	x, y int
}

var (
	dirUp    = posInt{0, -1}
	dirRight = posInt{1, 0}
	dirDown  = posInt{0, 1}
	dirLeft  = posInt{-1, 0}
)

var directions = [4]posInt{
	dirUp, dirRight, dirDown, dirLeft,
}

var dirNames = map[posInt]string{
	dirRight: "right",
	dirDown:  "down",
	dirLeft:  "left",
	dirUp:    "up",
}

var dirTokens = map[posInt]byte{
	dirRight: '>',
	dirDown:  'v',
	dirLeft:  '<',
	dirUp:    '^',
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

func (p posInt) mul(n int) posInt {
	p.x *= n
	p.y *= n
	return p
}

func (p posInt) mul2D(n posInt) posInt {
	p.x *= n.x
	p.y *= n.y
	return p
}

// keep in mind that its all int ops, no fracture
func (p posInt) div(n int) posInt {
	p.x /= n
	p.y /= n
	return p
}

// keep in mind that its all int ops, no fracture
func (p posInt) div2D(n posInt) posInt {
	p.x /= n.x
	p.y /= n.y
	return p
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

func (p posInt) rotateLeft() posInt {
	// 1,0 -> 0,-1 -> -1,0 -> 0,1 -> 1,0
	p.x, p.y = p.y, -p.x
	// not pointer, so original struct is unchanged
	return p
}

func (p posInt) abs() posInt {
	if p.x < 0 {
		p.x *= -1
	}
	if p.y < 0 {
		p.y *= -1
	}
	return p
}

func (p posInt) manhDist(other posInt) uint {
	d := other.sub(p).abs()
	return uint(d.x + d.y)
}
