package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"zhmaxo/advent-of-code-2024/days"
)

func main() {
	if err := aoc2024Main(); err != nil {
		fmt.Printf("err: %v\n", err)
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func aoc2024Main() error {
	day, file, err := getCliArgs()
	if err != nil {
		return err
	}

	reader, err := getCurrentInputReader(day, file)
	if err != nil {
		return err
	}
	defer reader.Close()

	solution, ok := days.DaySolutions[uint8(day)]
	if !ok {
		err = fmt.Errorf("no day with number %v", day)
		return err
	}

	err = solution.ReadData(reader)
	if err != nil {
		return err
	}

	answer1, err := solution.SolvePt1()
	if err != nil {
		return err
	}

	fmt.Printf("pt1 answer: %v\n", answer1)

	answer2, err := solution.SolvePt2()
	if err != nil {
		return err
	}

	fmt.Printf("pt2 answer: %v\n", answer2)
	return nil
}

func getCliArgs() (day uint, file string, err error) {
	dayNumber := flag.Uint("day", 0, "specify day to solve")
	filename := flag.String("file", "", "specify input filename here")
	flag.Parse()

	switch {
	case dayNumber == nil, filename == nil:
		err = fmt.Errorf("unexpected nil flag values")
		return
	}

	file = *filename
	day = *dayNumber

	const maxDayNumber = uint(25)
	switch {
	case day == 0:
		err = fmt.Errorf("you should specify non-zero day number")
	case day > maxDayNumber:
		err = fmt.Errorf("day number cannot be more than %v (get %v)", maxDayNumber, dayNumber)
	}
	return
}

func getCurrentInputReader(day uint, file string) (io.ReadCloser, error) {
	if file != "" {
		return os.Open(file)
	}
	// TODO: http get input via API
	return nil, fmt.Errorf("cant get input reader for day %v", day)
}
