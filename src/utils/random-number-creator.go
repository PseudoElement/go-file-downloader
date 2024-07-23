package custom_utils

import "math/rand"

func CreateRandomNumber(min int, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
