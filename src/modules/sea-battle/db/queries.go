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
		SELECT r.id, r.room_name, r.positions, r.created_at, p.email, p.id, p.is_owner
		FROM seabattle_rooms r 
		LEFT JOIN seabattle_players p
		ON r.room_name = p.room_name;
	`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dbRow := DB_PlayerWithRoomJoinRow{}
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

		roomsData = append(roomsData, dbRow)
	}

	return roomsData, nil
}

func (q SeaBattleQueries) CreateRoom(roomName string) (DB_NewCreatedRoom, error) {
	var newRoom DB_NewCreatedRoom
	query := fmt.Sprintf(`
		INSERT INTO seabattle_rooms(room_name, positions) 
		VALUES($1, $2)
		RETURNING id, room_name, created_at;
	`)
	err := q.db.QueryRow(query, roomName, "").Scan(&newRoom.RoomId, &newRoom.RoomName, &newRoom.CreatedAt)
	if err != nil {
		return newRoom, fmt.Errorf("Error in CreateRoom. Error: %s", err.Error())
	}

	return newRoom, nil
}

func (q SeaBattleQueries) DeleteRoom(roomId string) error {
	query := fmt.Sprintf("DELETE FROM seabattle_rooms WHERE id=$1;")
	_, err := q.db.Exec(query, roomId)
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

func (q SeaBattleQueries) CheckPlayerAlreadyExists(playerEmail string) (bool, error) {
	var player any
	row := q.db.QueryRow("SELECT * FROM seabattle_player WHERE email=$1;", playerEmail)

	if err := row.Scan(&player); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return true, err
		}
	}

	return true, nil
}

func (q SeaBattleQueries) ConnectPlayerToRoom(email string, roomName string, isOwner bool) (string, error) {
	var playerId string
	query := fmt.Sprintf(`
		INSERT INTO seabattle_players(email, room_name, is_owner) 
		VALUES($1, $2, $3)
		RETURNING id;
	`)

	err := q.db.QueryRow(query, email, roomName, isOwner).Scan(&playerId)
	if err != nil {
		return "", fmt.Errorf("Error in ConnectPlayerToRoom. Error: %s", err.Error())
	}

	return playerId, nil
}

func (q SeaBattleQueries) DisconnectPlayerFromRoom(email string, roomName string) error {
	query := fmt.Sprintf("DELETE FROM seabattle_players WHERE email=$1 AND room_name=$2;")
	_, err := q.db.Exec(query, email, roomName)
	if err != nil {
		return fmt.Errorf("Error in DisconnectPlayerFromRoom. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) UpdatePositions(newPositions string, roomId string) error {
	query := fmt.Sprintf(`
		UPDATE seabattle_rooms
		SET positions=$1 
		WHERE id=$2;
	`)
	_, err := q.db.Exec(query, newPositions, roomId)
	if err != nil {
		return fmt.Errorf("Error in UpdatePositions. Error: %s", err.Error())
	}

	return nil
}

func (q SeaBattleQueries) ChangeOwnerStatus(playerId string, isOwner bool) error {
	query := fmt.Sprintf(`
		UPDATE seabattle_players
		SET is_owner=$1
		WHERE id=$2;
	`)
	_, err := q.db.Exec(query, isOwner, playerId)
	if err != nil {
		return fmt.Errorf("Error in ChangeOwnerStatus. Error: %s", err.Error())
	}

	return nil
}
