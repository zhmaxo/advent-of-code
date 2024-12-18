package days

type rect struct {
	topLeft, size posInt
}

// exclude higher bound
func (r rect) contains(p posInt) bool {
	tl := r.topLeft
	br := r.size.add(tl)
	if tl.x > p.x || tl.y > p.y {
		return false
	}

	if br.x <= p.x || br.y <= p.y {
		return false
	}
	return true
}

func (r rect) area() int {
	return r.size.x * r.size.y
}

func (r rect) loopInside(p posInt) posInt {
	if r.size.x < 1 || r.size.y < 1 {
		// can't loop inside the invalid rect
		return p
	}

	tl := r.topLeft
	br := r.size.add(tl)

	for p.x < tl.x {
		p.x += r.size.x
	}
	for p.x >= br.x {
		p.x -= r.size.x
	}

	for p.y < tl.y {
		p.y += r.size.y
	}
	for p.y >= br.y {
		p.y -= r.size.y
	}
	return p
}
