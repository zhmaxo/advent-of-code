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

	5: batchTestCase(`47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`, "143", "123"),

	6: batchTestCase(`....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`, "41", "6"),

	7: batchTestCase(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`, "3749", "11387"),

	8: batchTestCase(`............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`, "14", "34"),
}

func batchTestCase(input, expect1, expect2 string) testCase {
	return testCase{
		input:   input,
		expect1: expect1,
		expect2: expect2,
	}
}

func TestUtils(t *testing.T) {
	p1, p2 := posInt{2, 5}, posInt{3, 1}
	assert(t, p1.add(p2).eq(posInt{5, 6}),
		"unexpected result %v + %v = %v (expected %v)",
		p1, p2, p1.add(p2), posInt{5, 6})
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
		assert(t, answer1 == tc.expect1,
			"%v is incorrect d%vp1 test answer!", answer1, d)

		answer2, err := sol.SolvePt2()
		asserrt(t, err)
		assert(t, answer2 == tc.expect2,
			"%v is incorrect d%vp2 test answer!", answer2, d)
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
