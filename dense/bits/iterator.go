package bits

import (
	"sync"
	"unsafe"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/iterable"

	mb "math/bits"
)

func (s *Bitset[W, V]) Iterate() (iterable.Iter[V], uint) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator[W, V]{
		// explicit construction to avoid window recalculation
		Bitset: Bitset[W, V]{
			lock: &sync.RWMutex{},
			bits: make([]W, len(s.bits)),
			pop:  s.pop,
		},
	}
	copy(it.bits, s.bits)

	return it, uint(len(s.bits))
}

type Iterator[W Width, V bitset.Value] struct {
	Bitset[W, V]

	index    uint
	wordBits []V
}

func (it *Iterator[W, V]) Next() (V, bool) {
	it.lock.Lock()
	defer it.lock.Unlock()

	if len(it.wordBits) == 0 {
		var window W

		// find next non-zero window
		for ; it.index < uint(len(it.bits)); it.index++ {
			window = it.bits[it.index]
			if window != 0 {
				break
			}
		}
		if it.index >= uint(len(it.bits)) {
			return 0, false
		}

		// precompute next values in this window
		it.wordBits = make([]V, mb.OnesCount64(uint64(window)))
		windowSize := uint(unsafe.Sizeof(W(0))) * 8

		wordIndex := 0
		for i := uint(0); i < windowSize; i += 1 {
			if window&(1<<i) != 0 {
				it.wordBits[wordIndex] = V(i + (windowSize * it.index))
				wordIndex++
			}
		}
		it.index++
	}

	val := it.wordBits[0]
	it.wordBits = it.wordBits[1:]

	return val, true
}

var (
	_ iterable.Iter[uint]     = (*Iterator[uint, uint])(nil)
	_ iterable.Iterable[rune] = (*Bitset[uint, rune])(nil)
)
