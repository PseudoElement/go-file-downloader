package custom_utils

import (
	"math/rand"
)

func CreateRandomByteForSql() byte {
	if CreateRandomBool() {
		return '1'
	}
	return '0'
}

func CreateRandomBool() bool {
	randByte := rand.Intn(2)
	if randByte == 1 {
		return true
	}
	return false
}
