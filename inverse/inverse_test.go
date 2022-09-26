package inverse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zblach/bitset/dense/bits"
)

func Test_Inverse(t *testing.T) {
	n := New(bits.NewUint8(8).Set(1, 2, 3))

	assert.False(t, n.Get(1))
	assert.True(t, n.Get(0))
}
