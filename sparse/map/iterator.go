package mapset

import (
	"sort"

	"github.com/zblach/bitset"
)

func (s *Bitset) Iterate() bitset.Iterator {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator{
		keys: make([]uint, s.pop),
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

type Iterator struct {
	keys  []uint
	index uint
}

func (it *Iterator) Next() (uint, bool) {
	if it.index >= uint(len(it.keys)) {
		return 0, false
	}
	val := it.keys[it.index]
	it.index++
	return val, true
}

var _ bitset.Iterator = (*Iterator)(nil)
