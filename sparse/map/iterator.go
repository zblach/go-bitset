package mapset

import (
	"sort"

	"github.com/zblach/bitset"
)

func (s *Bitset[V]) Iterate() bitset.Iter[V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator[V]{
		keys: make([]V, s.pop),
	}

	i := 0
	for k := range s.values {
		it.keys[i] = k
		i++
	}

	sort.Slice(it.keys, func(i, j int) bool {
		return it.keys[i] < it.keys[j]
	})

	return it
}

type Iterator[V bitset.Value] struct {
	keys  []V
	index uint
}

func (it *Iterator[V]) Next() (V, bool) {
	if it.index >= uint(len(it.keys)) {
		return 0, false
	}
	val := it.keys[it.index]
	it.index++
	return val, true
}

var _ bitset.Iter[uint16] = (*Iterator[uint16])(nil)
