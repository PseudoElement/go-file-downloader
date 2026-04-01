package voicechat

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("[ConnectionService_CreateRoom] room created: %+v", *voiceRoom)

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
	voiceRoom, ok := cs.rooms[roomId]
	if !ok {
		return fmt.Errorf("Room with id %s not found.", roomId)
	}

	if voiceRoom.maxPeers >= len(voiceRoom.peers) {
		return fmt.Errorf("Room is full.")
	}

	found := cs.findPeer(voiceRoom, peerName)
	if found != nil {
		return fmt.Errorf("User %s already connected.", peerName)
	}

	wsCommands := CreatePeerCommandsMap(voiceRoom)
	isHost := peerName == voiceRoom.hostName
	peer := NewPeer(peerName, "", isHost, wsCommands)
	voiceRoom.AddPeer(peer)

	err := peer.Connect(context.TODO(), w, req)

	return err
}

func (cs *ConnectionService) findPeer(voiceRoom *VoiceRoom, peerName string) *Peer {
	for _, peer := range voiceRoom.peers {
		if peer.name == peerName {
			return peer
		}
	}
	return nil
}

func (cs *ConnectionService) ToRoomsModel() []models.VoiceRoom {
	rooms := make([]models.VoiceRoom, 0)
	for _, room := range cs.rooms {
		peers := make([]models.Peer, len(room.peers))
		for j, peer := range room.peers {
			peers[j] = models.Peer{
				Descriptor: peer.descriptor,
				Name:       peer.name,
				Id:         peer.id,
				IsHost:     peer.isHost,
			}
		}

		roomModel := models.VoiceRoom{
			Name:     room.name,
			Id:       room.id,
			MaxPeers: room.maxPeers,
			HostName: room.hostName,
		}
		rooms = append(rooms, roomModel)
	}
	return rooms
}
