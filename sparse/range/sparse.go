package rangeset

import (
	"sync"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/mixin/logical"
	"github.com/zblach/go-bitset/sparse/range/sparse_set"
)

type Bitset[V bitset.Value] struct {
	logical.IterableMixin[V]

	lock *sync.RWMutex

	sets sparse_set.Set[V]
	pop  uint
}

// And implements bitset.Logical
func (a *Bitset[V]) And(b *Bitset[V]) (aAndB *Bitset[V]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	aAndB = New[V]()

	it_a, _ := a.Iterate()
	it_b, _ := b.Iterate()

	val_a, next_a := it_a.Next()
	val_b, next_b := it_b.Next()

	for next_a && next_b {
		switch {
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

	return
}

// Or implements bitset.Logical
func (a *Bitset[V]) Or(b *Bitset[V]) (aOrB *Bitset[V]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset[V]
	if a.pop > b.pop {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aOrB = &Bitset[V]{
		lock: &sync.RWMutex{},
		sets: make(sparse_set.Set[V], len(long.sets)),
		pop:  long.pop,
	}
	copy(aOrB.sets, long.sets)

	it_s, _ := short.Iterate()
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
		sets: sparse_set.Set[V]{},
	}
}

func (s *Bitset[V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.sets = make(sparse_set.Set[V], 0)
}

func (s *Bitset[V]) Copy() *Bitset[V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	clone := New[V]()
	clone.pop = s.pop
	copy(clone.sets, s.sets)

	return clone
}

// Get implements bitset.Bitset
func (s *Bitset[V]) Get(index V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, st := range s.sets {
		if st.Contains(index) {
			return true
		}
	}

	return false
}

// Set implements bitset.Bitset
func (s *Bitset[V]) Set(indices ...V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if s.sets.Insert(index) {
			s.pop++
		}
	}
}

// Unset implements bitset.Bitset
func (s *Bitset[V]) Unset(indices ...V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if s.sets.Remove(index) {
			s.pop--
		}
	}
}

var (
	_ bitset.Bitset[byte]                = (*Bitset[byte])(nil)
	_ bitset.Binary[rune, *Bitset[rune]] = (*Bitset[rune])(nil)
	_ bitset.Inspect[uint]               = (*Bitset[uint])(nil)
)
