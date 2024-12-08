package seabattle_queries

import (
	"database/sql"
	"fmt"
)

type SeaBattleQueries struct {
	db *sql.DB
}

func New(db *sql.DB) SeaBattleQueries {
	return SeaBattleQueries{
		db: db,
	}
}

func (q SeaBattleQueries) createRoomsTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_rooms (
			id SERIAL NOT NULL PRIMARY KEY, 
            name TEXT,
			positions TEXT,
            CONSTRAINT uc_room UNIQUE (id, name)
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in CreateRoomsTable. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) createPlayersTable() error {
	_, err := q.db.Exec(`
		CREATE TABLE IF NOT EXISTS seabattle_players (
			id SERIAL NOT NULL PRIMARY KEY, 
            email TEXT,
			room_name TEXT,
            CONSTRAINT uc_player UNIQUE (id, email, room_name)
		);
	`)
	if err != nil {
		return fmt.Errorf("Error in CreatePlayersTable. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) CreateRoom(roomName string) error {
	query := fmt.Sprintf("INSERT INTO seabattle_rooms(room_name) VALUES($1);")
	_, err := q.db.Exec(query, roomName)
	if err != nil {
		return fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) DeleteRoom(roomName string) error {
	query := fmt.Sprintf("DELETE FROM seabattle_rooms(room_name) VALUES($1);")
	_, err := q.db.Exec(query, roomName)
	if err != nil {
		return fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) ConnectPlayerToRoom(email string, roomName string) error {
	query := fmt.Sprintf("INSERT INTO seabattle_players(email, room_name) VALUES($1, $2);")
	_, err := q.db.Exec(query, email, roomName)
	if err != nil {
		return fmt.Errorf("Error in ConnectPlayerToRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) DisconnectPlayerFromRoom(email string, roomName string) error {
	query := fmt.Sprintf("DELETE FROM seabattle_players WHERE email=$1 AND room_name=$2;")
	_, err := q.db.Exec(query, email, roomName)
	if err != nil {
		return fmt.Errorf("Error in DisconnectPlayerFromRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) UpdatePositions(newPositions string, roomName string) error {
	query := fmt.Sprintf("INSERT INTO seabattle_rooms(positions) VALUES($1) WHERE room_name=$2;")
	_, err := q.db.Exec(query, newPositions, roomName)
	if err != nil {
		return fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return nil
}
