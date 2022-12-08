package bools

import (
	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/iterable"
)

func (s *Bitset[V]) Iterate() (iterable.Iter[V], uint) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	len := uint(len(s.bits))

	it := &Iterator[V]{
		Bitset: *New[V](len),
	}
	it.pop = s.pop
	copy(it.bits, s.bits)

	return it, len
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

var (
	_ iterable.Iter[uint]     = (*Iterator[uint])(nil)
	_ iterable.Iterable[rune] = (*Bitset[rune])(nil)
)
