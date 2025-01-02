package seabattle_queries

import (
	"fmt"

	db_interfaces "github.com/pseudoelement/go-file-downloader/src/db/db-interfaces"
)

func (q SeaBattleQueries) CreateTables() error {
	if err := q.createErrorsTable(); err != nil {
		return err
	}
	if err := q.createRoomsTable(); err != nil {
		return err
	}
	if err := q.createPlayersTable(); err != nil {
		return err
	}

	return nil
}

func (q SeaBattleQueries) createErrorsTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_errors (
			id SERIAL NOT NULL PRIMARY KEY, 
            error VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in createErrorsTable. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) createRoomsTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_rooms (
			id SERIAL NOT NULL PRIMARY KEY, 
            room_name TEXT,
			positions TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
            CONSTRAINT uc_room UNIQUE (id, room_name)
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in createRoomsTable. Error: %s", err.Error())
	}

	return nil
}

/*
resp: {
	rooms: {
		room_id_1: {room_name: "room_1", room_id: "room_id_1", created_at: "12 Dec 2025", players: [Player, Player]},
		room_id_2: {room_name: "room_2", room_id: "room_id_2", created_at: "13 Dec 2025", players: [Player, Player]}
	}
}
*/

func (q SeaBattleQueries) createPlayersTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_players (
			id SERIAL NOT NULL PRIMARY KEY, 
            email TEXT,
			room_name TEXT,
			is_owner BOOLEAN NOT NULL,
            CONSTRAINT uc_player UNIQUE (id, email, room_name)
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in CreatePlayersTable. Error: %s", err.Error())
	}

	return nil
}

var _ db_interfaces.TableCreator = (*SeaBattleQueries)(nil)
