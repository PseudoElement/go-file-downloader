package seabattle_queries

type DB_PlayerWithRoomJoinRow struct {
	RoomPositions string
	RoomId        string
	RoomName      string
	CreatedAt     string
	PlayerEmail   string
	PlayerId      string
	IsOwner       bool
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
