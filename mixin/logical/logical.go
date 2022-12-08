package logical

import (
	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/iterable"
	"golang.org/x/exp/slices"
)

// BitsetMixin can be included in a bitset definition to get the associated functions for free.
type BitsetMixin[V bitset.Value] struct {
	bitset.Bitset[V]
}

func (l BitsetMixin[V]) Any(val V, vals ...V) bool {
	// TODO: factor out locking
	if l.Get(val) {
		return true
	}
	for _, v := range vals {
		if l.Get(v) {
			return true
		}
	}
	return false
}

func (l BitsetMixin[V]) All(val V, vals ...V) bool {
	// TODO: factor out locking
	if !l.Get(val) {
		return false
	}
	for _, v := range vals {
		if !l.Get(v) {
			return false
		}
	}
	return true
}

type IterableMixin[V bitset.Value] struct {
	iterable.Iterable[V]
}

func (i IterableMixin[V]) Any(val V, vals ...V) bool {
	valcopy := make([]V, 0, len(vals)+1)
	valcopy[0] = val
	copy(valcopy[1:], vals)
	slices.Sort(valcopy)

	it, _ := i.Iterate()
	v, ok := it.Next()

	for ok && len(valcopy) > 0 {
		switch {
		case v == valcopy[0]:
			return true
		case v < valcopy[0]:
			v, ok = it.Next()
		case v > valcopy[0]:
			valcopy = valcopy[1:]
		}
	}

	return false
}

func (i IterableMixin[V]) All(val V, vals ...V) bool {
	valcopy := make([]V, 0, len(vals))
	valcopy[0] = val
	copy(valcopy, vals)
	slices.Sort(valcopy)

	it, _ := i.Iterate()
	v, ok := it.Next()

	for ok && len(valcopy) > 0 {
		switch {
		case v == valcopy[0]:
			valcopy = valcopy[1:]
			fallthrough
		case v < valcopy[0]:
			v, ok = it.Next()
		case v > valcopy[0]:
			return false
		}
	}
	return true
}

type LogicalMixin[V bitset.Value] interface {
	Any(V, ...V) bool
	All(V, ...V) bool
}

var (
	_ LogicalMixin[byte] = (*BitsetMixin[byte])(nil)
	_ LogicalMixin[rune] = (*IterableMixin[rune])(nil)
)
