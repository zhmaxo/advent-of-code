package days

import (
	"strings"
	"testing"
)

type testCase struct {
	input   string
	expect1 string
	expect2 string
}

var testCases = map[uint8]testCase{
	1: batchTestCase(`3   4
  4   3
  2   5
  1   3
  3   9
  3   3`,
		"11", "31"),

	2: batchTestCase(`7 6 4 2 1
  1 2 7 8 9
  9 7 6 2 1
  1 3 2 4 5
  8 6 4 4 1
  1 3 6 7 9`, "2", "4"),

	3: batchTestCase(`xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`,
		"161", "48"),

	4: batchTestCase(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`,
		"18", "9"),
}

func batchTestCase(input, expect1, expect2 string) testCase {
	return testCase{
		input:   input,
		expect1: expect1,
		expect2: expect2,
	}
}

func TestDays(t *testing.T) {
	for d, tc := range testCases {
		sol, ok := DaySolutions[d]
		if !ok {
			t.Logf("there's registered day %v but no solution registered to check", d)
			continue
		}
		reader := strings.NewReader(tc.input)
		err := sol.ReadData(reader)
		asserrt(t, err)

		answer1, err := sol.SolvePt1()
		asserrt(t, err)
		assert(t, answer1 == tc.expect1, "%v is incorrect part1 test answer!", answer1)

		answer2, err := sol.SolvePt2()
		asserrt(t, err)
		assert(t, answer2 == tc.expect2, "%v is incorrect part2 test answer!", answer2)
	}
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
