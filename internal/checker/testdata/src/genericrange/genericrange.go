package genericrange

type Maps interface{ ~map[int]int }
type NestedMaps interface{ Maps }

//mikoto:pure
func mapRange[T NestedMaps]() int {
	m := T{1: 1, 2: 2}
	for k := range m { // want "cannot range over a map, channel, or iterator"
		return k
	}
	return 0
}

type NamedMap map[int]int

//mikoto:pure
func namedMapRange[T NamedMap]() int {
	var m T
	for k := range m { // want "cannot range over a map, channel, or iterator"
		return k
	}
	return 0
}

//mikoto:pure
func channelRange[T interface{ ~chan int | ~<-chan int }]() {
	var ch T
	for range ch { // want "cannot range over a map, channel, or iterator"
	}
}

//mikoto:pure
func iteratorRange[T ~func(func(int) bool)]() {
	var iterator T
	for range iterator { // want "cannot range over a map, channel, or iterator"
	}
}

//mikoto:pure
func arrayRange[T ~[2]int]() int {
	a := T{1, 2}
	n := 0
	for _, v := range a {
		n += v
	}
	return n
}

//mikoto:pure
func integerRange[T ~int]() int {
	n := 0
	for i := range T(3) {
		n += int(i)
	}
	return n
}

// Intersections must exclude map terms that do not satisfy the full constraint.
type First interface{ ~[]int | ~map[int]int }
type Second interface{ ~[]int | ~map[string]int }
type SliceIntersection interface {
	First
	Second
}

//mikoto:pure
func sliceRange[T SliceIntersection]() int {
	a := T{1, 2}
	n := 0
	for _, v := range a {
		n += v
	}
	return n
}

type NamedSlice []int
type ExactSliceIntersection interface {
	~[]int | ~map[int]int
	NamedSlice
}

//mikoto:pure
func exactSliceRange[T ExactSliceIntersection]() int {
	a := T{1, 2}
	n := 0
	for _, v := range a {
		n += v
	}
	return n
}
