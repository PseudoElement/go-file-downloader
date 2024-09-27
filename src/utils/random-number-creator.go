package custom_utils

import "math/rand"

func CreateRandomNumber(min int64, max int64) int64 {
	if min == max {
		return min
	}
	return int64(rand.Intn(int(max-min))) + min
}
