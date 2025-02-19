package seabattle

type RoomPlayers struct {
	CurrentPlayer *Player
	Enemy         *Player
}

type ErrorForDB struct {
	Msg string
}
