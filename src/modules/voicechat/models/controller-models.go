package models

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type CreateRoomReqBody struct {
	RoomName       string `json:"room_name"`
	MaxPeers       int    `json:"max_peers"`
	HostDescriptor string `json:"host_descriptor"`
	HostName       string `json:"host_name"`
}

type WsActionJson struct {
	// CONNECT, DISCONNECT
	Action string `json:"action"`
	Peer   struct {
		PeerName       string `json:"peer_name"`
		PeerDescriptor string `json:"peer_descriptor"`
	} `json:"peer"`
}
