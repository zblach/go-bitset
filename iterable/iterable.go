package iterable

import (
	"math"
	"sync"

	"github.com/zblach/go-bitset"
)

// The Iterable[V] interface allows for enumeration over a bitset. It's not available everywhere.
type Iterable[V bitset.Value] interface {
	Iterate() (it Iter[V], size uint)
}

// experimental. see: https://github.com/golang/go/discussions/54245. think python generators.
type Iter[V bitset.Value] interface {
	Next() (V, bool)
}

// Values collects all elements out of a bitset by instantiating and exhausting an iterator.
func Values[V bitset.Value](s Iterable[V]) (vals []V) {
	it, size := s.Iterate()
	vals = make([]V, 0, size)

	for v, ok := it.Next(); ok; v, ok = it.Next() {
		vals = append(vals, v)
	}

	return
}

// Copy copies the values from `src` into `dst`, clearing the contents of `dst` beforehand.
// These arguments need to track the same kind of value, but do not need to be the same implementation.
func Copy[V bitset.Value](dst bitset.Bitset[V], src Iterable[V]) {
	dst.Clear()
	dst.Set(Values(src)...)
}

// wrapper structs for heterogenous AND/OR operations
type groupIterator[V bitset.Value] struct {
	iters []Iterable[V]
}
type (
	andIterator[V bitset.Value] groupIterator[V]
	orIterator[V bitset.Value]  groupIterator[V]
)

type groupIter[V bitset.Value] struct {
	lock  *sync.RWMutex
	iters []*peekIter[V]
}
type (
	andIter[V bitset.Value] groupIter[V]
	orIter[V bitset.Value]  groupIter[V]
)

// an Iter[V], which remembers the last value read
type peekIter[V bitset.Value] struct {
	Iter[V]
	curr V
}

func (it *peekIter[V]) Next() (V, bool) {
	v, ok := it.Iter.Next()
	if ok {
		it.curr = v
	}
	return v, ok
}

var (
	_ Iter[byte]   = (*andIter[byte])(nil)
	_ Iter[uint16] = (*orIter[uint16])(nil)
)

// Next implements Iter for all 's'
func (and *andIter[V]) Next() (V, bool) {
	and.lock.Lock()
	defer and.lock.Unlock()

	var _v V
	if len(and.iters) == 0 {
		return _v, false
	}

	for {
		// find minimum value
		minIndex, maxIndex := 0, 0
		for i, it := range and.iters[1:] {
			if and.iters[minIndex].curr > it.curr {
				minIndex = i + 1
			} else if and.iters[maxIndex].curr < it.curr {
				maxIndex = i + 1
			}
		}

		if minIndex == maxIndex {
			// all iterators have this value. it's a return candidate.
			ret := and.iters[0].curr
			for _, it := range and.iters {
				// increment all iterators
				if _, ok := it.Next(); !ok {
					// if any of the iterators is exhausted, we're done. no more 'and's possible.
					and.iters = nil
					break
				}

			}

			return ret, true

		} else {
			// increment all iterators with the same 'low value'
			minVal := and.iters[minIndex].curr
			for _, it := range and.iters {
				if it.curr == minVal {
					if _, ok := it.Next(); !ok {
						// and if any of them are exhausted, we're done
						and.iters = nil
						return _v, false
					}
				}
			}
		}
	}
}

func (or *orIter[V]) Next() (V, bool) {
	or.lock.Lock()
	defer or.lock.Unlock()

	var _v V
	if len(or.iters) == 0 {
		return _v, false
	}

	// find the lowest value
	minIndex := 0
	for i := 1; i < len(or.iters); i++ {
		if or.iters[minIndex].curr > or.iters[i].curr {
			minIndex = i
		}
	}

	ret := or.iters[minIndex].curr
	// increment all iterators with that same value to avoid dupes
	nextIters := make([]*peekIter[V], 0, len(or.iters))
	for _, it := range or.iters {
		if it.curr != ret {
			nextIters = append(nextIters, it)
			continue
		}
		if _, ok := it.Next(); ok {
			// ignore exhausted iterators
			nextIters = append(nextIters, it)
		}
	}
	or.iters = nextIters

	return ret, true
}

func And[V bitset.Value](s1, s2 Iterable[V], s ...Iterable[V]) Iterable[V] {
	return andIterator[V]{iters: append(s[:], s1, s2)}
}

func Or[V bitset.Value](s1, s2 Iterable[V], s ...Iterable[V]) Iterable[V] {
	return orIterator[V]{iters: append(s[:], s1, s2)}
}

// Iterate implements Iterable
func (ai andIterator[V]) Iterate() (Iter[V], uint) {
	gi, min, _ := newIter(ai.iters...)
	if min == 0 {
		// then one of the iterators is empty. this means no values at all.
		gi.iters = nil
	}
	it := andIter[V](gi)
	return &it, min
}

func (oi orIterator[V]) Iterate() (Iter[V], uint) {
	gi, _, max := newIter(oi.iters...)
	it := orIter[V](gi)
	return &it, max

}

var _ Iterable[rune] = (*andIterator[rune])(nil)

// newIter is an internal utility function for iterator composition. the return values are
// gi - groupIter -- a 'Next()'-able object to query the next value from
// min - the minimum size of all overlapping sets. smallest of any individual iterator
// max - the maximum size of all overlapping sets. sum of all iterator lengths.
func newIter[V bitset.Value](s ...Iterable[V]) (gi groupIter[V], min, max uint) {
	gi = groupIter[V]{
		lock:  &sync.RWMutex{},
		iters: make([]*peekIter[V], 0, len(s)),
	}

	if len(s) == 0 {
		return
	}

	min, max = math.MaxUint, 0

	for _, iter := range s {
		it, siz := iter.Iterate()

		if siz < min {
			min = siz
		}
		max += siz

		if v, ok := it.Next(); ok { // zero size iterators are not included
			gi.iters = append(gi.iters, &peekIter[V]{
				Iter: it,
				curr: v,
			})
		}
	}

	return

}
