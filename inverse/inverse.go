package inverse

import "github.com/zblach/bitset"

type Not struct {
	s bitset.Bitset
}

func New(s bitset.Bitset) bitset.Bitset {
	switch n := s.(type) {
	case *Not:
		return n.s
	}
	return &Not{
		s: s,
	}
}

func (n *Not) Clear() {
	n.s.Clear()
}

// Get implements bitset.Bitset
func (n *Not) Get(index uint) bool {
	return !n.s.Get(index)
}

// Set implements bitset.Bitset
func (n *Not) Set(indices ...uint) bitset.Bitset {
	n.s.Unset(indices...)
	return n
}

// Unset implements bitset.Bitset
func (n *Not) Unset(indices ...uint) bitset.Bitset {
	n.s.Set(indices...)
	return n
}

var _ bitset.Bitset = (*Not)(nil)
