package models

import "encoding/json"

type WsAction string

type WsMsgJson struct {
	Action WsAction        `json:"action"`
	Data   json.RawMessage `json:"data"`
}

/*---------------------------------------------Message Data--------------------------------------------------------------- */

type ConnectionDataFromClient struct {
	ConnectedUserName string `json:"connected_user_name"`
	RoomId            string `json:"room_id"`
}

type DisconnectionDataFromClient struct {
	DisconnectedUserName string `json:"disconnected_user_name"`
	DisconnectedUserId   string `json:"disconnected_user_id"`
}

type OfferDataFromClient struct {
	OfferingUserId         string `json:"offering_user_id"`
	OfferingUserDescriptor string `json:"offering_user_descriptor"`
	TargetUserId           string `json:"target_user_id"`
}

type AnswerDataFromClient struct {
	AnsweringUserId         string `json:"answering_user_id"`
	AnsweringUserDescriptor string `json:"answering_user_descriptor"`
	TargetUserId            string `json:"target_user_id"`
}

type MicrophoneToggledDataFromClient struct {
	ToggledUserId string `json:"toggled_user_id"`
	MicEnabled    bool   `json:"mic_enabled"`
}

/*------------------------------------------------------------------------------------------------------------ */
