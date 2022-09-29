package bitset

import "golang.org/x/exp/constraints"

type Value interface {
	constraints.Unsigned |
		~uint | /* uint8 | */ uint16 | uint32 | uint64 |

		~rune | ~byte // a.k.a int32 / uint8 respectively

}

type Bitset[V Value] interface {
	Get(index V) bool

	Set(indices ...V) Bitset[V]
	Unset(indices ...V) Bitset[V]

	Clear()
}

func Copy[V Value](dst Bitset[V], src Inspect[V]) {
	dst.Clear()
	it := src.Iterate()
	i, vals := 0, make([]V, src.Pop())
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		vals[i] = n
		i++
	}
	dst.Set(vals...)
}

type Logical[V Value, T Bitset[V]] interface {
	// these functions are not expected to modify A or B
	And(b T) (aAndB T)
	Or(b T) (aOrB T)
}

type Size interface {
	// int and not uint for consistency's sake :(
	Len() int
	Cap() int
}

// experimental. see: https://github.com/golang/go/discussions/54245
type Iter[V Value] interface {
	Next() (V, bool)
}

type Inspect[V Value] interface {
	Size
	Pop() uint
	Iterate() Iter[V]
}
