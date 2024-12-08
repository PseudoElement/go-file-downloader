package seabattle_queries

import "fmt"

func (q SeaBattleQueries) createErrorsTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_errors (
			id SERIAL NOT NULL PRIMARY KEY, 
            error VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in CreateRoomsTable. Error: %s", err.Error())
	}

	return nil
}

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
