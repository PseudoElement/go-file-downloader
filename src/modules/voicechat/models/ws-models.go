package models

import "encoding/json"

type WsMsgJson struct {
	// CONNECT, DISCONNECT
	Action WsAction        `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type WsErrorMsg struct {
	Error string `json:"error"`
}

type ConnectionData struct {
	PeerName       string `json:"peer_name,omitempty"`
	PeerDescriptor string `json:"peer_descriptor,omitempty"`
}

type DisconnectionData struct {
	DisconnectedPeerName string `json:"disconnected_peer_name"`
	NewHostName          string `json:"new_host_name"`
	NewHostId            string `json:"new_host_id"`
}

type WsConnectionMsg struct {
	// CONNECT
	Action WsAction       `json:"action"`
	Data   ConnectionData `json:"data"`
}

type WsDisconnectionMsg struct {
	// DISCONNECT
	Action WsAction          `json:"action"`
	Data   DisconnectionData `json:"data"`
}
