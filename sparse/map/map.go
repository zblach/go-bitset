package mapset

import (
	"sync"

	"github.com/zblach/bitset"
)

type noneT struct{}

var none = noneT{}

type Bitset struct {
	lock *sync.RWMutex

	values map[uint]noneT
	pop    uint
}

func New() *Bitset {
	return &Bitset{
		lock:   &sync.RWMutex{},
		values: map[uint]noneT{},
	}
}

func (s *Bitset) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.values = map[uint]noneT{}
	s.pop = 0
}

// Get implements bitset.Bitset
func (s *Bitset) Get(index uint) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.values[index]
	return ok
}

// Set implements bitset.Bitset
func (s *Bitset) Set(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return s
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if _, ok := s.values[index]; !ok {
			s.values[index] = none
			s.pop += 1
		}
	}

	return s
}

// Unset implements bitset.Bitset
func (s *Bitset) Unset(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return s
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if _, ok := s.values[index]; ok {
			delete(s.values, index)
			s.pop -= 1
		}
	}

	return s
}

// And implements bitset.Logical
func (a *Bitset) And(b *Bitset) (aAndB *Bitset) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset
	if a.pop > b.pop {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aAndB = New()

	for v := range short.values {
		if _, ok := long.values[v]; ok {
			aAndB.values[v] = none
			aAndB.pop += 1
		}
	}

	return
}

// Or implements bitset.Logical
func (a *Bitset) Or(b *Bitset) (aOrB *Bitset) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset
	if a.pop > b.pop {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aOrB = New()

	for v := range long.values {
		aOrB.values[v] = none
	}
	aOrB.pop = long.pop

	for v := range short.values {
		if _, ok := aOrB.values[v]; !ok {
			aOrB.values[v] = none
			aOrB.pop += 1
		}
	}

	return
}

// Cap implements bitset.Inspect
func (s *Bitset) Cap() int {
	return s.Len()
}

// Len implements bitset.Inspect
func (s *Bitset) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.values)
}

// Pop implements bitset.Inspect
func (s *Bitset) Pop() uint {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.pop
}

var _ bitset.Bitset = (*Bitset)(nil)
var _ bitset.Logical[*Bitset] = (*Bitset)(nil)
var _ bitset.Inspect = (*Bitset)(nil)
