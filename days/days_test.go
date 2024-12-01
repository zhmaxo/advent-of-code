package days

import (
	"strings"
	"testing"
)

func TestDay1(t *testing.T) {
	testInput := `3   4
  4   3
  2   5
  1   3
  3   9
  3   3`

	reader := strings.NewReader(testInput)
	// first part
	s := day1Solution{}
	err := s.ReadData(reader)
	asserrt(t, err)

	answer1, err := s.SolvePt1()
	asserrt(t, err)
	assert(t, answer1 == "11", "%v is incorrect part1 test answer!", answer1)

	answer2, err := s.SolvePt2()
	asserrt(t, err)
	assert(t, answer2 == "31", "%v is incorrect part2 test answer", answer2)
}

func assert(t *testing.T, condition bool, failMsg string, args ...any) {
	if !condition {
		t.Fatalf(failMsg, args...)
	}
}

func asserrt(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func asserrtf(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%v:  %v", msg, err)
	}
}
