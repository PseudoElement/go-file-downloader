package common

import (
	"math/rand"
	"slices"
	"time"
)

func RandomString() string {
	rand.NewSource(time.Now().UnixNano())
	length := 36
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	var shouldAddDash = func(idx int) bool {
		dashIdxs := []int{8, 13, 18, 23}
		return slices.Contains(dashIdxs, idx)
	}
	for i := 0; i < length; i++ {
		if shouldAddDash(i) {
			result[i] = '-'
			i++
		}
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}
