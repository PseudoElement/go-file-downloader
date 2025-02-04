package seabattle_queries

import "database/sql"

type DB_PlayerWithRoomJoinRow struct {
	RoomPositions string
	RoomId        string
	RoomName      string
	CreatedAt     string
	PlayerEmail   sql.NullString
	PlayerId      sql.NullInt64
	IsOwner       sql.NullBool
}

type DB_NewCreatedRoom struct {
	RoomId    string
	RoomName  string
	CreatedAt string
}

type DB_Player struct {
	RoomName    string
	PlayerEmail string
	PlayerId    string
	IsOwner     bool
}
