package seabattle

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

func GetPlayersFromRoom(playerEmail string, room *Room) (RoomPlayers, bool) {
	playersOfRoom := RoomPlayers{}
	if len(room.players) == 0 {
		return playersOfRoom, true
	}

	for _, player := range room.players {
		if player.info.email == playerEmail {
			playersOfRoom.CurrentPlayer = player
		} else {
			playersOfRoom.Enemy = player
		}
	}
	fmt.Println("playersOfRoom ==> ", playersOfRoom)

	return playersOfRoom, false
}

func MockHttpReq() *http.Request {
	return new(http.Request)
}

func MockRespWriter() http.ResponseWriter {
	return struct{ http.ResponseWriter }{}
}

func MockPlayer() *Player {
	return &Player{
		info:      PlayerInfo{},
		positions: "",
		room:      new(Room),
		rooms:     []*Room{},
	}
}

func IsPlayerAlreadyAddedToRoomFromDB(room *Room, email string) bool {
	for _, player := range room.players {
		if player.info.email == email {
			return true
		}
	}
	return false
}

func IsPlayerAlreadyConnectedToRoom(room *Room, email string) bool {
	for _, player := range room.players {
		if player.info.email == email && player.Conn() != nil {
			return true
		}
	}
	return false
}

func IsWin(steppingPlayerEmail string, room *Room) bool {
	enemy := GetEnemyInRoom(steppingPlayerEmail, room)

	var strikedEnemyCellsCount int8
	for _, runeChar := range enemy.positions {
		if string(runeChar) == STRIKED_CELL_SYMBOL {
			strikedEnemyCellsCount++
		}
	}

	return strikedEnemyCellsCount >= TOTAL_CELL_WITH_SHIPS_COUNT
}

func IsShipKilled(enemy *Player, step NewStepReqMsg) bool {
	row := step.Step[:1]    // A
	column := step.Step[1:] // 1

	to1 := common.ToInt(column) - 1
	to10 := common.ToInt(column) + 1
	// prevCell was empty
	var isLeftShipEnd bool
	// nextCell was empty
	var isRightShipEnd bool
	//ROW A1, A2, A3...
	for idx := 0; to10 < 10 || to1 > 0; idx++ {
		prevCellInRow := row + strconv.Itoa(to1)
		nextCellInRow := row + strconv.Itoa(to10)

		prevCellExpression := fmt.Sprintf("%s[^,]*,", prevCellInRow)
		r1, _ := regexp.Compile(prevCellExpression)
		prevCellValue := r1.FindString(enemy.positions)

		nextCellExpression := fmt.Sprintf("%s[^,]*,", nextCellInRow)
		r2, _ := regexp.Compile(nextCellExpression)
		nextCellValue := r2.FindString(enemy.positions)

		if !strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) {
			isLeftShipEnd = true
		}
		if !strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) {
			isRightShipEnd = true
		}

		if isLeftShipEnd && isRightShipEnd {
			if idx == 0 {
				break
			} else {
				return true
			}
		}

		matchPrev := strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(prevCellValue, STRIKED_CELL_SYMBOL)
		matchNext := strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(nextCellValue, STRIKED_CELL_SYMBOL)

		if !matchNext || !matchPrev {
			return false
		}

		if to1 > 0 {
			to1--
		}
		if to10 < 10 {
			to10++
		}
	}

	//column (A1, B1, C1...)
	letters := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	minLettersIdx := 0
	maxLettersIdx := len(letters) - 1
	var columnIdx int
	for i, l := range letters {
		if string(l) == string(row) {
			columnIdx = i
		}
	}

	toJ := columnIdx + 1
	toA := columnIdx - 1
	// nextCell was empty
	var isBottomShipEnd bool
	// prevCell was empty
	var isTopShipEnd bool
	for idx := 0; toJ < 10 || toA > 0; idx++ {
		isPrevOutOfMap := toA < minLettersIdx
		isNextOutOfMap := toJ > maxLettersIdx
		// A1+*,A2,A3
		var matchPrev bool = isPrevOutOfMap
		var matchNext bool = isNextOutOfMap
		prevCellValue := ""
		nextCellValue := ""
		if !isPrevOutOfMap {
			prevLetter := letters[toA]
			prevCellInColumn := prevLetter + column

			prevCellExpression := fmt.Sprintf("%s[^,]*,", prevCellInColumn)
			r1, _ := regexp.Compile(prevCellExpression)
			prevCellValue = r1.FindString(enemy.positions)

			matchPrev = strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(prevCellValue, STRIKED_CELL_SYMBOL)

		}
		if !isNextOutOfMap {
			nextLetter := letters[toJ]
			nextCellInColumn := nextLetter + column

			nextCellExpression := fmt.Sprintf("%s[^,]*,", nextCellInColumn)
			r2, _ := regexp.Compile(nextCellExpression)
			nextCellValue = r2.FindString(enemy.positions)

			matchNext = strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(nextCellValue, STRIKED_CELL_SYMBOL)
		}

		if !strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) {
			isTopShipEnd = true
		}
		if !strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) {
			isBottomShipEnd = true
		}

		if (isTopShipEnd || isPrevOutOfMap) && (isBottomShipEnd || isNextOutOfMap) {
			return true
		}

		if !matchNext || !matchPrev {
			return false
		}

		if toA > -1 {
			toA--
		}
		if toJ < 11 {
			toJ++
		}
	}

	return false
}

func GetPlayerInRoomByEmail(steppingPlayerEmail string, room *Room) *Player {
	for _, player := range room.players {
		if player.info.email == steppingPlayerEmail {
			return player
		}
	}

	return nil
}

func GetEnemyInRoom(steppingPlayerEmail string, room *Room) *Player {
	for _, pl := range room.players {
		if pl.info.email != steppingPlayerEmail {
			return pl
		}
	}

	return nil
}
