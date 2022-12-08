package bools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bools_Sizes(t *testing.T) {
	assert.Len(t, New[uint](1).bits, 1)
	assert.Len(t, New[uint](8).bits, 8)
	assert.Len(t, New[uint8](64).bits, 64)

	assert.Len(t, New[rune](0).bits, 0)
}

const (
	T = true
	F = false
)

func Test_Bools_SetUnset(t *testing.T) {
	e := New[uint8](8)

	e.Set(0, 3, 7)

	assert.Len(t, e.bits, 8)
	assert.Equal(t, []bool{T, F, F, T, F, F, F, T}, e.bits)

	e.Set(8)
	assert.Len(t, e.bits, 9)
	assert.Equal(t, T, e.bits[8])

	e.Set(24)
	assert.Len(t, e.bits, 25)
	assert.Equal(t, T, e.bits[24])
	assert.True(t, e.Get(24))

	e.Unset(8)
	assert.Len(t, e.bits, 25)
	assert.Equal(t, F, e.bits[8])

	e.Unset(64)
	assert.Len(t, e.bits, 25)
	assert.False(t, e.Get(64))
}

func Test_Bools_Logical(t *testing.T) {
	a := New[uint](0)
	b := New[uint](0)

	a.Set(1, 3, 5, 6, 7)
	b.Set(0, 2, 4, 6, 7)

	aAndB := a.And(b)
	aOrB := a.Or(b)

	assert.Len(t, aAndB.bits, 8)
	assert.Len(t, aOrB.bits, 8)

	assert.Equal(t, []bool{T, T, T, T, T, T, T, T}, aOrB.bits)
	assert.Equal(t, []bool{F, F, F, F, F, F, T, T}, aAndB.bits)

	b.Set(10)

	assert.Len(t, a.bits, 8)
	assert.Len(t, b.bits, 11)

	aAndB = a.And(b)
	aOrB = a.Or(b)

	assert.Len(t, aAndB.bits, 8)
	assert.Len(t, aOrB.bits, 11)

	bAndA := b.And(a)
	bOrA := b.Or(a)

	assert.Len(t, bAndA.bits, 8)
	assert.Len(t, bOrA.bits, 11)
}

func Test_Bools_Inspect(t *testing.T) {
	a := New[uint](2)
	a.Set(1, 2, 3, 1)
	a.Unset(2)

	assert.Equal(t, 4, a.Len())
	assert.GreaterOrEqual(t, a.Cap(), 4)
	assert.Equal(t, uint(2), a.Pop())
}
