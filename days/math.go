package days

func pow(n int, p uint8) int {
	res := 1
	for range p {
		res *= n
	}
	return res
}
