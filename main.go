package main

import (
	"flag"
	"fmt"
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
	dayNumber := flag.Uint("day", 0, "specify day to solve")
	filename := flag.String("file", "", "specify input filename here")
	flag.Parse()

	switch {
	case dayNumber == nil:
		fallthrough
	case filename == nil:
		return fmt.Errorf("unexpected nil flag values")
	}

	const maxDayNumber = uint(25)
	switch day := *dayNumber; {
	case day == 0:
		return fmt.Errorf("you should specify non-zero day number")
	case day > maxDayNumber:
		return fmt.Errorf("day number cannot be more than %v (get %v)", maxDayNumber, dayNumber)
	}

	if *filename != "" {
		reader, err := days.ReadFile(*filename)
		if err != nil {
			return err
		}

		answer, err := days.DaySolutions[uint8(*dayNumber)](reader)
		if err != nil {
			return err
		}

		fmt.Printf("answer: %v", answer)
	} else {
		println("no filename specified")
	}
	return nil
}
