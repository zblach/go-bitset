package bools

import (
	"sync"

	"github.com/zblach/bitset"
)

// Bitset is a boolean-slice-backed bit array. It favors speed over size.
type Bitset struct {
	lock *sync.RWMutex
	bits []bool

	pop uint
}

// New creates a new boolean bitset with an initial size of size.
func New(size uint) *Bitset {
	return &Bitset{
		lock: &sync.RWMutex{},
		bits: make([]bool, size),
	}
}

func (s *Bitset) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.bits = make([]bool, 0)
	s.pop = 0
}

// Get returns whether or not a value is set in the underlying bool slice.
// Getting a value outside of what's stored automatically returns false.
func (b *Bitset) Get(index uint) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if index >= uint(len(b.bits)) {
		return false
	}
	return b.bits[index]
}

// Set one or more values in the bitset.
// The bitset will be expanded if necessary.
func (b *Bitset) Set(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return b
	}

	maxIndex := uint(0)
	for _, index := range indices {
		if index > maxIndex {
			maxIndex = index
		}
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	b.growright(maxIndex)

	for _, index := range indices {
		if !b.bits[index] {
			b.bits[index] = true
			b.pop += 1
		}
	}
	return b
}

// Unset one or more values in the bitset.
// Indices outside of range are ignored.
func (b *Bitset) Unset(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return b
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	for _, index := range indices {
		if index >= uint(len(b.bits)) {
			continue
		}
		if b.bits[index] {
			b.bits[index] = false
			b.pop -= 1
		}
	}
	return b
}

// growright expands the underlying storage, if necessary
func (b *Bitset) growright(newSize uint) {
	ulen := uint(len(b.bits))
	if newSize >= ulen {
		b.bits = append(b.bits, make([]bool, (newSize-ulen+1))...)
	}
}

var _ bitset.Bitset = (*Bitset)(nil)

func (a *Bitset) And(b *Bitset) *Bitset {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var minSize uint
	var short, long *Bitset
	if len(a.bits) > len(b.bits) {
		minSize = uint(len(b.bits))
		short, long = b, a
	} else {
		minSize = uint(len(a.bits))
		short, long = a, b
	}

	result := &Bitset{
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

func (a *Bitset) Or(b *Bitset) *Bitset {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var maxSize uint
	var short, long *Bitset
	if len(a.bits) > len(b.bits) {
		maxSize = uint(len(a.bits))
		short, long = b, a
	} else {
		maxSize = uint(len(b.bits))
		short, long = a, b
	}

	result := &Bitset{
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

var _ bitset.Logical[*Bitset] = (*Bitset)(nil)

func (b *Bitset) Len() int {
	return len(b.bits)
}
func (b *Bitset) Cap() int {
	return cap(b.bits)
}
func (b *Bitset) Pop() uint {
	return b.pop
}

var _ bitset.Inspect = (*Bitset)(nil)
