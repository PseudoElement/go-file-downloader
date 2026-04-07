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
	ConnectedPeerName string `json:"connected_peer_name"`
	ConnectedPeerId   string `json:"connected_peer_id"`
}

type DisconnectionDataToClient struct {
	DisconnectedUserName string `json:"disconnected_user_name"`
	DisconnectedUserId   string `json:"disconnected_user_id"`
	NewHostName          string `json:"new_host_name"`
	NewHostId            string `json:"new_host_id"`
}

type RoomDataToClient struct {
	Room VoiceRoom `json:"room"`
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
	Action WsAction            `json:"action"`
	Data   OfferDataFromClient `json:"data"`
}

type WsAnswerMessageToClient struct {
	Action WsAction             `json:"action"`
	Data   AnswerDataFromClient `json:"data"`
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

type WsRoomRemovedMsgToClient struct {
	WsRoomCreatedMsgToClient
}

/*------------------------------------------------------------------------------------------------------------------------------------ */
