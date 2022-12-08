package iterable_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zblach/go-bitset/dense/bits"
	"github.com/zblach/go-bitset/dense/bools"
	"github.com/zblach/go-bitset/iterable"
)

func Test_AndIterable(t *testing.T) {
	a := bits.New[uint8, rune](0)
	a.Set(1, 2, 4, 8, 16, 22)

	b := bits.New[uint32, rune](0)
	b.Set(2, 4, 6, 8, 10, 22)

	c := bits.New[uint16, rune](0)
	c.Set(1, 2, 3, 5, 8, 11, 13, 22)

	andIter := iterable.And[rune](a, b, c)

	assert.EqualValues(t, []rune{2, 8, 22}, iterable.Values(andIter))
}

func Test_OrIterable(t *testing.T) {
	a := bits.New[uint8, rune](0)
	a.Set(1, 2, 4, 8, 16, 22)

	b := bits.New[uint32, rune](0)
	b.Set(2, 4, 6, 8, 10, 22)

	c := bits.New[uint16, rune](0)
	c.Set(1, 2, 3, 5, 8, 11, 13, 22)

	andIter := iterable.Or[rune](a, b, c)

	assert.EqualValues(t, []rune{1, 2, 3, 4, 5, 6, 8, 10, 11, 13, 16, 22}, iterable.Values(andIter))
}

func Test_HeterogenousCopy(t *testing.T) {
	src := bools.New[uint](30)
	dst := bits.NewUint(0)

	src.Set(1, 2, 3, 4, 10, 11, 13, 15)

	iterable.Copy[uint](dst, src)

	assert.Equal(t, dst.Pop(), uint(8))

}
