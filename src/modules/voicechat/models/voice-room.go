package models

type Peer struct {
	Descriptor string `json:"descriptor"`
	Name       string `json:"name"`
	IsHost     bool   `json:"is_host"`
	Id         string `json:"id"`
}

type VoiceRoom struct {
	Peers    []Peer `json:"peers"`
	Name     string `json:"name"`
	Id       string `json:"id"`
	MaxPeers int    `json:"max_peers"`
	HostName string `json:"host_name"`
}
