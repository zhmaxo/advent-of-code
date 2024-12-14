package days

import "fmt"

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

type byteField struct {
	data [][]byte
	rect rect
}

func scanField(reader ioReader) (field byteField, err error) {
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()

	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		// consistency assurance
		if field.rect.size.x == 0 {
			field.rect.size.x = len(line)
		} else if len(line) != field.rect.size.x {
			err = fmt.Errorf("inconsistent line lens: %v but expected %v",
				len(line), field.rect.size.x)
			return
		}

		row := make([]byte, len(line))
		copy(row, line)
		field.data = append(field.data, row)
		field.rect.size.y++
	}
	return
}
