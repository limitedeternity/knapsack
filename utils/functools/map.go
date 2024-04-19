package functools

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, 0, len(ts))
	for _, t := range ts {
		result = append(result, fn(t))
	}

	return result
}
