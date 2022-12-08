package rangeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zblach/go-bitset/sparse/range/sparse_set"
)

func TestSetUnset(t *testing.T) {
	s := New[uint]()

	s.Set(4, 5)
	s.Unset(6)

	assert.Equal(t, sparse_set.Set[uint]{{4, 5}}, s.sets)
}

func TestSplit(t *testing.T) {
	s := New[uint]()

	s.Set(2, 3, 4)

	assert.Equal(t, sparse_set.Set[uint]{{2, 4}}, s.sets)

	s.Unset(3)

	assert.Equal(t, sparse_set.Set[uint]{{2, 2}, {4, 4}}, s.sets)
}

func TestSparseSetInsert(t *testing.T) {
	ss := sparse_set.Set[uint]{}

	ss.Insert(4)
	ss.Insert(7)
	ss.Insert(2)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 2}, {4, 4}, {7, 7}}, ss)
	// ordered

	ss.Insert(4)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 2}, {4, 4}, {7, 7}}, ss)
	// no-op

	ss.Insert(5)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 2}, {4, 5}, {7, 7}}, ss)
	// extend {4, 4} -> {4, 5}, no merge

	ss.Insert(3)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 5}, {7, 7}}, ss)
	// extend {2, 3} -> {2, 4}, merge with {4, 5} to become {2, 5}

	ss.Insert(10)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 5}, {7, 7}, {10, 10}}, ss)
	// add {10, 10}

	ss.Insert(6)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 7}, {10, 10}}, ss)
	// extend {2, 5} -> {2, 6}, merge with {7, 7} to become {2, 7}

	ss.Insert(1)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 7}, {10, 10}}, ss)
	// extend {2, 7} -> {1, 7}, no merge as it's first elem

	ss.Insert(9)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 7}, {9, 10}}, ss)
	// extend {10, 10} -> {9, 10}, no merge

	ss.Insert(8)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 10}}, ss)
	// extend {1, 7} -> {1, 8}, merge with {9, 10}

	ss.Insert(11)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 11}}, ss)
	// extend {1, 10} -> {1, 11}. no merge as it's final element

}

func TestSparseSetRemove(t *testing.T) {
	ss := sparse_set.Set[uint]{{1, 2}, {4, 8}, {11, 15}}

	ss.Remove(0)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 2}, {4, 8}, {11, 15}}, ss)
	// no-op

	ss.Remove(1)
	assert.EqualValues(t, sparse_set.Set[uint]{{2, 2}, {4, 8}, {11, 15}}, ss)
	// modify {1, 2} -> {2, 2}

	ss.Remove(2)
	assert.EqualValues(t, sparse_set.Set[uint]{{4, 8}, {11, 15}}, ss)
	// Remove {2, 2}

	ss.Remove(12)
	assert.EqualValues(t, sparse_set.Set[uint]{{4, 8}, {11, 11}, {13, 15}}, ss)
	// split {11, 15} -> {11, 11} + {13, 15}

	ss.Remove(6)
	assert.EqualValues(t, sparse_set.Set[uint]{{4, 5}, {7, 8}, {11, 11}, {13, 15}}, ss)
	// split {4, 8} -> {4, 5} + {7, 8}

	ss.Remove(13)
	assert.EqualValues(t, sparse_set.Set[uint]{{4, 5}, {7, 8}, {11, 11}, {14, 15}}, ss)

	ss.Remove(15)
	ss.Remove(14)
	assert.EqualValues(t, sparse_set.Set[uint]{{4, 5}, {7, 8}, {11, 11}}, ss)
	// Remove last element
}

func TestRangeIterator(t *testing.T) {
	ss := New[uint]()
	ss.Set(1, 2, 4, 5, 6, 7, 8)

	assert.EqualValues(t, sparse_set.Set[uint]{{1, 2}, {4, 8}}, ss.sets)

	vs := make([]uint, 0, ss.Pop())

	it, _ := ss.Iterate()
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		vs = append(vs, n)
	}

	assert.EqualValues(t, []uint{1, 2, 4, 5, 6, 7, 8}, vs)
}

func TestLogical(t *testing.T) {
	a := New[uint]()
	a.Set(1, 3, 6, 8)

	b := New[uint]()
	b.Set(2, 4, 6, 7, 8)

	aAndB := a.And(b)
	aOrB := a.Or(b)

	assert.EqualValues(t, sparse_set.Set[uint]{{6, 6}, {8, 8}}, aAndB.sets)
	assert.EqualValues(t, sparse_set.Set[uint]{{1, 4}, {6, 8}}, aOrB.sets)
}
