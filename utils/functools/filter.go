package functools

func Filter[T any](ts []T, fn func(T) bool) []T {
	var filtered []T
	for _, t := range ts {
		if fn(t) {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
