package seabattle

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type EventHandlers struct {
	room Room
}

func (eh *EventHandlers) HandleNewMsg(msgBody SocketRequestMsg[any]) error {
	switch data := msgBody.Data.(type) {
	case NewStepReqMsg:
		return eh.handlePlayerStep(msgBody.Email, data)
	case PlayerPositionsMsg:
		return eh.handlePlayerSetPositions(msgBody.Email, data.PlayerPositions)
	default:
		return fmt.Errorf("Unknown msgBody type.")
	}
}

func (eh *EventHandlers) queries() seabattle_queries.SeaBattleQueries {
	return eh.room.queries
}

func (eh *EventHandlers) getPlayerByEmail(steppingPlayerEmail string) *Player {
	for _, player := range eh.room.players {
		if player.info.email == steppingPlayerEmail {
			return player
		}
	}

	return nil
}

func (eh *EventHandlers) getEnemy(steppingPlayerEmail string) *Player {
	var enemy *Player
	for _, pl := range eh.room.players {
		if pl.info.email != steppingPlayerEmail {
			enemy = pl
		}
	}

	return enemy
}

func (eh *EventHandlers) sendMessageToClients(msg any) {
	for _, player := range eh.room.players {
		if err := player.Conn().WriteJSON(msg); err != nil {
			eh.queries().SaveNewError(player.room.id, player.info.id, err.Error())
		}
	}
}

func (eh *EventHandlers) handleConnection(email string) error {
	player := eh.getPlayerByEmail(email)
	if err := eh.queries().ConnectPlayerToRoom(player.info.email, player.room.name); err != nil {
		return err
	}

	for _, player := range eh.room.players {
		msg := SocketRespMsg[ConnectPlayerResp]{
			Message:    fmt.Sprintf("Player %s connected to %s.", player.info.email, player.room.name),
			ActionType: CONNECT_PLAYER,
			Data: ConnectPlayerResp{
				Email: player.info.email,
				Id:    player.info.id,
			},
		}
		if err := player.Conn().WriteJSON(msg); err != nil {
			eh.queries().SaveNewError(player.room.id, player.info.id, err.Error())
		}
	}

	return nil
}

func (eh *EventHandlers) handlePlayerSetPositions(email string, playerPositions string) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)

	player.positions = playerPositions
	allPositions := player.info.id + ": " + playerPositions + PLAYER_POSITIONS_SEPARATOR + enemy.info.id + ": " + enemy.positions

	if err := eh.queries().UpdatePositions(allPositions, eh.room.name); err != nil {
		eh.queries().SaveNewError(player.room.id, player.info.id, err.Error())
	}

	msg := SocketRespMsg[PlayerSetPositionsResp]{
		Message:    fmt.Sprintf("Player %s set positions.", player.info.email),
		ActionType: SET_PLAYER_POSITIONS,
	}
	for _, player := range eh.room.players {
		if err := player.Conn().WriteJSON(msg); err != nil {
			eh.queries().SaveNewError(player.room.id, player.info.id, err.Error())
		}
	}

	return nil
}

func (eh *EventHandlers) handlePlayerStep(email string, step NewStepReqMsg) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)

	expression := fmt.Sprintf("%s.*,", step.Step)
	r, _ := regexp.Compile(expression)
	selectedCellValue := r.FindString(enemy.positions)
	cellValueWithoutComma := selectedCellValue[:len(selectedCellValue)-1]

	if strings.Contains(cellValueWithoutComma, ".") || strings.Contains(cellValueWithoutComma, "*") {
		eh.handleAlreadyChecked(player, enemy, step)
	} else if strings.Contains(cellValueWithoutComma, "+") {
		eh.handleHit(player, enemy, step, cellValueWithoutComma)
	} else {
		eh.handleMiss(player, enemy, step, cellValueWithoutComma)
	}

	return nil
}

func (eh *EventHandlers) handleAlreadyChecked(player *Player, enemy *Player, step NewStepReqMsg) error {
	steppingPlayerMsg := SocketRespMsg[PlayerStepResp]{
		Message:    "You've already selected this cell. Try again!",
		ActionType: STEP,
		Data: PlayerStepResp{
			Email:  player.info.email,
			Id:     player.info.id,
			Step:   step.Step,
			Result: ALREADY_CHECKED,
		},
	}
	if err := player.Conn().WriteJSON(steppingPlayerMsg); err != nil {
		eh.queries().SaveNewError(player.room.id, player.info.id, err.Error())
	}

	return nil
}

func (eh *EventHandlers) handleHit(player *Player, enemy *Player, step NewStepReqMsg, cellValue string) error {
	if err := eh.updatePlayerPositions(player, enemy, step, cellValue, HIT); err != nil {
		return err
	}

	msg := SocketRespMsg[PlayerStepResp]{
		Message: fmt.Sprintf(
			"Player %s hit ship of player %s in cell %s.",
			player.info.email,
			enemy.info.email,
			step.Step,
		),
		ActionType: STEP,
		Data: PlayerStepResp{
			Email:  player.info.email,
			Id:     player.info.id,
			Step:   step.Step,
			Result: HIT,
		},
	}

	if eh.isShipKilled(enemy, step) {
		msg.Message = strings.Replace(msg.Message, "hit", "killed", 1)
		msg.Data.Result = KILL
	}

	eh.sendMessageToClients(msg)

	if eh.isWin(player.info.email) {
		winMsg := SocketRespMsg[any]{
			Message:    fmt.Sprintf("Player %s won the game.", player.info.email),
			ActionType: WIN_GAME,
		}
		eh.sendMessageToClients(winMsg)
	}

	return nil
}

func (eh *EventHandlers) handleMiss(player *Player, enemy *Player, step NewStepReqMsg, cellValue string) error {
	if err := eh.updatePlayerPositions(player, enemy, step, cellValue, MISS); err != nil {
		return err
	}

	msg := SocketRespMsg[PlayerStepResp]{
		Message:    fmt.Sprintf("Player %s missed in cell %s.", player.info.email, step.Step),
		ActionType: STEP,
		Data: PlayerStepResp{
			Email:  player.info.email,
			Id:     player.info.id,
			Step:   step.Step,
			Result: MISS,
		},
	}
	eh.sendMessageToClients(msg)

	return nil
}

func (eh *EventHandlers) updatePlayerPositions(player *Player, enemy *Player, step NewStepReqMsg, cellValue string, stepResult string) error {
	if stepResult == HIT || stepResult == KILL {
		enemy.positions = strings.Replace(enemy.positions, cellValue, fmt.Sprintf("%s*", cellValue), 1)
	} else {
		// MISS
		enemy.positions = strings.Replace(enemy.positions, cellValue, fmt.Sprintf("%s.", cellValue), 1)
	}

	newAllPositions := player.info.id + ": " + player.positions + PLAYER_POSITIONS_SEPARATOR + enemy.info.id + ": " + enemy.positions
	if err := eh.queries().UpdatePositions(newAllPositions, eh.room.name); err != nil {
		return err
	}

	return nil
}

func (eh *EventHandlers) isShipKilled(enemy *Player, step NewStepReqMsg) bool {
	column := step.Step[:1]
	row := step.Step[1:]

	to1 := common.ToInt(row) - 1
	to10 := common.ToInt(row) + 1
	for idx := 0; to10 <= 10 || to1 >= 0; idx++ {
		prevCellInColumn := column + strconv.Itoa(to1)
		nextCellInColumn := column + strconv.Itoa(to10)

		prevCellExpression := fmt.Sprintf("%s.*,", prevCellInColumn)
		r1, _ := regexp.Compile(prevCellExpression)
		prevCellValue := r1.FindString(enemy.positions)

		nextCellExpression := fmt.Sprintf("%s.*,", nextCellInColumn)
		r2, _ := regexp.Compile(nextCellExpression)
		nextCellValue := r2.FindString(enemy.positions)

		if !strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) && !strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) {
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

		if to1 > -1 {
			to1--
		}
		if to10 < 11 {
			to10++
		}
	}

	//HANDLE row (A1, B1, C1...)
	letters := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	var columnIdx int
	for i, l := range letters {
		if string(l) == string(column) {
			columnIdx = i + 1
		}
	}

	toJ := columnIdx
	toA := columnIdx
	for idx := 0; toJ <= 10 || toA >= 0; idx++ {
		prevLetter := letters[toA]
		nextLetter := letters[toJ]
		prevCellInRow := prevLetter + row
		nextCellInRow := nextLetter + row

		prevCellExpression := fmt.Sprintf("%s.*,", prevCellInRow)
		r1, _ := regexp.Compile(prevCellExpression)
		prevCellValue := r1.FindString(enemy.positions)

		nextCellExpression := fmt.Sprintf("%s.*,", nextCellInRow)
		r2, _ := regexp.Compile(nextCellExpression)
		nextCellValue := r2.FindString(enemy.positions)

		// here we suppose that ship located in one ROW (exm. A1+,B1+,C1+), cause it's not returned value in column loop.
		// if it's not contain ship even on first iteration - we suppose here is 1-cell ship.
		if !strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) && !strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) {
			return true
		}

		matchPrev := strings.Contains(prevCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(prevCellValue, STRIKED_CELL_SYMBOL)
		matchNext := strings.Contains(nextCellValue, CELL_WITH_SHIP_SYMBOL) == strings.Contains(nextCellValue, STRIKED_CELL_SYMBOL)

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

func (eh *EventHandlers) isWin(steppingPlayerEmail string) bool {
	enemy := eh.getEnemy(steppingPlayerEmail)

	var strikedEnemyCellsCount int8
	for _, runeChar := range enemy.positions {
		if string(runeChar) == STRIKED_CELL_SYMBOL {
			strikedEnemyCellsCount++
		}
	}

	return strikedEnemyCellsCount >= TOTAL_CELL_WITH_SHIPS_COUNT
}
