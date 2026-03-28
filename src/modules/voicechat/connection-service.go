package voicechat

import "github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"

type ConnectionService struct {
	rooms map[string]*VoiceRoom
}

func (cs *ConnectionService) CreateRoom(reqBody models.CreateRoomReqBody) *VoiceRoom {
	voiceRoom := NewVoiceRoom(reqBody.RoomName, reqBody.MaxPeers)
	hostPeer := NewPeer(reqBody.HostName, reqBody.HostDescriptor, true)
	voiceRoom.AddPeer(hostPeer)

	cs.rooms[voiceRoom.id] = voiceRoom

	return voiceRoom
}

func (cs *ConnectionService) ConnectToRoom() {}
