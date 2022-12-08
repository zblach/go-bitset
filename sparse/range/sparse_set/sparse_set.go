package sparse_set

import "golang.org/x/exp/constraints"

// sparseSet is an ordered series of disjoint sparseRange items. overlapping ranges are coalesced.
type Set[V constraints.Integer] []Range[V]

func (s *Set[V]) Insert(val V) bool {
	for i, r := range *s {
		switch {
		case r.Contains(val):
			// val exists already
			return false

		case r.Start-1 == val:
			// we're not contained by the previous set, but this one can be extended forward
			(*s)[i].Start -= 1
			return true

		case r.Start > val:
			// doesn't exist in the set up to here. insert it at this index.
			*s = append(*s, Range[V]{})
			for j := len(*s) - 1; j > i; j-- {
				(*s)[j] = (*s)[j-1]
			}
			(*s)[i] = Range[V]{val, val}
			return true

		case r.End+1 == val:
			// extend this element, and check if it collides into the next element (if it exists)
			(*s)[i].End += 1
			if i+1 == len(*s) {
				// no next element
				return true
			}
			if (*s)[i+1].Start == val+1 {
				// merge [i], [i+1]
				(*s)[i].End = (*s)[i+1].End
				// remove [i+1]
				(*s) = append((*s)[:i+1], (*s)[i+2:]...)
			}
			return true
		}
	}
	// it's none of any of the above cases. add it to the end of the list
	*s = append(*s, Range[V]{val, val})
	return true
}

func (s *Set[V]) Remove(val V) bool {
	for i, r := range *s {
		switch {
		case r.Start > val:
			return false
			// element not present to be removed. no-op.

		case r.Start == val:
			// r can have start modified
			(*s)[i].Start += 1
			if (*s)[i].Start > r.End {
				// remove invalid element.
				if i == 0 {
					(*s) = (*s)[1:]
				} else {
					for j := i; j < len(*s)-1; j++ {
						(*s)[j] = (*s)[j+1]
					}
					(*s) = (*s)[:len(*s)-1]
				}
			}

		case r.End == val:
			// r can have end modified
			(*s)[i].End -= 1
			// if the element was invalid here, then it was one element long, and would have been caught by the previous case.

		case r.Contains(val):
			// split r into two parts, excluding val
			(*s) = append(*s, Range[V]{})
			for j := len(*s) - 1; j > i; j-- {
				(*s)[j] = (*s)[j-1]
			}
			(*s)[i].End = val - 1

			(*s)[i+1].Start = val + 1
			(*s)[i+1].End = r.End
		}
	}
	return true
}
