package bitset

type Bitset interface {
	Get(index uint) bool

	Set(indices ...uint) Bitset
	Unset(indices ...uint) Bitset

	Clear()
}

func Copy(dst Bitset, src Inspect) {
	dst.Clear()
	it := src.Iterate()
	i, vals := 0, make([]uint, src.Pop())
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		vals[i] = n
		i++
	}
	dst.Set(vals...)
}

type Logical[T Bitset] interface {
	And(b T) (aAndB T)
	Or(b T) (aOrB T)
}

type Size interface {
	Len() int
	Cap() int
}

// https://github.com/golang/go/discussions/54245
type Iterator interface {
	Next() (uint, bool)
}

type Inspect interface {
	Size
	Pop() uint
	Iterate() Iterator
}
