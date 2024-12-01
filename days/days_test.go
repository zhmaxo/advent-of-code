package days

import "testing"

func TestDay1(t *testing.T) {
	testInput := `3   4
  4   3
  2   5
  1   3
  3   9
  3   3`

	reader := StrToReader(testInput)
	answer, err := solve_day1(reader)
	if err != nil {
		t.Fatal(err)
	}
	if answer != "11" {
		t.Fatalf("%v is incorrect test answer!", answer)
	}
	t.Log("as expected")
}
