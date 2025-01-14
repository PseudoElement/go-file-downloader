package seabattle_queries

import "fmt"

func (q SeaBattleQueries) SaveNewError(roomId string, errorMsg string) {
	info := fmt.Sprintf(
		"Error in handleConnection: room_id - %s, error - %s",
		roomId,
		errorMsg,
	)

	_, err := q.db.Exec(`INSERT INTO seabattle_errors(error) VALUES($1);`, info)
	if err != nil {
		fmt.Errorf("Error in SaveNewError. Error: %s", err.Error())
	}
}
