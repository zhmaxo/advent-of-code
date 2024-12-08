package days

type rect struct {
	topLeft, size posInt
}

// exclude upper bound
func (r rect) contains(p posInt) bool {
	if r.topLeft.x > p.x || r.topLeft.y > p.y {
		return false
	}

	if r.size.x <= p.x || r.size.y <= p.y {
		return false
	}
	return true
}
