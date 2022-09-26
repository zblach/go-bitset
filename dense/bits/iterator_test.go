package bits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate_Uint8(t *testing.T) {
	s := NewUint8(40)     // length = 5
	s.Set(1, 2, 4, 6)     // word[0]
	s.Set(20, 21, 22, 23) // word[2]

	it := s.Iterate()

	i, values := 0, make([]uint, s.Pop())
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		values[i] = n
		i++
	}

	assert.EqualValues(t, []uint{1, 2, 4, 6, 20, 21, 22, 23}, values)
}

func TestIterate_Uint64(t *testing.T) {
	s := NewUint64(40)    // length = 1
	s.Set(1, 2, 4, 6)     // word[0]
	s.Set(20, 21, 22, 23) // word[0]

	it := s.Iterate()

	i, values := 0, make([]uint, s.Pop())
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		values[i] = n
		i++
	}

	assert.EqualValues(t, []uint{1, 2, 4, 6, 20, 21, 22, 23}, values)
}
