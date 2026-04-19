package models

type WsMsgToClientJson struct {
	Action WsAction `json:"action"`
	Data   any      `json:"data"`
}

type WsErrorMsg struct {
	Error string `json:"error"`
}

/*---------------------------------------------Message Data--------------------------------------------------------------- */

type ConnectionDataToClient struct {
	ConnectedUserName string `json:"connected_user_name"`
	ConnectedUserId   string `json:"connected_user_id"`
	RoomId            string `json:"room_id"`
}

type DisconnectionDataToClient struct {
	DisconnectedUserName string `json:"disconnected_user_name"`
	DisconnectedUserId   string `json:"disconnected_user_id"`
	RoomId               string `json:"room_id"`
	NewHostName          string `json:"new_host_name"`
	NewHostId            string `json:"new_host_id"`
}

type RoomDataToClient struct {
	Room VoiceRoom `json:"room"`
}

type OfferDataToClient struct {
	OfferingUserId         string `json:"offering_user_id"`
	OfferingUserDescriptor string `json:"offering_user_descriptor"`
}

type AnswerDataToClient struct {
	AnsweringUserId         string `json:"answering_user_id"`
	AnsweringUserDescriptor string `json:"answering_user_descriptor"`
}

/*------------------------------------------------Messages to client------------------------------------------------------------ */

type WsConnectionMsgToOtherClient struct {
	Action WsAction               `json:"action"`
	Data   ConnectionDataToClient `json:"data"`
}

type WsConnectionMsgToNewConnectedClient struct {
	Action WsAction         `json:"action"`
	Data   RoomDataToClient `json:"data"`
}

type WsDisconnectionMsgToClient struct {
	Action WsAction                  `json:"action"`
	Data   DisconnectionDataToClient `json:"data"`
}

type WsOfferMessageToClient struct {
	Action WsAction          `json:"action"`
	Data   OfferDataToClient `json:"data"`
}

type WsAnswerMessageToClient struct {
	Action WsAction           `json:"action"`
	Data   AnswerDataToClient `json:"data"`
}

type WsMicrophoneToggledMessageToClient struct {
	Action WsAction                        `json:"action"`
	Data   MicrophoneToggledDataFromClient `json:"data"`
}

/*----------------------------------------------------Global rooms messages-------------------------------------------------------- */

type WsErrorMsgToClient struct {
	Action WsAction   `json:"action"`
	Data   WsErrorMsg `json:"data"`
}

type WsRoomCreatedMsgToClient struct {
	Action WsAction         `json:"action"`
	Data   RoomDataToClient `json:"data"`
}

type WsUserJoinedMsgToClient struct {
	Action WsAction               `json:"action"`
	Data   ConnectionDataToClient `json:"data"`
}

type WsRoomRemovedMsgToClient struct {
	WsRoomCreatedMsgToClient
}

/*------------------------------------------------------------------------------------------------------------------------------------ */
