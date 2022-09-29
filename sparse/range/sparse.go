package rangeset

import (
	"sync"

	"github.com/zblach/go-bitset"
	"golang.org/x/exp/constraints"
)

type Bitset[V bitset.Value] struct {
	lock *sync.RWMutex

	sets sparseSet[V]
	pop  uint
}

// And implements bitset.Logical
func (a *Bitset[V]) And(b *Bitset[V]) (aAndB *Bitset[V]) {
	aAndB = New[V]()

	it_a, it_b := a.Iterate(), b.Iterate()

	val_a, next_a := it_a.Next()
	val_b, next_b := it_b.Next()

	for {
		switch {
		case !next_a, !next_b:
			return
		case val_a == val_b:
			aAndB.Set(val_a)
			val_a, next_a = it_a.Next()
			val_b, next_b = it_b.Next()
		case val_a < val_b:
			val_a, next_a = it_a.Next()
		case val_a > val_b:
			val_b, next_b = it_b.Next()
		}
	}
}

// Or implements bitset.Logical
func (a *Bitset[V]) Or(b *Bitset[V]) (aOrB *Bitset[V]) {
	var short, long *Bitset[V]
	if a.pop > b.pop {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aOrB = &Bitset[V]{
		lock: &sync.RWMutex{},
		sets: make(sparseSet[V], len(long.sets)),
		pop:  long.pop,
	}
	copy(aOrB.sets, long.sets)

	it_s := short.Iterate()
	for v, ok := it_s.Next(); ok; v, ok = it_s.Next() {
		aOrB.Set(v)
	}

	return
}

// Cap implements bitset.Inspect
func (s *Bitset[V]) Cap() int {
	return int(s.pop)
}

// Len implements bitset.Inspect
func (s *Bitset[V]) Len() int {
	return int(s.pop)
}

// Pop implements bitset.Inspect
func (s *Bitset[V]) Pop() uint {
	return s.pop
}

func New[V bitset.Value]() *Bitset[V] {
	return &Bitset[V]{
		lock: &sync.RWMutex{},
		sets: []sparseRange[V]{},
	}
}

func (s *Bitset[V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.sets = make([]sparseRange[V], 0)
}

// Get implements bitset.Bitset
func (s *Bitset[V]) Get(index V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, st := range s.sets {
		if st.contains(index) {
			return true
		}
	}

	return false
}

// Set implements bitset.Bitset
func (s *Bitset[V]) Set(indices ...V) bitset.Bitset[V] {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if s.sets.insert(index) {
			s.pop++
		}
	}
	return s
}

// Unset implements bitset.Bitset
func (s *Bitset[V]) Unset(indices ...V) bitset.Bitset[V] {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if s.sets.remove(index) {
			s.pop--
		}
	}
	return s
}

var (
	_ bitset.Bitset[byte]                 = (*Bitset[byte])(nil)
	_ bitset.Logical[rune, *Bitset[rune]] = (*Bitset[rune])(nil) // TBD
	_ bitset.Inspect[uint]                = (*Bitset[uint])(nil)
)

type sparseSet[V constraints.Integer] []sparseRange[V]

func (s *sparseSet[V]) insert(val V) bool {
	for i, r := range *s {
		switch {
		case r.contains(val):
			// val exists already
			return false
		case r.start-1 == val:
			// we're not contained by the previous set, but this one can be extended forward
			(*s)[i].start -= 1
			return true

		case r.start > val:
			// doesn't exist in the set up to here. insert it at this index.
			*s = append(*s, sparseRange[V]{})
			for j := len(*s) - 1; j > i; j-- {
				(*s)[j] = (*s)[j-1]
			}
			(*s)[i] = sparseRange[V]{val, val}
			return true

		case r.end+1 == val:
			// extend this element, and check if it collides into the next element (if it exists)
			(*s)[i].end += 1
			if i+1 == len(*s) {
				return true
			}
			if (*s)[i+1].start == val+1 {
				(*s)[i].end = (*s)[i+1].end
				(*s) = append((*s)[:i+1], (*s)[i+2:]...)
			}
			return true
		}
	}
	// it's none of any of the above cases. add it to the end of the list
	*s = append(*s, sparseRange[V]{val, val})
	return true
}

func (s *sparseSet[V]) remove(val V) bool {
	for i, r := range *s {
		switch {
		case r.start > val:
			return false
			// element not present to be removed. no-op.

		case r.start == val:
			// r can have start modified
			(*s)[i].start += 1
			if (*s)[i].start > r.end {
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

		case r.end == val:
			// r can have end modified
			(*s)[i].end -= 1
			// if the element was invalid here, then it was one element long, and would have been caught by the previous case.

		case r.contains(val):
			// split r into two parts, excluding val
			(*s) = append(*s, sparseRange[V]{})
			for j := len(*s) - 1; j > i; j-- {
				(*s)[j] = (*s)[j-1]
			}
			(*s)[i].end = val - 1

			(*s)[i+1].start = val + 1
			(*s)[i+1].end = r.end
		}
	}
	return true
}

type sparseRange[V constraints.Integer] struct {
	start, end V // inclusive
}

func (s *sparseRange[V]) contains(val V) bool {
	return val >= s.start && val <= s.end
}
