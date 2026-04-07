package models

// global events
const (
	ROOM_CREATED WsAction = "ROOM_CREATED"
	ROOM_REMOVED WsAction = "ROOM_REMOVED"
	ERROR        WsAction = "ERROR"
)

// inner room events from client
const (
	CONNECT    WsAction = "CONNECT"
	DISCONNECT WsAction = "DISCONNECT"
	OFFER      WsAction = "OFFER"
	ANSWER     WsAction = "ANSWER"
)

// inner room events to client
const (
	YOU_CONNECTED     WsAction = "YOU_CONNECTED"
	USER_CONNECTED    WsAction = "USER_CONNECTED"
	USER_DISCONNECTED WsAction = "USER_DISCONNECTED"
	OFFER_CREATED     WsAction = "OFFER_CREATED"
	ANSWER_CREATED    WsAction = "ANSWER_CREATED"
)
