package seabattle

type SocketRespMsg[T any] struct {
	Message    string `json:"message"`
	ActionType string `json:"action_type"`
	Data       T      `json:"data"`
}

type ConnectPlayerResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
}

type PlayerReadyResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
}

type PlayerStepResp struct {
	Email string `json:"player_email"`
	Id    string `json:"player_id"`
	Step  string `json:"step"`
	// Hit/Missed/Killed/Finish
	Result string `json:"step_result"`
}

type PlayerSetPositionsResp struct{}

type RoomsListResp struct {
	// key is roomId
	Rooms map[string]RoomsListRoomResp `json:"rooms"`
}

type RoomsListRoomResp struct {
	RoomName  string                `json:"room_name"`
	CreatedAt string                `json:"created_at"`
	Players   []RoomsListPlayerResp `json:"players"`
}

type RoomsListPlayerResp struct {
	PlayerId    string `json:"player_id"`
	PlayerEmail string `json:"player_name"`
	IsOwner     bool   `json:"is_owner"`
}
