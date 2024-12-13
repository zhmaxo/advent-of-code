package days

import (
	"bytes"
	"slices"
)

func startsWith(line, prefix []byte) bool {
	return len(line) >= len(prefix) &&
		bytes.Equal(line[:len(prefix)], prefix)
}

func endsWith(line, ending []byte) bool {
	return len(line) >= len(ending) &&
		bytes.Equal(line[len(line)-len(ending):], ending)
}

func substrIdx(line, substr []byte) (idx int) {
	idx = slices.Index(line, substr[0])
	if idx < 0 {
		return
	}
	endIdx := idx + len(substr)
	if endIdx >= len(line) {
		idx = -1
		return
	}
	if !bytes.Equal(line[idx:endIdx], substr) {
		idx = -1
	}
	return
}
