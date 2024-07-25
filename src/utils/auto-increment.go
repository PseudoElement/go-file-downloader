package custom_utils

func AutoIncrement(start int) func() int {
	count := start
	return func() int {
		if start == count {
			count++
			return start
		}
		value := count
		count++
		return value
	}
}
