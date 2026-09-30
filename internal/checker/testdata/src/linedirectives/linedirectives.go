package linedirectives

func long() int { // want "function long is 12 lines .limit 8."
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
//line imaginary.go:1
	return n
}

//line reset.go:1
func short() int {
//line inflated.go:1000
	return 1
}
