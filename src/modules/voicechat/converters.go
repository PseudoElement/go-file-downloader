package voicechat

import "github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"

func ApiRoomToClientRoom(room *VoiceRoom) models.VoiceRoom {
	return models.VoiceRoom{
		Name:     room.name,
		Id:       room.id,
		MaxPeers: room.maxPeers,
		HostName: room.hostName,
	}
}

func ApiPeerToClientPeer(peer *Peer) models.Peer {
	return models.Peer{
		Descriptor: peer.descriptor,
		Name:       peer.name,
		Id:         peer.id,
		IsHost:     peer.isHost,
	}
}

func ApiRoomsToClientRooms(rooms map[string]*VoiceRoom) []models.VoiceRoom {
	clientRooms := make([]models.VoiceRoom, 0)
	for _, room := range rooms {
		peers := make([]models.Peer, len(room.peers))
		for idx, peer := range room.peers {
			peers[idx] = ApiPeerToClientPeer(peer)
		}
		roomModel := ApiRoomToClientRoom(room)
		clientRooms = append(clientRooms, roomModel)
	}
	return clientRooms
}
