package bools

import (
	"sync"

	"github.com/zblach/bitset"
)

func (s *Bitset) Iterate() bitset.Iterator {
	it := &Iterator{
		Bitset: Bitset{
			lock: &sync.RWMutex{},
			pop:  s.pop,
			bits: make([]bool, len(s.bits)),
		},
		index: 0,
	}
	copy(it.bits, s.bits)

	return it
}

type Iterator struct {
	Bitset
	index uint
}

func (it *Iterator) Next() (uint, bool) {
	it.lock.Lock()
	defer it.lock.Unlock()

	for ; it.index < uint(len(it.bits)); it.index++ {
		if it.bits[it.index] {
			return it.index, true
		}
	}
	return 0, false
}
