package custom_utils

import (
	"math/rand"
)

func CreateRandomBool() bool {
	randByte := rand.Intn(2)
	if randByte == 1 {
		return true
	} else {
		return false
	}
}
