package pointeraccess

import "unsafe"

//mikoto:pure
func read(addr uintptr) int {
	p := (*[1]int)(unsafe.Pointer(addr))
	return p[0] // want "cannot dereference a pointer"
}

//mikoto:pure
func write(addr uintptr, n int) {
	p := (*[1]int)(unsafe.Pointer(addr))
	for p[0] = range n { // want "cannot assign through a reference" "cannot dereference a pointer"
	}
}

//mikoto:pure
func rangeValue(addr uintptr) int {
	p := (*[1]int)(unsafe.Pointer(addr))
	for _, v := range p { // want "cannot dereference a pointer"
		return v
	}
	return 0
}

//mikoto:pure
func slice(addr uintptr) int {
	p := (*[1]int)(unsafe.Pointer(addr))
	return p[:][0] // want "cannot dereference a pointer"
}

type Value struct{ N int }

//mikoto:pure
func field(addr uintptr) int {
	p := (*Value)(unsafe.Pointer(addr))
	return p.N // want "cannot dereference a pointer"
}

//mikoto:pure
func (v Value) Get() int { return v.N }

//mikoto:pure
func method(addr uintptr) int {
	p := (*Value)(unsafe.Pointer(addr))
	return p.Get() // want "cannot dereference a pointer"
}

//mikoto:pure
func genericRead[T ~*[1]int](addr uintptr) int {
	p := T(unsafe.Pointer(addr))
	return p[0] // want "cannot dereference a pointer"
}

type ArrayPointer[T any] *[1]T

// Instantiating a pointer type is not itself an indirect memory access.
//
//mikoto:pure
func pointerType() int {
	var p ArrayPointer[int]
	_ = p
	return 1
}

//mikoto:pure
func namedPointerRead(addr uintptr) int {
	p := ArrayPointer[int](unsafe.Pointer(addr))
	return p[0] // want "cannot dereference a pointer"
}

// Both range variables must follow the ordinary assignment rules.
//
//mikoto:pure
func rangeAssignments() {
	a := [2]int{}
	for a[0], a[1] = range [2]int{1, 2} { // want "cannot assign through a reference" "cannot assign through a reference"
	}
}

// Constant-length iteration does not read through an array pointer.
//
//mikoto:pure
func rangeIndices() int {
	var p *[2]int
	n := 0
	for i, _ := range p {
		n += i
	}
	return n
}

//mikoto:pure
func localValues() int {
	a := [2]int{1, 2}
	s := a[:]
	v := Value{N: 3}
	return a[0] + s[1] + v.N + v.Get()
}
