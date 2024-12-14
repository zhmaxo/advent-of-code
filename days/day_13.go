package days

import (
	"regexp"
	"strconv"
)

func init() {
	DaySolutions[13] = &day13Solution{}
}

type day13Solution struct {
	equations []clawEquation
}

type clawEquation struct {
	A, B  posInt
	prize posInt
}

func (s *day13Solution) HasData() bool {
	return true
}

func (s *day13Solution) ReadData(reader ioReader) (err error) {
	const (
		startSize    = 512
		linesInBlock = 4
		regStr       = `.*X.([0-9]*), Y.([0-9]*)`
	)
	s.equations = make([]clawEquation, 0, startSize)
	r := regexp.MustCompile(regStr)
	sc := ProcessReader(reader)
	defer func() { err = RefineError(err) }()
	lineNum := 0
	var ce clawEquation
	var line []byte
	for err == nil {
		line, _, err = sc.ReadLine()
		if err != nil {
			break
		}

		ln := (lineNum + linesInBlock) % linesInBlock
		switch ln {
		case linesInBlock - 1:
		default:
			match := r.FindSubmatch(line)
			if len(match) < 3 {
				continue
			}
			var p posInt
			p.x, err = strconv.Atoi(string(match[1]))
			if err != nil {
				break
			}
			p.y, err = strconv.Atoi(string(match[2]))
			if err != nil {
				break
			}
			switch ln {
			case 0: // parse Button A
				ce.A = p
			case 1: // parse Button B
				ce.B = p
			case 2: // parse Prize
				ce.prize = p
				s.equations = append(s.equations, ce)
			}
		}
		lineNum++
	}
	return
}

func (s *day13Solution) SolvePt1() (answer string, err error) {
	const (
		priceA = 3
		priceB = 1
	)
	result := 0
	for _, e := range s.equations {
		// Xa*A+Xb*B = X <=> Xa*A = X - Xb*B <=> A = (X - Xb*B)/Xa
		// Ya*A+Yb*B = Y <=> Yb*B = Y - Ya*A <=> B = (Y - Ya*A)/Yb
		// A = (Yb*X - Xb*Y)/(Xa*Yb - Xb*Ya)
		// B = (Y - Ya*A)/Yb
		Xa, Xb, X := e.A.x, e.B.x, e.prize.x
		Ya, Yb, Y := e.A.y, e.B.y, e.prize.y
		A := (Yb*X - Xb*Y) / (Xa*Yb - Xb*Ya)
		B := (Y - Ya*A) / Yb

		checkPrizeX := A*Xa + B*Xb
		checkPrizeY := A*Ya + B*Yb
		check := checkPrizeX == X && checkPrizeY == Y

		if !check {
			// pf("equation %v cant be solved with integers [A=%v, B=%v is wrong]: (%v, %v != %v, %v)\n",
			// 	e, A, B, checkPrizeX, checkPrizeY, e.prize.x, e.prize.y)
			continue
		}
		// pf("result: A=%v, B=%v\n", A, B)
		result += A*priceA + B*priceB
	}
	answer = Stringify(result)
	return
}

func (s *day13Solution) SolvePt2() (answer string, err error) {
	const (
		priceA     = 3
		priceB     = 1
		correction = 10000000000000
	)
	result := float64(0)
	for i, e := range s.equations {
		// Xa*A+Xb*B = X <=> Xa*A = X - Xb*B <=> A = (X - Xb*B)/Xa
		// Ya*A+Yb*B = Y <=> Yb*B = Y - Ya*A <=> B = (Y - Ya*A)/Yb
		// A = (Yb*X - Xb*Y)/(Xa*Yb - Xb*Ya)
		// B = (Y - Ya*A)/Yb
		Xa, Xb, X := float64(e.A.x), float64(e.B.x), float64(e.prize.x)
		Ya, Yb, Y := float64(e.A.y), float64(e.B.y), float64(e.prize.y)
		X += correction
		Y += correction
		A := (Yb*X - Xb*Y) / (Xa*Yb - Xb*Ya)
		B := (Y - Ya*A) / Yb

		// checkPrizeX := A*Xa + B*Xb
		// checkPrizeY := A*Ya + B*Yb
		c1 := A-float64(uint64(A)) == 0
		c2 := B-float64(uint64(B)) == 0
		check := c1 && c2

		pf("%v: %v (%v and %v) A=%v (%v), B=%v (%v)\n",
			i, check, c1, c2, uint64(A), A-float64(uint64(A)), uint64(B), B-float64(uint64(B)))
		if !check {
			// pf("equation %v cant be solved with integers [A=%v, B=%v is wrong]: (%v, %v != %v, %v)\n",
			// 	e, A, B, checkPrizeX, checkPrizeY, X, Y)
			continue
		}
		// pf("result: A=%v, B=%v\n", A, B)
		prevRes := result
		result += A*priceA + B*priceB
		if result < prevRes {
			pf("oops overflow from %v + %v to %v", prevRes, A*priceA+B*priceB, result)
		}
	}
	answer = Stringify(uint64(result))
	return
}

// func (s *day13Solution) ReadData_DEP(reader ioReader) (err error) {
// 	const (
// 		l1start  = "Button A: X+"
// 		l2start  = "Button B: X+"
// 		sep1and2 = ", Y+"
// 		l3start  = "Prize: X="
// 		sep3     = ", Y="
// 	)
// 	type lineParseData struct {
// 		start, sep []byte
// 	}
//
// 	splitsByEmptyLines := true
// 	parseRules := []lineParseData{
// 		{[]byte(l1start), []byte(sep1and2)},
// 		{[]byte(l2start), []byte(sep1and2)},
// 		{[]byte(l3start), []byte(sep3)},
// 	}
//
// 	func(rules, byte) (res posInt, err error) {
// 		return
// 	}
//
// 	sc := ProcessReader(reader)
// 	defer func() { err = RefineError(err) }()
//
// 	parseErr := fmt.Errorf("can't parse line")
//
// 	lineNum := 0
// 	linesInBlock := len(parseRules)
// 	if splitsByEmptyLines {
// 		linesInBlock++
// 	}
// 	var line []byte
// 	for err == nil {
// 		line, _, err = sc.ReadLine()
// 		if err != nil {
// 			break
// 		}
//
// 		ruleIdx := lineNum % linesInBlock
// 		if ruleIdx >= len(parseRules) {
// 			// empty line
// 			continue
// 		}
//
// 		rule := parseRules[ruleIdx]
// 		var parseRes posInt
// 		if !startsWith(line, rule.start) {
// 			err = parseErr
// 			break
// 		}
// 	}
// 	return
// }
