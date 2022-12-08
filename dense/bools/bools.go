package bools

import (
	"sync"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/mixin/logical"
)

// Bitset is a boolean-slice-backed bit array. It favors speed over size.
type Bitset[V bitset.Value] struct {
	logical.IterableMixin[V]

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

func (s *Bitset[V]) Copy() *Bitset[V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	clone := New[V](s.pop)
	copy(clone.bits, s.bits)

	return clone
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
func (s *Bitset[V]) Set(indices ...V) {
	if len(indices) == 0 {
		return
	}

	maxIndex := indices[0]
	for i := 1; i < len(indices); i++ {
		if indices[i] > maxIndex {
			maxIndex = indices[i]
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.growright(uint64(maxIndex))

	for _, index := range indices {
		if !s.bits[index] {
			s.bits[index] = true
			s.pop += 1
		}
	}
}

// Unset one or more values in the bitset.
// Indices outside of range are ignored.
func (s *Bitset[V]) Unset(indices ...V) {
	if len(indices) == 0 {
		return
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
}

// growright expands the underlying storage, if necessary
func (b *Bitset[V]) growright(newSize uint64) {
	ulen := uint64(len(b.bits))
	if newSize >= ulen {
		b.bits = append(b.bits, make([]bool, (newSize-ulen+1))...)
	}
}

var _ bitset.Bitset[uint] = (*Bitset[uint])(nil)

func (a *Bitset[V]) And(b *Bitset[V]) (aAndB *Bitset[V]) {
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

	aAndB = New[V](minSize)

	for i, v := range short.bits {
		if v && long.bits[i] {
			aAndB.bits[i] = true
			aAndB.pop += 1
		}
	}

	return
}

func (a *Bitset[V]) Or(b *Bitset[V]) (aOrB *Bitset[V]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset[V]
	if len(a.bits) > len(b.bits) {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aOrB = New[V](long.pop)

	result := &Bitset[V]{
		lock: &sync.RWMutex{},
		bits: make([]bool, len(long.bits)),
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

var _ bitset.Binary[uint, *Bitset[uint]] = (*Bitset[uint])(nil)

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
