package seabattle_queries

import "fmt"

func (q SeaBattleQueries) SaveNewError(roomName string, errorMsg string) {
	info := fmt.Sprintf(
		"Error: room_id - %s, error - %s",
		roomName,
		errorMsg,
	)

	_, err := q.db.Exec(`INSERT INTO seabattle_errors(error) VALUES($1);`, info)
	if err != nil {
		fmt.Errorf("Error in SaveNewError. Error: %s", err.Error())
	}
}

func (q SeaBattleQueries) SaveAutoRoomDeletion(roomInfo DB_PlayerWithRoomJoinRow) {
	info := fmt.Sprintf("Room with ID %v deleted with data: player_email - %s, player_id - %s, room_name - %s.",
		roomInfo.RoomId,
		roomInfo.PlayerEmail,
		roomInfo.PlayerId,
		roomInfo.RoomName)
	_, err := q.db.Exec(`INSERT INTO seabattle_errors(error) VALUES($1);`, info)
	if err != nil {
		fmt.Errorf("Error in SaveAutoRoomDeletion. Error: %s", err.Error())
	}
}
