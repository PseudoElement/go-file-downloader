package seabattle

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type EventHandlers struct {
	room  *Room
	rooms []*Room
}

func NewEventHandlers(room *Room, rooms []*Room) EventHandlers {
	return EventHandlers{room: room, rooms: rooms}
}

func (eh *EventHandlers) HandleNewMsg(msgBody SocketRequestMsg[any]) error {
	switch msgBody.ActionType {
	//  @TODO fix ERROR when send STEP action from client
	case STEP:
		var stepData NewStepReqMsg
		bytes, err := json.Marshal(msgBody.Data)
		log.Println(err)
		json.Unmarshal(bytes, &stepData)
		return eh.handlePlayerStep(msgBody.Email, stepData)
	case SET_PLAYER_POSITIONS:
		var setPositionsData PlayerPositionsMsg
		bytes, err := json.Marshal(msgBody.Data)
		log.Println(err)
		json.Unmarshal(bytes, &setPositionsData)
		return eh.handlePlayerSetPositions(msgBody.Email, setPositionsData.PlayerPositions)
	default:
		fmt.Errorf("Unknown msgBody type.")
		eh.sendMessageToClients(struct {
			Message string `json:"message"`
		}{Message: "Unknown msgBody type."})
		return nil
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
		if player.Conn() != nil {
			if err := player.Conn().WriteJSON(msg); err != nil {
				eh.queries().SaveNewError(player.room.id, err.Error())
			}
		}
	}
}

func (eh *EventHandlers) sendMessageToEnemy(enemy *Player, msg any) {
	if enemy != nil && enemy.Conn() != nil {
		if err := enemy.Conn().WriteJSON(msg); err != nil {
			eh.queries().SaveNewError(enemy.room.id, err.Error())
		}
	}
}

func (eh *EventHandlers) handleConnection(email string) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)

	var yourData PlayerInfoForClientOnConnection
	var enemyData PlayerInfoForClientOnConnection
	if player != nil {
		yourData = PlayerInfoForClientOnConnection{
			PlayerId:    player.info.id,
			PlayerEmail: player.info.email,
			IsOwner:     player.info.isOwner,
		}
	}
	if enemy != nil {
		enemyData = PlayerInfoForClientOnConnection{
			PlayerId:    enemy.info.id,
			PlayerEmail: enemy.info.email,
			IsOwner:     enemy.info.isOwner,
		}
	}

	msg := SocketRespMsg[ConnectPlayerResp]{
		Message:    fmt.Sprintf("Player %s connected to room %s.", player.info.email, player.room.name),
		ActionType: CONNECT_PLAYER,
		Data: ConnectPlayerResp{
			RoomId:    eh.room.id,
			RoomName:  eh.room.name,
			CreatedAt: eh.room.created_at,
			YourData:  yourData,
			EnemyData: enemyData,
		},
	}
	eh.sendMessageToEnemy(enemy, msg)

	return nil
}

func (eh *EventHandlers) handleDisconnection(email string) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)
	wasOwner := player.info.isOwner

	if err := eh.queries().DisconnectPlayerFromRoom(player.info.email, player.room.name); err != nil {
		return err
	}
	delete(eh.room.players, player.info.id)

	if wasOwner && enemy != nil {
		if err := eh.queries().ChangeOwnerStatus(enemy.info.id, true); err != nil {
			eh.queries().SaveNewError(eh.room.id, err.Error())
		}
		enemy.MakeAsOwner()
	}

	if isEmptyRoom := len(eh.room.players) == 0; isEmptyRoom {
		if err := eh.queries().DeleteRoom(eh.room.id); err != nil {
			eh.queries().SaveNewError(eh.room.id, err.Error())
		}

		eh.rooms = slice_utils_module.Filter(eh.rooms, func(r *Room, idx int) bool {
			return r.id != eh.room.id
		})
	}

	msg := SocketRespMsg[DisconnectPlayerResp]{
		Message:    fmt.Sprintf("Player %s disconnected from %s.", player.info.email, player.room.name),
		ActionType: DISCONNECT_PLAYER,
		Data: DisconnectPlayerResp{
			RoomId:   player.room.id,
			RoomName: player.room.name,
			Email:    player.info.email,
			Id:       player.info.id,
		},
	}
	eh.sendMessageToClients(msg)

	return nil
}

func (eh *EventHandlers) handlePlayerSetPositions(email string, playerPositions string) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)

	player.positions = playerPositions
	var enemyPositions string
	if enemy != nil {
		enemyPositions = enemy.positions
	}

	allPositions := player.info.id + ": " + playerPositions + PLAYER_POSITIONS_SEPARATOR + enemy.info.id + ": " + enemyPositions

	if err := eh.queries().UpdatePositions(allPositions, eh.room.name); err != nil {
		eh.queries().SaveNewError(player.room.id, err.Error())
	}

	msg := SocketRespMsg[PlayerSetPositionsResp]{
		Message:    fmt.Sprintf("Player %s set positions.", player.info.email),
		ActionType: SET_PLAYER_POSITIONS,
		Data: PlayerSetPositionsResp{
			Email: player.info.email,
			Id:    player.info.id,
		},
	}
	for _, player := range eh.room.players {
		if err := player.Conn().WriteJSON(msg); err != nil {
			eh.queries().SaveNewError(player.room.id, err.Error())
		}
	}

	return nil
}

func (eh *EventHandlers) handlePlayerStep(email string, step NewStepReqMsg) error {
	player := eh.getPlayerByEmail(email)
	enemy := eh.getEnemy(email)
	// FIX enemy is nil
	expression := fmt.Sprintf("%s[^,]*,", step.Step)
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
		Message:    fmt.Sprintf("Player %s already checked this cell. Next time check another one!", player.info.email),
		ActionType: STEP,
		Data: PlayerStepResp{
			Email:  player.info.email,
			Id:     player.info.id,
			Step:   step.Step,
			Result: ALREADY_CHECKED,
		},
	}
	if err := player.Conn().WriteJSON(steppingPlayerMsg); err != nil {
		eh.queries().SaveNewError(player.room.id, err.Error())
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
		winMsg := SocketRespMsg[GameWinResp]{
			Message:    fmt.Sprintf("Player %s won the game.", player.info.email),
			ActionType: WIN_GAME,
			Data: GameWinResp{
				WinnerEmail: player.info.email,
				WinnerId:    player.info.id,
			},
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
