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

func (q SeaBattleQueries) GetRoomsList() ([]DB_PlayerWithRoomJoinRow, error) {
	roomsData := make([]DB_PlayerWithRoomJoinRow, 0, 1000)
	rows, err := q.db.Query(`
		SELECT r.id, r.name, r.positions, r.created_at, p.email, p.id, p.is_owner
		FROM seabattle_rooms r 
		LEFT JOIN seabattle_players p
		ON r.room_name = p.room_name;
	`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dbRow := new(DB_PlayerWithRoomJoinRow)
		if err := rows.Scan(
			&dbRow.RoomId,
			&dbRow.RoomName,
			&dbRow.RoomPositions,
			&dbRow.CreatedAt,
			&dbRow.PlayerEmail,
			&dbRow.PlayerId,
			&dbRow.IsOwner,
		); err != nil {
			if err == sql.ErrNoRows {
				return roomsData, nil
			}
			return nil, err
		}
	}

	return roomsData, nil
}

func (q SeaBattleQueries) CreateRoom(roomName string) (DB_NewCreatedRoom, error) {
	var newRoom DB_NewCreatedRoom
	query := fmt.Sprintf(`
		INSERT INTO seabattle_rooms(room_name) 
		VALUES($1)
		RETURNING id, room_name, created_at;
	`)
	err := q.db.QueryRow(query, roomName).Scan(&newRoom.RoomId, &newRoom.CreatedAt, &newRoom.RoomName)
	if err != nil {
		return newRoom, fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return newRoom, nil
}

func (q SeaBattleQueries) DeleteRoom(roomName string) error {
	query := fmt.Sprintf("DELETE FROM seabattle_rooms(room_name) VALUES($1);")
	_, err := q.db.Exec(query, roomName)
	if err != nil {
		return fmt.Errorf("Error in DeleteRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) CheckRoomAlreadyExists(roomName string) (bool, error) {
	var room any
	row := q.db.QueryRow("SELECT * FROM seabattle_rooms WHERE room_name=$1;", roomName)

	if err := row.Scan(&room); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return true, err
		}
	}

	return true, nil
}

func (q SeaBattleQueries) ConnectPlayerToRoom(email string, roomName string, isOwner bool) error {
	query := fmt.Sprintf("INSERT INTO seabattle_players(email, room_name, is_owner) VALUES($1, $2, $3);")
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
	query := fmt.Sprintf("INSERT INTO seabattle_rooms(positions) VALUES($1) WHERE =room_name$2;")
	_, err := q.db.Exec(query, newPositions, roomName)
	if err != nil {
		return fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return nil
}
