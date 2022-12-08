package rangeset

import (
	"sync"

	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/iterable"
	"github.com/zblach/go-bitset/sparse/range/sparse_set"
)

// Iterate implements bitset.Inspect
func (s *Bitset[V]) Iterate() (iterable.Iter[V], uint) {
	it := &Iterator[V]{
		b: &Bitset[V]{
			lock: &sync.RWMutex{},
			sets: make(sparse_set.Set[V], len(s.sets)),
		},
		setRange: *sparse_set.NewRange[V](1, 0), // illegal. will be replaced on first call
	}
	copy(it.b.sets, s.sets)

	return it, s.pop
}

type Iterator[V bitset.Value] struct {
	b *Bitset[V]

	setIndex int
	setRange sparse_set.Range[V]
}

func (it *Iterator[V]) Next() (V, bool) {
	it.b.lock.Lock()
	defer it.b.lock.Unlock()

	if it.setRange.Start > it.setRange.End {
		if it.setIndex >= len(it.b.sets) {
			return 0, false
		}
		it.setRange.Start = it.b.sets[it.setIndex].Start
		it.setRange.End = it.b.sets[it.setIndex].End
		it.setIndex++
	}

	val := it.setRange.Start
	it.setRange.Start++
	return val, true
}

var (
	_ iterable.Iter[uint]     = (*Iterator[uint])(nil)
	_ iterable.Iterable[rune] = (*Bitset[rune])(nil)
)
