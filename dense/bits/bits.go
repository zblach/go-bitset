package bits

import (
	"sync"
	"unsafe"

	mb "math/bits"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/mixin/logical"
)

// Width is the underlying word-size for storing bit fields.
// It mostly matters for allocation and indexing
type Width interface {
	uint8 | uint16 | uint32 | uint64 | // explicit sizes
		uintptr | ~uint // machine-optimized sizes
}

// Basic internal type aliases.
type (
	Uint   = Bitset[uint, uint]
	Uint8  = Bitset[uint8, uint]
	Uint16 = Bitset[uint16, uint]
	Uint32 = Bitset[uint32, uint]
	Uint64 = Bitset[uint64, uint]
)

// NOTE: can't currently use a generic parameter as a type alias.

// handy aliases for instantiation
var (
	NewUint   = New[uint, uint]
	NewUint8  = New[uint8, uint]
	NewUint16 = New[uint16, uint]
	NewUint32 = New[uint32, uint]
	NewUint64 = New[uint64, uint]
)

// Bitset is a threadsafe container for storing a set of bits.
type Bitset[W Width, V bitset.Value] struct {
	logical.IterableMixin[V]

	lock *sync.RWMutex

	bits []W
	pop  uint
}

// New instantiates a new bitset with an initial size of size.
// This size parameter refers to the number of bitset elements, not the underlying storage.
func New[W Width, V bitset.Value](size uint) *Bitset[W, V] {
	var width uint
	if size == 0 {
		width = 0
	} else {
		var bits W
		width, bits = indexToTuple[W](size - 1)
		if bits > 0 {
			width += 1
		}
	}
	return &Bitset[W, V]{
		lock: &sync.RWMutex{},
		bits: make([]W, width),
	}
}

// Clear unsets all elements in the bitset, and sets the internal size to zero.
func (s *Bitset[W, V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.bits = make([]W, 0)
	s.pop = 0
}

// Copy returns a deep copy of the bitset.
func (s *Bitset[W, V]) Copy() *Bitset[W, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	clone := &Bitset[W, V]{
		lock: &sync.RWMutex{},
		pop:  s.pop,
		bits: make([]W, len(s.bits)),
	}
	copy(clone.bits, s.bits)

	return clone
}

// Get returns whether or not a value is set in the underlying bitset.
// Getting a value outside of what's stored automatically returns false.
func (s *Bitset[W, V]) Get(index V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	elem, bit := indexToTuple[W](uint(index))
	if elem >= uint(len(s.bits)) {
		return false
	}
	return (s.bits[elem] & bit) != 0
}

// Set one or more values in the bitset.
// The bitset will be expanded if necessary.
func (s *Bitset[W, V]) Set(indices ...V) {
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

	s.growright(uint(maxIndex))

	for _, index := range indices {
		elem, bit := indexToTuple[W](uint(index))
		if (s.bits[elem] & bit) == 0 {
			s.bits[elem] |= bit
			s.pop += 1
		}
	}
}

// Unset one or more values in the bitset.
// Indices outside of range are ignored.
func (s *Bitset[W, V]) Unset(indices ...V) {
	if len(indices) == 0 {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		elem, bit := indexToTuple[W](uint(index))
		if elem >= uint(len(s.bits)) {
			continue
		}

		if (s.bits[elem] & bit) != 0 {
			s.bits[elem] &= ^bit
			s.pop -= 1
		}
	}
}

// indexToTuple is a utility function to compute the memory element and bitmask, based on index.
func indexToTuple[W Width](index uint) (elem uint, bit W) {
	wbits := uint(unsafe.Sizeof(W(0)) << 3)

	return index / wbits, 1 << (index & (wbits - 1)) // equiv to index % wbits
}

// growright expands the underlying storage, if necessary
func (s *Bitset[W, V]) growright(newSize uint) {
	newLen, _ := indexToTuple[W](newSize)
	ulen := uint(len(s.bits))
	if newLen >= ulen {
		chunk := make([]W, newLen-ulen+1)
		s.bits = append(s.bits, chunk...)
	}
}

// interface adherence validation
var (
	_ bitset.Bitset[uint] = (*Uint)(nil)
	_ bitset.Bitset[uint] = (*Uint8)(nil)
	_ bitset.Bitset[uint] = (*Uint16)(nil)
	_ bitset.Bitset[uint] = (*Uint32)(nil)
	_ bitset.Bitset[uint] = (*Uint64)(nil)

	_ bitset.Binary[uint, *Uint]   = (*Uint)(nil)
	_ bitset.Binary[uint, *Uint8]  = (*Uint8)(nil)
	_ bitset.Binary[uint, *Uint16] = (*Uint16)(nil)
	_ bitset.Binary[uint, *Uint32] = (*Uint32)(nil)
	_ bitset.Binary[uint, *Uint64] = (*Uint64)(nil)

	_ bitset.Inspect[uint] = (*Uint)(nil)
	_ bitset.Inspect[uint] = (*Uint8)(nil)
	_ bitset.Inspect[uint] = (*Uint16)(nil)
	_ bitset.Inspect[uint] = (*Uint32)(nil)
	_ bitset.Inspect[uint] = (*Uint64)(nil)
)

// And computes and returns the intersection of two bitsets.
// It does not modify either bitset.
func (a *Bitset[W, V]) And(b *Bitset[W, V]) (aAndB *Bitset[W, V]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset[W, V]
	if len(a.bits) > len(b.bits) {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aAndB = &Bitset[W, V]{
		lock: &sync.RWMutex{},
		bits: make([]W, len(short.bits)),
	}

	for i, bits := range short.bits {
		aAndB.bits[i] = long.bits[i] & bits
		aAndB.pop += uint(mb.OnesCount64(uint64(aAndB.bits[i])))
	}

	return
}

// Or computes and returns the union of two bitsets.
// It does not modify either bitset.
func (a *Bitset[W, V]) Or(b *Bitset[W, V]) (aOrB *Bitset[W, V]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var short, long *Bitset[W, V]
	if len(a.bits) > len(b.bits) {
		short, long = b, a
	} else {
		short, long = a, b
	}

	aOrB = &Bitset[W, V]{
		lock: &sync.RWMutex{},
		bits: make([]W, len(long.bits)),
	}
	for i, sbits := range short.bits {
		aOrB.bits[i] = long.bits[i] | sbits
		aOrB.pop += uint(mb.OnesCount64(uint64(aOrB.bits[i])))
	}
	for i := len(short.bits); i < len(long.bits); i++ {
		aOrB.bits[i] = long.bits[i]
		aOrB.pop += uint(mb.OnesCount64(uint64(aOrB.bits[i])))
	}

	return
}

// Inspection functions

// Len is the used number of bits in the underlying data store (rounded up to word size).
func (s *Bitset[V, W]) Len() int {
	return len(s.bits) * int(unsafe.Sizeof(W(0))) * 8
}

// Cap is the available number of bits in the underlying data store (rounded up to word size).
func (s *Bitset[V, W]) Cap() int {
	return cap(s.bits) * int(unsafe.Sizeof(W(0))) * 8
}

// Pop is the number of bits set in the underlying data store.
func (s *Bitset[V, W]) Pop() uint {
	return s.pop
}
