package custom_utils

import (
	"math"
	"math/rand"
	"strings"
)

const (
	possibleCons   = "bcdfghjklmnpqrstvwxyz"
	possibleVowels = "aeiou"
)

func CreateRandowWordForSqlTable(minLength int, maxLength int, startUpperCase bool) string {
	return "'" + CreateRandomWord(minLength, maxLength, startUpperCase) + "'"
}

func CreateRandomWord(minLength int, maxLength int, startUpperCase bool) string {
	word := ""
	randomLength := CreateRandomNumber(minLength, maxLength)

	for ind := range randomLength {
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
