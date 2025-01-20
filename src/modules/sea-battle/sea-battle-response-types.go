package seabattle

type SocketRespMsg[T any] struct {
	Message    string `json:"message"`
	ActionType string `json:"action_type"`
	Data       T      `json:"data"`
}

type ConnectPlayerResp struct {
	RoomId              string                          `json:"room_id"`
	RoomName            string                          `json:"room_name"`
	CreatedAt           string                          `json:"created_at"`
	SteppingPlayerEmail string                          `json:"stepping_player_email"`
	Player1             PlayerInfoForClientOnConnection `json:"player_1"`
	Player2             PlayerInfoForClientOnConnection `json:"player_2"`
}

type DisconnectPlayerResp struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
	Email    string `json:"player_email"`
	Id       string `json:"player_id"`
}

type PlayerReadyResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
}

type GameWinResp struct {
	WinnerEmail string `json:"winner_email"`
	WinnerId    string `json:"winner_id"`
}

type PlayerStepResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
	Step  string `json:"step"`
	// Hit/Missed/Killed/Finish
	Result string `json:"step_result"`
}

type PlayerSetPositionsResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
}

type ErrorResp struct{}

type RoomsListResp struct {
	// key is roomId
	Rooms map[string]RoomsListRoomResp `json:"rooms"`
}

type RoomsListRoomResp struct {
	RoomId    string                 `json:"room_id"`
	RoomName  string                 `json:"room_name"`
	CreatedAt string                 `json:"created_at"`
	Players   *[]RoomsListPlayerResp `json:"players"`
}

type RoomsListPlayerResp struct {
	PlayerId    string `json:"player_id"`
	PlayerEmail string `json:"player_email"`
	IsOwner     bool   `json:"is_owner"`
}

type PlayerInfoForClientOnConnection struct {
	PlayerId    string `json:"player_id"`
	PlayerEmail string `json:"player_email"`
	IsOwner     bool   `json:"is_owner"`
}
