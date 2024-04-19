package functools

type Pair[T, U any] struct {
	First  T
	Second U
}

func Zip[T, U any](ts []T, us []U) []Pair[T, U] {
	if len(ts) != len(us) {
		panic("slices have different length")
	}

	pairs := make([]Pair[T, U], 0, len(ts))
	for i := range ts {
		pairs = append(pairs, Pair[T, U]{ts[i], us[i]})
	}

	return pairs
}
