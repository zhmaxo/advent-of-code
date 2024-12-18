package days

import "fmt"

type byteField struct {
	data [][]byte
	rect rect
}

func (f byteField) getValueAt(p posInt) byte {
	return f.data[p.y][p.x]
}

func (f byteField) setValueAt(p posInt, v byte) {
	f.data[p.y][p.x] = v
}

func newField(area rect) byteField {
	f := byteField{rect: area}
	f.data = make([][]byte, area.size.y)
	for i := range f.data {
		f.data[i] = make([]byte, area.size.x)
	}
	return f
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

type walkableField struct {
	byteField
	passable, obstacle byte
}

func (this walkableField) isPassable(p posInt) bool {
	// NOTE: or != this.obstacle ?
	return this.rect.contains(p) && this.getValueAt(p) == this.passable
}

func (this *walkableField) getNode(p posInt) (walkableNode, bool) {
	node := walkableNode{parent: this, posInt: p}
	return node, this.isPassable(p)
}

type walkableNode struct {
	parent *walkableField
	posInt
}

func (this walkableNode) gotoNode(p posInt) walkableNode {
	this.posInt = p
	return this
}

func (this walkableNode) passableNeighbors() (
	neighbors [4]posInt, passable [4]bool, count int,
) {
	p := this.parent
	if p == nil {
		return
	}
	neighbors = this.neighbors()
	count = 0
	for i, n := range neighbors {
		if p.isPassable(n) {
			passable[i] = true
			count++
		}
	}
	return
}
