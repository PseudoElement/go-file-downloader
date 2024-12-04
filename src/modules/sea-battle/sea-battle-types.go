package seabattle

// socket actions
const (
	CONNECT_PLAYER       = "CONNECT_PLAYER"
	DISCONNECT_PLAYER    = "DISCONNECT_PLAYER"
	START_GAME           = "START_GAME"
	END_GAME             = "END_GAME"
	STEP                 = "STEP"
	SET_PLAYER_POSITIONS = "SET_PLAYER_POSITIONS"
	ERROR                = "ERROR"
)

// step results
const (
	MISS   = "MISS"
	STRIKE = "STRIKE"
	KILL   = "KILL"
)

// Types: START_GAME, END_GAME, STEP, CONNECT_PLAYER, DISCONNECT_PLAYER, SET_PLAYER_POSITIONS
type SocketMsg struct {
	Type string `json:"type"`
}

type UpdatePlayerPositionsMsg struct {
	SocketMsg
	Email           string `json:"email"`
	PlayerPositions string `json:"player_positions"`
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
	/* who chooses cell */
	ActorEmail string `json:"actor_email"`
	/* on who field is assaulted */
	TargetEmail string `json:"target_email"`
	Cell        string `json:"cell"`
	/* MISS, STRIKE, KILL */
	Result string `json:"result"`
}

type UpdatePlayerPositionsMsgResp struct {
	/* who changed positions */
	Email string `json:"email"`
	/* not empty in error occured*/
	ErrorMsg string `json:"error"`
}
