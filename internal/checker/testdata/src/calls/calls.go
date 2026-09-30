package calls

//mikoto:pure
func constant[T any]() int { return 1 }

//mikoto:pure
func pair[T, U any]() int { return 2 }

type Box[T any] struct{ value int }

//mikoto:pure
func (b Box[T]) Value() int { return b.value }

//mikoto:pure
func genericCalls(n int) int {
	return constant[int]() + pair[int, string]() + Box[int]{value: n}.Value()
}

//mikoto:pure
func square(n int) int { return n * n }

//mikoto:pure
func parenthesized(n int) int {
	(n) = 4
	(n)++
	(n)--
	return (square)(n) + (len)("abc") + (constant[int])() + (pair[int, string])()
}

//mikoto:pure
func parenthesizedRange(n int) int {
	for (n) = range 3 {
	}
	return n
}

func unmarked[T any]() int { return 1 }

//mikoto:pure
func rejectUnmarked() int {
	return (unmarked[int])() // want "can only call another marked pure function"
}

type ImpureBox[T any] struct{ value int }

func (b ImpureBox[T]) Value() int { return b.value }

//mikoto:pure
func rejectUnmarkedMethod() int {
	return ImpureBox[int]{}.Value() // want "can only call another marked pure function"
}

// Function values must not be authorized by looking through ordinary indexing.
//
//mikoto:pure
func rejectIndexedFunction() int {
	f := [1]func(int) int{square}
	return f[0](2) // want "can only call another marked pure function"
}

// Parentheses must not authorize a forbidden builtin.
//
//mikoto:pure
func rejectBuiltin() {
	(println)(1) // want "cannot call builtin println"
}
