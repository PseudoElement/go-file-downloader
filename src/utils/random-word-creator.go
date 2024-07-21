package utils

import (
	"math"
	"math/rand"
	"strings"
)

const (
	possibleCons   = "bcdfghjklmnpqrstvwxyz"
	possibleVowels = "aeiou"
)

func CreateRandomWord(length int, startUpperCase bool) string {
	word := ""

	for ind := range length {
		isVowel := rand.Intn(100) > 50

		var letters string
		if isVowel {
			letters = possibleVowels
		} else {
			letters = possibleCons
		}

		randInd := int(math.Max(0, float64(rand.Intn(len(letters))-1)))
		randomLetter := string(letters[randInd])
		if startUpperCase && ind == 0 {
			randomLetter = strings.ToUpper(randomLetter)
		}

		word += randomLetter
	}

	return word
}
