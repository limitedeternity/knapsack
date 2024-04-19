package functools

func Reduce[T, M any](ts []T, fn func(M, T) M, init M) M {
	acc := init
	for _, t := range ts {
		acc = fn(acc, t)
	}

	return acc
}
