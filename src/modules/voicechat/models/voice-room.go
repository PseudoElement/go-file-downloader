package models

type User struct {
	Name   string `json:"name"`
	IsHost bool   `json:"is_host"`
	Id     string `json:"id"`
	Muted  bool   `json:"muted"`
}

type VoiceRoom struct {
	Users    []User `json:"users"`
	Name     string `json:"name"`
	Id       string `json:"id"`
	MaxUsers int    `json:"max_users"`
	HostName string `json:"host_name"`
}
