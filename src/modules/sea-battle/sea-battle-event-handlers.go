package seabattle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
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
		bytes, _ := json.Marshal(msgBody.Data)
		json.Unmarshal(bytes, &stepData)
		return eh.handlePlayerStep(msgBody.Email, stepData)
	case SET_PLAYER_POSITIONS:
		var setPositionsData PlayerPositionsMsg
		bytes, _ := json.Marshal(msgBody.Data)
		json.Unmarshal(bytes, &setPositionsData)
		return eh.handlePlayerSetPositions(msgBody.Email, setPositionsData.PlayerPositions)
	case READY:
		return eh.handlePlayerReady(msgBody.Email)
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
	player := GetPlayerInRoomByEmail(email, eh.room)
	enemy := GetEnemyInRoom(email, eh.room)

	if player == nil {
		player = MockPlayer()
	}
	if enemy == nil {
		enemy = MockPlayer()
	}

	msg := SocketRespMsg[ConnectPlayerResp]{
		Message:    fmt.Sprintf("Player %s connected to room %s.", player.info.email, player.room.name),
		ActionType: CONNECT_PLAYER,
		Data: ConnectPlayerResp{
			RoomId:    eh.room.id,
			RoomName:  eh.room.name,
			CreatedAt: eh.room.created_at,
			YourData: PlayerInfoForClientOnConnection{
				PlayerId:    player.info.id,
				PlayerEmail: player.info.email,
				IsOwner:     player.info.isOwner,
			},
			EnemyData: PlayerInfoForClientOnConnection{
				PlayerId:    enemy.info.id,
				PlayerEmail: enemy.info.email,
				IsOwner:     enemy.info.isOwner,
			},
		},
	}
	eh.sendMessageToEnemy(enemy, msg)

	return nil
}

func (eh *EventHandlers) handleDisconnection(email string) error {
	player := GetPlayerInRoomByEmail(email, eh.room)
	enemy := GetEnemyInRoom(email, eh.room)
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

func (eh *EventHandlers) handlePlayerReady(email string) error {
	player := GetPlayerInRoomByEmail(email, eh.room)
	player.setReadyStatus(true)

	msg := SocketRespMsg[PlayerReadyResp]{
		Message:    fmt.Sprintf("Player %s is ready.", email),
		ActionType: READY,
		Data: PlayerReadyResp{
			Email: player.info.email,
			Id:    player.info.id,
		},
	}
	eh.sendMessageToClients(msg)

	return nil
}

func (eh *EventHandlers) handlePlayerSetPositions(email string, playerPositions string) error {
	player := GetPlayerInRoomByEmail(email, eh.room)
	enemy := GetEnemyInRoom(email, eh.room)

	player.positions = playerPositions
	if enemy == nil {
		enemy = MockPlayer()
	}

	allPositions := player.info.id + ": " + playerPositions + PLAYER_POSITIONS_SEPARATOR + enemy.info.id + ": " + enemy.positions

	if err := eh.queries().UpdatePositions(allPositions, eh.room.id); err != nil {
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

	eh.sendMessageToClients(msg)

	return nil
}

func (eh *EventHandlers) handlePlayerStep(email string, step NewStepReqMsg) error {
	player := GetPlayerInRoomByEmail(email, eh.room)
	enemy := GetEnemyInRoom(email, eh.room)
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

	if IsShipKilled(enemy, step) {
		msg.Message = strings.Replace(msg.Message, "hit", "killed", 1)
		msg.Data.Result = KILL
	}

	eh.sendMessageToClients(msg)

	if IsWin(player.info.email, eh.room) {
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
	if err := eh.queries().UpdatePositions(newAllPositions, eh.room.id); err != nil {
		return err
	}

	return nil
}
