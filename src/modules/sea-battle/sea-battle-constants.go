package seabattle

import (
	"strconv"
)

const PLAYER_FIELDS_SEPARATOR = "___"

func CreateMockFields(player1 string, player2 string) string {
	s := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	creator := func(player string) string {
		str := player + ": "
		letterIdx := 0
		for count := 1; count < 11; count++ {
			letter := s[letterIdx]
			countStr := strconv.Itoa(count)
			str += letter + countStr + ","
			if count == 10 {
				if letter == s[len(s)-1] {
					break
				}
				count = 0
				letterIdx++
			}
		}
		return str[:len(str)-1]
	}

	fieldsPlayer1 := creator(player1)
	fieldsPlayer2 := creator(player2)

	allFields := fieldsPlayer1 + PLAYER_FIELDS_SEPARATOR + fieldsPlayer2

	return allFields
}
