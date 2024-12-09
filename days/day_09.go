package days

import "fmt"

func init() {
	DaySolutions[9] = &day9Solution{}
}

type day9Solution struct {
	files  []memBlock
	spaces []memBlock
}

// represents sequence of memory positions
// uses for store spaces info and ordered files
type memBlock struct {
	start, length uint
}

// represents any memory sequence that stores moved file chunk
// NOTE: there's no info about wholeness and whatever else
type fileBlock struct {
	memBlock
	fileID uint
}

func (mb memBlock) checksum(fileID uint) uint {
	var result uint
	for i := mb.start; i < mb.start+mb.length; i++ {
		result += i * fileID
		fmt.Printf("%v * %v = %v\n", i, fileID, i*fileID)
	}
	return uint(result)
}

func (s *day9Solution) HasData() bool {
	return true
}

func (s *day9Solution) ReadData(reader ioReader) (err error) {
	s.files = make([]memBlock, 0, 2048)
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()
	var b byte
	currentPos := uint(0)
	var isFile bool
	for err == nil {
		isFile = !isFile
		b, err = sc.ReadByte()
		if err != nil {
			break
		}
		blockLen := uint(b - '0')
		if blockLen == 0 {
			continue
		}
		mb := memBlock{uint(currentPos), blockLen}
		if isFile {
			fmt.Printf("add file(%v) p%v, len%v\n", len(s.files), mb.start, mb.length)
			s.files = append(s.files, mb)
		} else {
			fmt.Printf("add space p%v, len%v\n", mb.start, mb.length)
			s.spaces = append(s.spaces, mb)
		}
		currentPos += blockLen
	}
	fmt.Printf("total files: %v\n", len(s.files))
	return
}

func (s *day9Solution) SolvePt1() (answer string, err error) {
	checkSum := uint64(0)
	// idOffset := 0
	spaceOffset := 0

	// TODO: copy before second part
	files, spaces := s.files, s.spaces
	movedFiles := make([]fileBlock, 0, len(s.files))
	// move each file from last to first
	for i := len(files) - 1; i > 0 && spaceOffset < len(spaces); {
		if spaces[spaceOffset].start > files[i].start {
			break
		}
		moveFileID := uint(i)
		fileid := &i
		var movedBlock memBlock
		movedBlock, _ = memmove(files, spaces, fileid, &spaceOffset)
		if movedBlock.length > 0 {
			movedFiles = append(movedFiles, fileBlock{movedBlock, moveFileID})
		}
	}

	pf("begin counting regular files checksum:\n")
	for i, mb := range files {
		if mb.length == 0 {
			break
		}

		fileID := uint(i)
		checkSum += uint64(mb.checksum(fileID))
	}

	pf("begin counting moved files checksum:\n")
	for _, fb := range movedFiles {
		checkSum += uint64(fb.checksum(fb.fileID))
	}
	// pf("check file %v; move target %v\n", blockID, curMvFile)
	answer = Stringify(checkSum)
	return
}

func (s *day9Solution) SolvePt2() (answer string, err error) {
	return
}

func memmove(files, spaces []memBlock, fileid, spaceid *int) (mvResult memBlock, spcused int) {
	curspcid := *spaceid
	curmvfid := *fileid
	if curspcid < len(spaces) {
		spc := &spaces[curspcid]
		mvf := &files[curmvfid]
		pf("move mem %v to space %v\n", mvf, spc)

		// lenDiff := spc.length - mvf.length
		if spc.length > mvf.length {
			mvResult = memBlock{spc.start, mvf.length}
			spc.length -= mvf.length
			spc.start += mvf.length
			mvf.length = 0
			curmvfid--
		} else {
			spcused = 1
			mvResult = *spc
			mvf.length -= spc.length
			spc.length = 0
			curspcid++
			if spc.length == mvf.length {
				curmvfid--
			}
		}
		*fileid = curmvfid
		*spaceid = curspcid
	}
	return
}
