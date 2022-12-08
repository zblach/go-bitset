package mapset

import (
	"sort"
	"sync"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/iterable"
)

func (s *Bitset[V]) Iterate() (iterable.Iter[V], uint) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	it := &Iterator[V]{
		lock: &sync.RWMutex{},
		keys: make([]V, 0, s.pop),
	}

	for k := range s.values {
		it.keys = append(it.keys, k)
	}

	// All other iterators are in order, but map access is random.
	// For consistency's sake, we sort the keys before we iterate.
	sort.Slice(it.keys, func(i, j int) bool {
		return it.keys[i] < it.keys[j]
	})

	return it, s.pop
}

type Iterator[V bitset.Value] struct {
	lock  *sync.RWMutex
	keys  []V
	index uint
}

func (it *Iterator[V]) Next() (V, bool) {
	it.lock.Lock()
	defer it.lock.Unlock()

	if it.index >= uint(len(it.keys)) {
		return 0, false
	}
	val := it.keys[it.index]
	it.index++
	return val, true
}

var (
	_ iterable.Iter[uint]     = (*Iterator[uint])(nil)
	_ iterable.Iterable[rune] = (*Bitset[rune])(nil)
)
