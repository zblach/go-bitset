package bits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Uint8s_Sizes(t *testing.T) {
	assert.Len(t, NewUint8(1).bits, 1)
	assert.Len(t, NewUint8(8).bits, 1)
	assert.Len(t, NewUint8(64).bits, 8)

	assert.Len(t, NewUint8(0).bits, 0)
}

func Test_Uint64s_Sizes(t *testing.T) {
	assert.Len(t, NewUint64(1).bits, 1)
	assert.Len(t, NewUint64(8).bits, 1)
	assert.Len(t, NewUint64(64).bits, 1)

	assert.Len(t, NewUint64(0).bits, 0)
}

func Test_Uint8_SetUnset(t *testing.T) {
	e := NewUint8(8)

	e.Set(0, 3, 7)

	assert.Len(t, e.bits, 1)
	//index bits stored as: 76543210
	assert.Equal(t, uint8(0b10001001), e.bits[0])

	e.Set(8)
	assert.Len(t, e.bits, 2)
	assert.Equal(t, uint8(0b00000001), e.bits[1])

	e.Set(24)
	assert.Len(t, e.bits, 4)
	assert.Equal(t, uint8(0b00000001), e.bits[3])
	assert.True(t, e.Get(24))

	e.Unset(8)
	assert.Len(t, e.bits, 4)
	assert.Equal(t, uint8(0b00000000), e.bits[1])

	e.Unset(64)
	assert.Len(t, e.bits, 4)
	assert.False(t, e.Get(64))
}

func Test_Uint8_Logical(t *testing.T) {
	a := NewUint8(0)
	b := NewUint8(0)

	a.Set(1, 3, 5, 6, 7)
	b.Set(0, 2, 4, 6, 7)

	aAndB := a.And(b)
	aOrB := a.Or(b)

	assert.Len(t, aAndB.bits, 1)
	assert.Len(t, aOrB.bits, 1)

	assert.Equal(t, uint8(0b11111111), aOrB.bits[0])
	assert.Equal(t, uint8(0b11000000), aAndB.bits[0])

	b.Set(10)

	assert.Len(t, a.bits, 1)
	assert.Len(t, b.bits, 2)

	aAndB = a.And(b)
	aOrB = a.Or(b)

	assert.Len(t, aAndB.bits, 1)
	assert.Len(t, aOrB.bits, 2)

	bAndA := b.And(a)
	bOrA := b.Or(a)

	assert.Len(t, bAndA.bits, 1)
	assert.Len(t, bOrA.bits, 2)
}

func Test_Uint64_Inspect(t *testing.T) {
	a := NewUint64(2)
	a.Set(1, 2, 3, 1)
	a.Unset(2)

	assert.Equal(t, 64, a.Len())
	assert.GreaterOrEqual(t, a.Cap(), 64)
	assert.Equal(t, uint(2), a.Pop())
}
