package slice_utils

func Map[T any, R any](slice []T, fn func(el T, index int) R) []R {
	sl := []R{}

	for i, el := range slice {
		value := fn(el, i)
		sl = append(sl, value)
	}

	return sl
}
