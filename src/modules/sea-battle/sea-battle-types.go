package seabattle

const (
	CONNECT_PLAYER    = "CONNECT_PLAYER"
	DISCONNECT_PLAYER = "DISCONNECT_PLAYER"
	START_GAME        = "START_GAME"
	END_GAME          = "END_GAME"
	STEP              = "STEP"
)

// Types: START_GAME, END_GAME, STEP, CONNECT_PLAYER, DISCONNECT_PLAYER
type SocketMsg struct {
	Type string `json:"type"`
}

type ConnectPlayerMsg struct {
	SocketMsg
	Email    string `json:"email"`
	RoomName string `json:"room_name"`
}
