package slice_utils

func Map[T any, R any](slice []T, fn func(el T, index int) R) []R {
	sl := []R{}

	for i, el := range slice {
		value := fn(el, i)
		sl = append(sl, value)
	}

	return sl
}

func IndexOf[T comparable](slice []T, val T) int {
	for idx, elasd := range slice {
		if elasd == val {
			return idx
		}
	}

	return -1
}
