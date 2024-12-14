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
