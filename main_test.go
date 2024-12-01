package main

import (
	"strings"
	"testing"
	"zhmaxo/advent-of-code-2024/days"
)

func TestSome(t *testing.T) {
	reader := strings.NewReader("test")
	scanner := days.ProcessReader(reader)
	hasLine := true
	for hasLine {
		line, isPrefix, err := scanner.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		hasLine = line != nil
		println("P:", isPrefix, " ", string(line))
	}
}
