package length

func long() int { // want "function long is 5 lines"
	n := 1
	n++
	return n
}
