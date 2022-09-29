package rangeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUnset(t *testing.T) {
	s := New[uint]()

	s.Set(4).Unset(6).Set(5)

	assert.Equal(t, sparseSet[uint]{{4, 5}}, s.sets)
}

func TestSplit(t *testing.T) {
	s := New[uint]()

	s.Set(2).Set(3).Set(4)

	assert.Equal(t, sparseSet[uint]{{2, 4}}, s.sets)

	s.Unset(3)

	assert.Equal(t, sparseSet[uint]{{2, 2}, {4, 4}}, s.sets)
}

func TestSparseSetInsert(t *testing.T) {
	ss := sparseSet[uint]{}

	ss.insert(4)
	ss.insert(7)
	ss.insert(2)
	assert.EqualValues(t, sparseSet[uint]{{2, 2}, {4, 4}, {7, 7}}, ss)
	// ordered

	ss.insert(4)
	assert.EqualValues(t, sparseSet[uint]{{2, 2}, {4, 4}, {7, 7}}, ss)
	// no-op

	ss.insert(5)
	assert.EqualValues(t, sparseSet[uint]{{2, 2}, {4, 5}, {7, 7}}, ss)
	// extend {4, 4} -> {4, 5}, no merge

	ss.insert(3)
	assert.EqualValues(t, sparseSet[uint]{{2, 5}, {7, 7}}, ss)
	// extend {2, 3} -> {2, 4}, merge with {4, 5} to become {2, 5}

	ss.insert(10)
	assert.EqualValues(t, sparseSet[uint]{{2, 5}, {7, 7}, {10, 10}}, ss)
	// add {10, 10}

	ss.insert(6)
	assert.EqualValues(t, sparseSet[uint]{{2, 7}, {10, 10}}, ss)
	// extend {2, 5} -> {2, 6}, merge with {7, 7} to become {2, 7}

	ss.insert(1)
	assert.EqualValues(t, sparseSet[uint]{{1, 7}, {10, 10}}, ss)
	// extend {2, 7} -> {1, 7}, no merge as it's first elem

	ss.insert(9)
	assert.EqualValues(t, sparseSet[uint]{{1, 7}, {9, 10}}, ss)
	// extend {10, 10} -> {9, 10}, no merge

	ss.insert(8)
	assert.EqualValues(t, sparseSet[uint]{{1, 10}}, ss)
	// extend {1, 7} -> {1, 8}, merge with {9, 10}

	ss.insert(11)
	assert.EqualValues(t, sparseSet[uint]{{1, 11}}, ss)
	// extend {1, 10} -> {1, 11}. no merge as it's final element

}

func TestSparseSetRemove(t *testing.T) {
	ss := sparseSet[uint]{{1, 2}, {4, 8}, {11, 15}}

	ss.remove(0)
	assert.EqualValues(t, sparseSet[uint]{{1, 2}, {4, 8}, {11, 15}}, ss)
	// no-op

	ss.remove(1)
	assert.EqualValues(t, sparseSet[uint]{{2, 2}, {4, 8}, {11, 15}}, ss)
	// modify {1, 2} -> {2, 2}

	ss.remove(2)
	assert.EqualValues(t, sparseSet[uint]{{4, 8}, {11, 15}}, ss)
	// remove {2, 2}

	ss.remove(12)
	assert.EqualValues(t, sparseSet[uint]{{4, 8}, {11, 11}, {13, 15}}, ss)
	// split {11, 15} -> {11, 11} + {13, 15}

	ss.remove(6)
	assert.EqualValues(t, sparseSet[uint]{{4, 5}, {7, 8}, {11, 11}, {13, 15}}, ss)
	// split {4, 8} -> {4, 5} + {7, 8}

	ss.remove(13)
	assert.EqualValues(t, sparseSet[uint]{{4, 5}, {7, 8}, {11, 11}, {14, 15}}, ss)

	ss.remove(15)
	ss.remove(14)
	assert.EqualValues(t, sparseSet[uint]{{4, 5}, {7, 8}, {11, 11}}, ss)
	// remove last element
}

func TestRangeIterator(t *testing.T) {
	ss := New[uint]()
	ss.Set(1, 2, 4, 5, 6, 7, 8)

	assert.EqualValues(t, sparseSet[uint]{{1, 2}, {4, 8}}, ss.sets)

	vs := make([]uint, 0, ss.Pop())

	it := ss.Iterate()
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

	assert.EqualValues(t, sparseSet[uint]{{6, 6}, {8, 8}}, aAndB.sets)
	assert.EqualValues(t, sparseSet[uint]{{1, 4}, {6, 8}}, aOrB.sets)
}
