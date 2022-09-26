package bits

import (
	"sync"
	"unsafe"

	"github.com/zblach/bitset"

	mb "math/bits"
)

func (s *Bitset[W]) Iterate() bitset.Iterator {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator[W]{
		lock: &sync.RWMutex{},
		copy: make([]W, len(s.bits)),
	}
	copy(it.copy, s.bits)

	return it
}

type Iterator[W Width] struct {
	lock *sync.RWMutex

	copy  []W
	index uint

	wordBits []uint
}

func (it *Iterator[W]) Next() (uint, bool) {
	it.lock.Lock()
	defer it.lock.Unlock()

	if len(it.wordBits) == 0 {
		var window W

		// find next non-zero window
		for ; it.index < uint(len(it.copy)); it.index++ {
			window = it.copy[it.index]
			if window != 0 {
				break
			}
		}
		if it.index >= uint(len(it.copy)) {
			return 0, false
		}

		it.wordBits = make([]uint, mb.OnesCount64(uint64(window)))
		windowSize := uint(unsafe.Sizeof(W(0))) * 8

		wordIndex := 0
		for i := uint(0); i < windowSize; i += 1 {
			if window&(1<<i) != 0 {
				it.wordBits[wordIndex] = i + (windowSize * it.index)
				wordIndex++
			}
		}
		it.index++
	}

	val := it.wordBits[0]
	it.wordBits = it.wordBits[1:]

	return val, true
}
