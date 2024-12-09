package days

import (
	"fmt"
	"sort"
)

func init() {
	DaySolutions[9] = &day9Solution{}
}

type day9Solution struct {
	spaceCache map[byte][]int

	files  []memBlock
	spaces []memBlock
}

func (s *day9Solution) HasData() bool {
	return true
}

func (s *day9Solution) ReadData(reader ioReader) (err error) {
	const startSize = 10000
	s.files = make([]memBlock, 0, startSize)
	s.spaces = make([]memBlock, 0, startSize)
	s.spaceCache = prepareSpacesCache(1024)
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
			l := byte(mb.length)
			s.spaceCache[l] = append(s.spaceCache[l], len(s.spaces))
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
	files, spaces := s.cpyData()
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

	// pf("begin counting regular files checksum:\n")
	for i, mb := range files {
		if mb.length == 0 {
			break
		}

		fileID := uint(i)
		checkSum += uint64(mb.checksum(fileID))
	}

	// pf("begin counting moved files checksum:\n")
	for _, fb := range movedFiles {
		checkSum += uint64(fb.checksum(fb.fileID))
	}
	// pf("check file %v; move target %v\n", blockID, curMvFile)
	answer = Stringify(checkSum)
	return
}

func (s *day9Solution) SolvePt2() (answer string, err error) {
	pf("part 2\n")
	checkSum := uint64(0)
	spaceOffset := 0
	files, spaces := s.files, s.spaces
	newSpaces := prepareSpacesCache(64)
	movedFiles := make([]fileBlock, len(files))
	for i := len(files) - 1; i > 0 && spaceOffset < len(spaces); {
		mb := &files[i]
		if spaces[spaceOffset].start > mb.start {
			break
		}
		targetLen := byte(mb.length)
		for ; targetLen < 10; targetLen++ {
			var spcPos int
			cacheToUse := newSpaces
			spaceToUse := newSpaces[targetLen]
			// select leftmost free space
			if len(spaceToUse) > 0 {
				spcPos = spaceToUse[0]
				altSpace := s.spaceCache[targetLen]
				if len(altSpace) > 0 {
					// less index -> less start pos
					if spcPos > altSpace[0] {
						cacheToUse = s.spaceCache
						spaceToUse = altSpace
						spcPos = spaceToUse[0]
					}
				}
			} else {
				cacheToUse = s.spaceCache
				spaceToUse = s.spaceCache[targetLen]
				if len(spaceToUse) < 1 {
					continue
				}
				spcPos = spaceToUse[0]
			}
			if spaces[spcPos].start > files[i].start {
				continue
			}
			cacheToUse[targetLen] = spaceToUse[1:]
			mvd, spcLeft := memmove(files, spaces, &i, &spcPos)
			movedFiles = append(movedFiles, fileBlock{mvd, uint(i + 1)})
			if spcLeft > 0 {
				spcs := append(newSpaces[spcLeft], spcPos)
				// should keep order to proper checks
				sort.Ints(spcs)
				newSpaces[spcLeft] = spcs
			}
			break
		}
		pf("not found space for mem block %v\n", files[i])
		i-- // not found space for that block
	}
	for id := 0; id < len(files); id++ {
		mb := files[id]
		if mb.length == 0 {
			continue
		}
		checkSum += uint64(mb.checksum(uint(id)))
	}

	for id := range movedFiles {
		fb := movedFiles[id]
		checkSum += uint64(fb.checksum(fb.fileID))
	}
	answer = Stringify(checkSum)
	return
}

func (s *day9Solution) cpyData() (files, spaces []memBlock) {
	files = make([]memBlock, len(s.files))
	copy(files, s.files)
	spaces = make([]memBlock, len(s.spaces))
	copy(spaces, s.spaces)
	return
}

func prepareSpacesCache(capacity int) (cache map[byte][]int) {
	const maxLength = 9
	cache = make(map[byte][]int, maxLength)
	for i := byte(1); i <= maxLength; i++ {
		cache[i] = make([]int, 0, capacity)
	}
	return
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

func memmove(files, spaces []memBlock, fileid, spaceid *int) (mvResult memBlock, spaceRemains byte) {
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
			spaceRemains = byte(spc.length)
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
