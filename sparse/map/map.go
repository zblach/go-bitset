package mapset

import (
	"sync"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/mixin/logical"
)

type noneT struct{}

var none = noneT{}

type Bitset[V bitset.Value] struct {
	logical.IterableMixin[V]

	lock *sync.RWMutex

	values map[V]noneT
	pop    uint
}

func New[V bitset.Value]() *Bitset[V] {
	return &Bitset[V]{
		lock:   &sync.RWMutex{},
		values: map[V]noneT{},
	}
}

func (s *Bitset[V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values = map[V]noneT{}
	s.pop = 0
}

func (s *Bitset[V]) Copy() *Bitset[V] {
	s.lock.Lock()
	defer s.lock.RUnlock()

	clone := New[V]()
	clone.pop = s.pop
	for k, v := range s.values {
		clone.values[k] = v
	}

	return clone
}

// Get implements bitset.Bitset
func (s *Bitset[V]) Get(index V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.values[index]
	return ok
}

// Set implements bitset.Bitset
func (s *Bitset[V]) Set(indices ...V) {
	if len(indices) == 0 {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if _, ok := s.values[index]; !ok {
			s.values[index] = none
			s.pop += 1
		}
	}
}

// Unset implements bitset.Bitset
func (s *Bitset[V]) Unset(indices ...V) {
	if len(indices) == 0 {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if _, ok := s.values[index]; ok {
			delete(s.values, index)
			s.pop -= 1
		}
	}
}

// And implements bitset.Logical
func (a *Bitset[V]) And(b *Bitset[V]) (aAndB *Bitset[V]) {
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

	aAndB = New[V]()

	for v := range short.values {
		if _, ok := long.values[v]; ok {
			aAndB.values[v] = none
			aAndB.pop += 1
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

	aOrB = long.Copy()

	for v := range short.values {
		if _, ok := aOrB.values[v]; !ok {
			aOrB.values[v] = none
			aOrB.pop += 1
		}
	}

	return
}

// Cap implements bitset.Inspect
func (s *Bitset[V]) Cap() int {
	return s.Len()
}

// Len implements bitset.Inspect
func (s *Bitset[V]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.values)
}

// Pop implements bitset.Inspect
func (s *Bitset[V]) Pop() uint {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.pop
}

// Interface adherence. Randomly-selected V types
var (
	_ bitset.Bitset[uint]                = (*Bitset[uint])(nil)
	_ bitset.Binary[rune, *Bitset[rune]] = (*Bitset[rune])(nil)
	_ bitset.Inspect[uint8]              = (*Bitset[uint8])(nil)
)
