package example

import "time"

var global = 2

//mikoto:pure
func square(n int) int { return n * n }

//mikoto:pure
func good(n int) int { return square(n) + len("a") }

//mikoto:pure
func bad(n int, p *int) int { // want "requires value-only parameters"
	global = n                    // want "cannot access package variable"
	*p = n                        // want "cannot assign through a reference" "cannot dereference a pointer"
	return int(time.Now().Unix()) // want "can only call another marked pure function" "can only call another marked pure function"
}
