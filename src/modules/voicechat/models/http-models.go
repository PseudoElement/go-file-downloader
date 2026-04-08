package models

type MessageJson struct {
	Message string `json:"message"`
}

type CreateRoomReqBody struct {
	RoomName string `json:"room_name"`
	MaxUsers int    `json:"max_users"`
	HostName string `json:"host_name"`
}

type MinimalRoomData struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

type CreateRoomRespBody struct {
	CreatedRoom MinimalRoomData `json:"created_room"`
}

type GetRoomsListRespBody struct {
	Rooms []VoiceRoom `json:"rooms"`
}

type GetRoomByIdRespBody struct {
	Room *VoiceRoom `json:"room"`
}
