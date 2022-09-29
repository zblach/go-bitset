package bools

import (
	"sync"

	"github.com/zblach/go-bitset"
)

func (s *Bitset[V]) Iterate() bitset.Iter[V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator[V]{
		Bitset: Bitset[V]{
			lock: &sync.RWMutex{},
			pop:  s.pop,
			bits: make([]bool, len(s.bits)),
		},
		index: 0,
	}
	copy(it.bits, s.bits)

	return it
}

type Iterator[V bitset.Value] struct {
	Bitset[V]
	index uint
}

func (it *Iterator[V]) Next() (V, bool) {
	it.lock.Lock()
	defer it.lock.Unlock()

	for ; it.index < uint(len(it.bits)); it.index++ {
		if it.bits[it.index] {
			val := V(it.index)
			it.index++
			return val, true
		}
	}
	return 0, false
}
