package seabattle

const (
	CONNECT_PLAYER       = "CONNECT_PLAYER"
	DISCONNECT_PLAYER    = "DISCONNECT_PLAYER"
	START_GAME           = "START_GAME"
	END_GAME             = "END_GAME"
	STEP                 = "STEP"
	SET_PLAYER_POSITIONS = "SET_PLAYER_POSITIONS"
)

// Types: START_GAME, END_GAME, STEP, CONNECT_PLAYER, DISCONNECT_PLAYER
type SocketMsg struct {
	Type string `json:"type"`
}

type SetPlayerPositionsMsg struct {
	SocketMsg
	Email     string
	Positions []struct {
		/* A1 */
		Cell string
	} `json:"positions"`
}

type ConnectPlayerMsg struct {
	SocketMsg
	Email    string `json:"email"`
	RoomName string `json:"room_name"`
}

type NewStepMsg struct {
	SocketMsg
	Email string `json:"email"`
	/* player step (example K1) */
	Cell string `json:"cell"`
}

type NewStepMsgResp struct {
	ActivePlayerEmail string `json:"active_player_email"`
	EnemyPlayerEmail  string `json:"enemy_player_email"`
	Cell              string `json:"cell"`
}
