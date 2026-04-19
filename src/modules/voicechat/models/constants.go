package models

// global events to client
const (
	ROOM_CREATED  WsAction = "ROOM_CREATED"
	ROOM_REMOVED  WsAction = "ROOM_REMOVED"
	USER_JOINED   WsAction = "USER_JOINED"
	USER_LEFT     WsAction = "USER_LEFT"
	GET_ALL_ROOMS WsAction = "GET_ALL_ROOMS"
	ERROR         WsAction = "ERROR"
)

// inner room events from client
const (
	CONNECT          WsAction = "CONNECT"
	DISCONNECT       WsAction = "DISCONNECT"
	OFFER            WsAction = "OFFER"
	ANSWER           WsAction = "ANSWER"
	USER_TOGGLED_MIC WsAction = "USER_TOGGLED_MIC"
)

// inner room events to client
const (
	YOU_CONNECTED     WsAction = "YOU_CONNECTED"
	USER_CONNECTED    WsAction = "USER_CONNECTED"
	USER_DISCONNECTED WsAction = "USER_DISCONNECTED"
	INCOMING_OFFER    WsAction = "INCOMING_OFFER"
	INCOMING_ANSWER   WsAction = "INCOMING_ANSWER"
)
