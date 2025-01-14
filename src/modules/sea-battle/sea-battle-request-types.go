package seabattle

type SocketRequestMsg[T any] struct {
	Email      string `json:"player_email"`
	ActionType string `json:"action_type"`
	Data       T      `json:"data"`
}

type NewStepReqMsg struct {
	/* player step (example K1) */
	Step string `json:"step"`
}

type PlayerReadyMsg struct{}

type PlayerPositionsMsg struct {
	PlayerPositions string `json:"player_positions"`
}
