package rangeset

import (
	"sync"

	"github.com/zblach/bitset"
)

// Iterate implements bitset.Inspect
func (s *Bitset[V]) Iterate() bitset.Iter[V] {
	it := &Iterator[V]{
		b: &Bitset[V]{
			lock: &sync.RWMutex{},
			sets: make(sparseSet[V], len(s.sets)),
		},
		setRange: sparseRange[V]{1, 0}, // illegal. will be replaced on first call
	}
	copy(it.b.sets, s.sets)

	return it
}

type Iterator[V bitset.Value] struct {
	b *Bitset[V]

	setIndex int
	setRange sparseRange[V]
}

func (it *Iterator[V]) Next() (V, bool) {
	it.b.lock.Lock()
	defer it.b.lock.Unlock()

	if it.setRange.start > it.setRange.end {
		if it.setIndex >= len(it.b.sets) {
			return 0, false
		}
		it.setRange.start = it.b.sets[it.setIndex].start
		it.setRange.end = it.b.sets[it.setIndex].end
		it.setIndex++
	}

	val := it.setRange.start
	it.setRange.start++
	return val, true
}
