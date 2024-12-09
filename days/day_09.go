package days

import "fmt"

func init() {
	DaySolutions[9] = &day9Solution{}
}

type day9Solution struct {
	files  []memBlock
	spaces []memBlock
}

type fileBlock struct {
	memBlock
	fileID uint
}

type memBlock struct {
	start, length uint
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
	curSpc := 0
	curMvFile := len(s.files) - 1
	for blockID := range s.files {
		mb := &s.files[blockID]
		pf("check file %v; move target %v\n", blockID, curMvFile)
		// seems like sane conversion
		fileID := uint(blockID)
		if blockID >= curMvFile {
			for len(s.spaces) > curSpc && s.spaces[curSpc].start < mb.start && blockID == curMvFile {
				mvRes := memMov(s.files, s.spaces, &curMvFile, &curSpc)
				if mvRes.length > 0 {
					checkSum += uint64(mvRes.checksum(fileID))
				}
				pf("%v > %v && %v < %v && %v == %v\n",
					len(s.spaces), curSpc, s.spaces[curSpc].start, fileID, blockID, curMvFile)
			}
			// pf("%v > %v && %v < %v && %v == %v\n",
			// 	len(s.spaces), curSpc, s.spaces[curSpc].start, fileID, blockID, curMvFile)
			checkSum += uint64(s.files[blockID].checksum(fileID))
			// collapse
			break
		} else {
			// keep count
			checkSum += uint64(mb.checksum(fileID))
			// cache mv file id
			fileID = uint(curMvFile)
			mvRes := memMov(s.files, s.spaces, &curMvFile, &curSpc)
			if mvRes.length > 0 {
				checkSum += uint64(mvRes.checksum(fileID))
			}
		}
	}
	answer = Stringify(checkSum)
	return
}

func (s *day9Solution) SolvePt2() (answer string, err error) {
	return
}

func memMov(files, spaces []memBlock, fileid, spaceid *int) (mvResult memBlock) {
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
