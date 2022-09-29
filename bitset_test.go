package bitset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zblach/go-bitset"
	"github.com/zblach/go-bitset/dense/bits"
	"github.com/zblach/go-bitset/dense/bools"
)

func TestHeterogenousCopy(t *testing.T) {
	src := bools.New[uint](30)
	dst := bits.NewUint(0)

	src.Set(1, 2, 3, 4, 10, 11, 13, 15)

	bitset.Copy[uint](dst, src)

	assert.Equal(t, dst.Pop(), uint(8))

}
