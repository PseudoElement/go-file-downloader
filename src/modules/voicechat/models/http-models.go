package models

type WsAction string

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
}

type CreateRoomReqBody struct {
	RoomName string `json:"room_name"`
	MaxPeers int    `json:"max_peers"`
	HostName string `json:"host_name"`
}

type CreateRoomRespBody struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
}
