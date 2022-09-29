package bools

import (
	"sync"

	"github.com/zblach/go-bitset"
)

// Bitset is a boolean-slice-backed bit array. It favors speed over size.
type Bitset[V bitset.Value] struct {
	lock *sync.RWMutex
	bits []bool

	pop uint
}

// New creates a new boolean bitset with an initial size of size.
func New[V bitset.Value](size uint) *Bitset[V] {
	return &Bitset[V]{
		lock: &sync.RWMutex{},
		bits: make([]bool, size),
	}
}

func (s *Bitset[V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.bits = make([]bool, 0)
	s.pop = 0
}

// Get returns whether or not a value is set in the underlying bool slice.
// Getting a value outside of what's stored automatically returns false.
func (s *Bitset[V]) Get(index V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if uint(index) >= uint(len(s.bits)) {
		return false
	}
	return s.bits[index]
}

// Set one or more values in the bitset.
// The bitset will be expanded if necessary.
func (s *Bitset[V]) Set(indices ...V) bitset.Bitset[V] {
	if len(indices) == 0 {
		return s
	}

	maxIndex := V(0)
	for _, index := range indices {
		if index > maxIndex {
			maxIndex = index
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.growright(uint(maxIndex))

	for _, index := range indices {
		if !s.bits[index] {
			s.bits[index] = true
			s.pop += 1
		}
	}
	return s
}

// Unset one or more values in the bitset.
// Indices outside of range are ignored.
func (s *Bitset[V]) Unset(indices ...V) bitset.Bitset[V] {
	if len(indices) == 0 {
		return s
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		if uint(index) >= uint(len(s.bits)) {
			continue
		}
		if s.bits[index] {
			s.bits[index] = false
			s.pop -= 1
		}
	}
	return s
}

// growright expands the underlying storage, if necessary
func (b *Bitset[V]) growright(newSize uint) {
	ulen := uint(len(b.bits))
	if newSize >= ulen {
		b.bits = append(b.bits, make([]bool, (newSize-ulen+1))...)
	}
}

var _ bitset.Bitset[uint] = (*Bitset[uint])(nil)

func (a *Bitset[V]) And(b *Bitset[V]) *Bitset[V] {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var minSize uint
	var short, long *Bitset[V]
	if len(a.bits) > len(b.bits) {
		minSize = uint(len(b.bits))
		short, long = b, a
	} else {
		minSize = uint(len(a.bits))
		short, long = a, b
	}

	result := &Bitset[V]{
		lock: &sync.RWMutex{},
		bits: make([]bool, minSize),
	}

	for i, v := range short.bits {
		if v && long.bits[i] {
			result.bits[i] = true
			result.pop += 1
		}
	}

	return result
}

func (a *Bitset[V]) Or(b *Bitset[V]) *Bitset[V] {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var maxSize uint
	var short, long *Bitset[V]
	if len(a.bits) > len(b.bits) {
		maxSize = uint(len(a.bits))
		short, long = b, a
	} else {
		maxSize = uint(len(b.bits))
		short, long = a, b
	}

	result := &Bitset[V]{
		lock: &sync.RWMutex{},
		bits: make([]bool, maxSize),
		pop:  long.pop,
	}
	copy(result.bits, long.bits)
	for i, v := range short.bits {
		if v && !result.bits[i] {
			result.bits[i] = true
			result.pop += 1
		}
	}

	return result
}

var _ bitset.Logical[uint, *Bitset[uint]] = (*Bitset[uint])(nil)

func (b *Bitset[V]) Len() int {
	return len(b.bits)
}
func (b *Bitset[V]) Cap() int {
	return cap(b.bits)
}
func (b *Bitset[V]) Pop() uint {
	return b.pop
}

var _ bitset.Inspect[uint] = (*Bitset[uint])(nil)
