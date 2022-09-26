package bits

import (
	"sync"
	"unsafe"

	mb "math/bits"

	"github.com/zblach/bitset"
)

// Width is the underlying word-size for storing bit fields.
// It mostly matters for allocation and indexing
type Width interface {
	uint8 | uint16 | uint32 | uint64 | // explicit sizes
		uintptr | ~uint // machine-optimized sizes
}

type (
	Uint   = Bitset[uint]
	Uint8  = Bitset[uint8]
	Uint16 = Bitset[uint16]
	Uint32 = Bitset[uint32]
	Uint64 = Bitset[uint64]
)

// handy aliases for instantiation
var (
	NewUint   = New[uint]
	NewUint8  = New[uint8]
	NewUint16 = New[uint16]
	NewUint32 = New[uint32]
	NewUint64 = New[uint64]
)

// Bitset is a threadsafe container for storing a set of bits.
type Bitset[W Width] struct {
	lock *sync.RWMutex

	bits []W
	pop  uint
}

// New instantiates a new bitset with an initial size of size.
// This size parameter refers to the number of bitset elements, not the underlying storage.
func New[W Width](size uint) *Bitset[W] {
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
	return &Bitset[W]{
		lock: &sync.RWMutex{},
		bits: make([]W, width),
	}
}

func (s *Bitset[W]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.bits = make([]W, 0)
	s.pop = 0
}

// Get returns whether or not a value is set in the underlying bitset.
// Getting a value outside of what's stored automatically returns false.
func (s *Bitset[W]) Get(index uint) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	elem, bit := indexToTuple[W](index)
	if elem >= uint(len(s.bits)) {
		return false
	}
	return (s.bits[elem] & bit) != 0
}

// Set one or more values in the bitset.
// The bitset will be expanded if necessary.
func (s *Bitset[W]) Set(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return s
	}

	maxIndex := uint(0)
	for _, index := range indices {
		if index > maxIndex {
			maxIndex = index
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.growright(maxIndex)

	for _, index := range indices {
		elem, bit := indexToTuple[W](index)
		if (s.bits[elem] & bit) == 0 {
			s.bits[elem] |= bit
			s.pop += 1
		}
	}

	return s
}

// Unset one or more values in the bitset.
// Indices outside of range are ignored.
func (s *Bitset[W]) Unset(indices ...uint) bitset.Bitset {
	if len(indices) == 0 {
		return s
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, index := range indices {
		elem, bit := indexToTuple[W](index)
		if elem >= uint(len(s.bits)) {
			continue
		}

		if (s.bits[elem] & bit) != 0 {
			s.bits[elem] &= ^bit
			s.pop -= 1
		}
	}
	return s
}

// indexToTuple is a utility function to compute the element and bit, based on index.
func indexToTuple[W Width](index uint) (elem uint, bit W) {
	wbits := uint(unsafe.Sizeof(W(0)) << 3)

	return index / wbits, 1 << (index & (wbits - 1))
}

// growright expands the underlying storage, if necessary
func (s *Bitset[W]) growright(newSize uint) {
	newLen, _ := indexToTuple[W](newSize)
	ulen := uint(len(s.bits))
	if newLen >= ulen {
		chunk := make([]W, newLen-ulen+1)
		s.bits = append(s.bits, chunk...)
	}
}

// interface adherence validation
var (
	_ bitset.Bitset = (*Uint)(nil)
	_ bitset.Bitset = (*Uint8)(nil)
	_ bitset.Bitset = (*Uint16)(nil)
	_ bitset.Bitset = (*Uint32)(nil)
	_ bitset.Bitset = (*Uint64)(nil)

	_ bitset.Logical[*Uint]   = (*Uint)(nil)
	_ bitset.Logical[*Uint8]  = (*Uint8)(nil)
	_ bitset.Logical[*Uint16] = (*Uint16)(nil)
	_ bitset.Logical[*Uint32] = (*Uint32)(nil)
	_ bitset.Logical[*Uint64] = (*Uint64)(nil)

	_ bitset.Inspect = (*Uint)(nil)
	_ bitset.Inspect = (*Uint8)(nil)
	_ bitset.Inspect = (*Uint16)(nil)
	_ bitset.Inspect = (*Uint32)(nil)
	_ bitset.Inspect = (*Uint64)(nil)
)

// And computes and returns the intersection of two bitsets.
// It does not modify either bitset.
func (a *Bitset[W]) And(b *Bitset[W]) (aAndB *Bitset[W]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var minSize uint
	var short, long *Bitset[W]
	if len(a.bits) > len(b.bits) {
		minSize = uint(len(b.bits))
		short, long = b, a
	} else {
		minSize = uint(len(a.bits))
		short, long = a, b
	}

	aAndB = &Bitset[W]{
		lock: &sync.RWMutex{},
		bits: make([]W, minSize),
	}

	for i, bits := range short.bits {
		aAndB.bits[i] = long.bits[i] & bits
		aAndB.pop += uint(mb.OnesCount64(uint64(aAndB.bits[i])))
	}

	return
}

// Or computes and returns the union of two bitsets.
// It does not modify either bitset.
func (a *Bitset[W]) Or(b *Bitset[W]) (aOrB *Bitset[W]) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	b.lock.RLock()
	defer b.lock.RUnlock()

	var maxSize uint
	var short, long *Bitset[W]
	if len(a.bits) > len(b.bits) {
		maxSize = uint(len(a.bits))
		short, long = b, a
	} else {
		maxSize = uint(len(b.bits))
		short, long = a, b
	}

	aOrB = &Bitset[W]{
		lock: &sync.RWMutex{},
		bits: make([]W, maxSize),
	}
	copy(aOrB.bits, long.bits)

	for i, bits := range short.bits {
		aOrB.bits[i] |= bits
		aOrB.pop += uint(mb.OnesCount64(uint64(aOrB.bits[i])))
	}

	return
}

// Inspection functions

// Len is the used number of bits in the underlying data store (rounded up to word size).
func (s *Bitset[W]) Len() int {
	return int(len(s.bits)*int(unsafe.Sizeof(W(0)))) * 8
}

// Cap is the available number of bits in the underlying data store (rounded up to word size).
func (s *Bitset[W]) Cap() int {
	return int(cap(s.bits)*int(unsafe.Sizeof(W(0)))) * 8
}

// Pop is the number of bits set in the underlying data store.
func (s *Bitset[W]) Pop() uint {
	return s.pop
}
