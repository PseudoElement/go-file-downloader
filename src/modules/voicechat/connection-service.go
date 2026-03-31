package voicechat

import (
	"context"
	"net/http"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

type ConnectionService struct {
	rooms map[string]*VoiceRoom
}

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		rooms: make(map[string]*VoiceRoom, 10),
	}
}

func (cs *ConnectionService) CreateRoom(reqBody models.CreateRoomReqBody) *VoiceRoom {
	voiceRoom := NewVoiceRoom(reqBody.RoomName, reqBody.MaxPeers, reqBody.HostName)
	cs.rooms[voiceRoom.id] = voiceRoom

	return voiceRoom
}

// 1. new WebSocket(room_id)
// 2. CONNECT msg with descriptor
func (cs *ConnectionService) ConnectToRoom(w http.ResponseWriter, req *http.Request) error {
	params, e := api_module.MapQueryParams(req, "room_id", "peer_name")
	if e != nil {
		return e
	}

	roomId := params["room_id"]
	peerName := params["peer_name"]
	voiceRoom := cs.rooms[roomId]

	wsCommands := CreatePeerCommandsMap(voiceRoom)
	isHost := peerName == voiceRoom.hostName
	peer := NewPeer(peerName, "", isHost, wsCommands)
	voiceRoom.AddPeer(peer)

	err := peer.Connect(context.TODO(), w, req)

	return err
}
