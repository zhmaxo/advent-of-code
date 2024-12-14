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

type byteField struct {
	data [][]byte
	rect rect
}

func (f byteField) getValueAt(p posInt) byte {
	return f.data[p.y][p.x]
}

func scanFieldFunc(reader ioReader, f func(line []byte, lineIdx int) ([]byte, error),
) (field byteField, err error) {
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()

	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		var row []byte
		row, err = f(line, field.rect.size.y)
		field.data = append(field.data, row)
		field.rect.size.y++

		// consistency assurance
		if field.rect.size.x == 0 {
			field.rect.size.x = len(row)
		} else if len(row) != field.rect.size.x {
			err = fmt.Errorf("inconsistent line lens: %v but expected %v",
				len(row), field.rect.size.x)
			return
		}
	}
	return
}

func scanFieldBytesAsIs(reader ioReader) (field byteField, err error) {
	return scanFieldFunc(reader,
		func(line []byte, _ int) ([]byte, error) {
			row := make([]byte, len(line))
			copy(row, line)
			return row, nil
		})
}
