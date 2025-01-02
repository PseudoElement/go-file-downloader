package seabattle_queries

import "fmt"

func (q SeaBattleQueries) SaveNewError(roomId string, player_id string, errorMsg string) {
	info := fmt.Sprintf(
		"Error in handleConnection: room_id - %s, player_id - %s, error - %s",
		roomId,
		player_id,
		errorMsg,
	)

	_, err := q.db.Exec(`INSERT INTO seabattle_errors(error) VALUES($1);`, info)
	if err != nil {
		fmt.Errorf("Error in SaveNewError. Error: %s", err.Error())
	}
}
