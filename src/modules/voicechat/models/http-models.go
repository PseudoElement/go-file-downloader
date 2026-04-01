package models

type MessageJson struct {
	Message string `json:"message"`
}

type WsAction string

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
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

type GetRoomsListRespBody struct {
	Response
	Data []VoiceRoom `json:"data"`
}
