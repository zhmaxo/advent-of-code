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

	9: batchTestCase(`2333133121414131402`, "1928", "2858"),

	11: batchTestCase(`125 17`, "55312", ""),

	12: batchTestCase(`RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`, "1930", ""),

	13: batchTestCase(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`, "480", ""),

	14: batchTestCase(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`, "12", ""),
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

		if tester, ok := sol.(Tester); ok {
			tester.ModifyForTest()
		}

		answer1, err := sol.SolvePt1()
		asserrt(t, err)
		assert(t, answer1 == tc.expect1,
			"%v is incorrect d%vp1 test answer!", answer1, d)

		answer2, err := sol.SolvePt2()
		asserrt(t, err)
		assert(t, tc.expect2 == "" || answer2 == tc.expect2,
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
