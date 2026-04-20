package voicechat

import "github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"

func ApiRoomToClientRoom(room *VoiceRoom) models.VoiceRoom {
	users := make([]models.User, len(room.users))
	for idx, user := range room.users {
		users[idx] = ApiUserToClientUser(user)
	}
	return models.VoiceRoom{
		Name:     room.name,
		Id:       room.id,
		MaxUsers: room.maxUsers,
		HostName: room.hostName,
		Users:    users,
	}
}

func ApiUserToClientUser(peer *User) models.User {
	return models.User{
		Name:   peer.name,
		Id:     peer.id,
		IsHost: peer.isHost,
		Muted:  peer.muted,
	}
}

func ApiRoomsToClientRooms(rooms map[string]*VoiceRoom) []models.VoiceRoom {
	clientRooms := make([]models.VoiceRoom, 0)
	for _, room := range rooms {
		users := make([]models.User, len(room.users))
		for idx, user := range room.users {
			users[idx] = ApiUserToClientUser(user)
		}
		roomModel := ApiRoomToClientRoom(room)
		clientRooms = append(clientRooms, roomModel)
	}
	return clientRooms
}
