package seabattle

import (
	"strconv"
)

// ERROR CODES
const (
	ROOM_ALREADY_EXISTS = 1
)

const PLAYER_POSITIONS_SEPARATOR = "___"
const STRIKED_CELL_SYMBOL = "*"
const CELL_WITH_SHIP_SYMBOL = "+"
const TOTAL_CELL_WITH_SHIPS_COUNT = 20

const (
	// FOR ALL
	READY                = "READY"
	CONNECT_PLAYER       = "CONNECT_PLAYER"
	DISCONNECT_PLAYER    = "DISCONNECT_PLAYER"
	STEP                 = "STEP"
	SET_PLAYER_POSITIONS = "SET_PLAYER_POSITIONS"
	WIN_GAME             = "WIN_GAME"
	// FOR SINGLE USER
	ERROR = "ERROR"
)

// Player step results
const (
	HIT             = "HIT"
	MISS            = "MISS"
	KILL            = "KILL"
	ALREADY_CHECKED = "ALREADY_CHECKED"
)

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

	allFields := fieldsPlayer1 + PLAYER_POSITIONS_SEPARATOR + fieldsPlayer2

	return allFields
}
