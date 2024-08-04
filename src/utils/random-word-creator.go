package custom_utils

import (
	"math"
	"math/rand"
	"strings"
)

const (
	possibleLetters = "bcdfghjklmnpqrstvwxyzaeiou"
)

func CreateRandomWord(minLength int, maxLength int, startUpperCase bool) string {
	var str = ""
	randomLength := CreateRandomNumber(minLength, maxLength)

	for ind := range randomLength {
		randInd := int(math.Max(0, float64(rand.Intn(len(possibleLetters))-1)))
		randomLetter := string(possibleLetters[randInd])
		if startUpperCase && ind == 0 {
			randomLetter = strings.ToUpper(randomLetter)
		}
		str += randomLetter
	}

	return str
}
