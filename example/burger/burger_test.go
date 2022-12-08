package burger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zblach/go-bitset/dense/bools"
)

type topping uint

const (
	_ topping = iota

	LETTUCE
	TOMATO
	ONION
	CHEESE
	MORE_CHEESE

	MUSTARD
	KETCHUP
	DONKEY_SAUCE
	BBQ_SAUCE

	BACON
	MORE_BACON
	EVEN_MORE_BACON
	EXTRA_BACON
	DOUBLE_EXTRA_BACON

	topping_count uint = iota - 1
)

type burger struct{ bools.Bitset[topping] }

func New() burger {
	return burger{*bools.New[topping](topping_count)}
}

var (
	base      = New()
	classic   = base.With(LETTUCE, TOMATO, MUSTARD)
	royale    = classic.With(CHEESE, ONION, DONKEY_SAUCE)
	baconizer = classic.With(BACON, MORE_BACON)
)

func (b burger) hasBacon() bool {
	return b.Any(BACON, MORE_BACON, EVEN_MORE_BACON, EXTRA_BACON, DOUBLE_EXTRA_BACON)
}

func (b burger) With(t ...topping) burger {
	new := b.Copy()
	new.Set(t...)
	return burger{*new}
}

func Test_Instantiation(t *testing.T) {
	classic := New().With(LETTUCE, TOMATO, MUSTARD)

	t.Log("I'll have uh.....")

	want := classic.With(BACON)
	t.Log(want)

	t.Log("but hold the tomato.")
	want.Unset(TOMATO)

	assert.False(t, want.Get(TOMATO))
	assert.True(t, want.Get(BACON))
}
