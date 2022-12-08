package sparse_set

import "golang.org/x/exp/constraints"

// sparseRange is an internal structure for storing range details.
type Range[V constraints.Integer] struct {
	Start, End V // inclusive
}

func NewRange[V constraints.Integer](start, end V) *Range[V] {
	return &Range[V]{start, end}
}

func (s *Range[V]) Contains(val V) bool {
	return val >= s.Start && val <= s.End
}

func (s *Range[V]) Valid() bool {
	return s.Start <= s.End
}
