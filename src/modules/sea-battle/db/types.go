package seabattle_queries

type DB_RoomOnPlayerJoinRow struct {
	RoomId      string
	RoomName    string
	CreatedAt   string
	PlayerEmail string
	PlayerId    string
	IsOwner     bool
}
