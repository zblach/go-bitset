package bitset

import (
	"golang.org/x/exp/constraints"
)

type Value interface {
	constraints.Unsigned |
		~rune // this is int32, but no runes are actually negative.
}

// Basic functionality of a bitset implementation
type Bitset[V Value] interface {
	Get(index V) bool

	Set(indices ...V)
	Unset(indices ...V)

	Clear()
}

// Binary are chainable binary operations over bitsets of the same type.
// these functions are not expected to modify A or B, and S is also meant to conform to Binary[V, S Bitset[V]]
type Binary[V Value, S Bitset[V]] interface {
	And(b S) (aAndB S) // self
	Or(b S) (aOrB S)   // self
}

// Size-related inspection referring to the underlying bitset data storage
type Size interface {
	// int and not uint for consistency's sake :(
	Len() int
	Cap() int
}

// Population count & size. TODO: move 'Pop' to Bitset proper?
type Inspect[V Value] interface {
	Size
	Pop() uint
}
