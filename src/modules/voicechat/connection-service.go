package voicechat

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

type ConnectionService struct {
	rooms     map[string]*VoiceRoom
	roomsChan chan models.WsMsgToClientJson
	conns     *[]*websocket.Conn
}

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		rooms:     make(map[string]*VoiceRoom, 10),
		roomsChan: make(chan models.WsMsgToClientJson),
	}
}

func (cs *ConnectionService) CreateRoom(reqBody models.CreateRoomReqBody) *VoiceRoom {
	voiceRoom := NewVoiceRoom(reqBody.RoomName, reqBody.MaxPeers, reqBody.HostName, cs.roomsChan)
	cs.rooms[voiceRoom.id] = voiceRoom
	log.Printf("[ConnectionService_CreateRoom] room created: %+v", *voiceRoom)

	roomModel := ApiRoomToClientRoom(voiceRoom)
	cs.roomsChan <- models.WsMsgToClientJson{
		Action: models.ROOM_CREATED,
		Data: models.RoomData{
			Room: roomModel,
		},
	}

	go voiceRoom.SetDeletionTimer(cs.rooms)

	return voiceRoom
}

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

	wsCommands := CreatePeerCommandsMap(voiceRoom, cs.rooms)
	isHost := peerName == voiceRoom.hostName
	peer := NewPeer(peerName, "", isHost, wsCommands)
	voiceRoom.AddPeer(peer)

	err := peer.Connect(context.TODO(), w, req)

	return err
}

func (cs *ConnectionService) ListenToRoomsChanges(w http.ResponseWriter, req *http.Request) error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil || conn == nil {
		return err
	}
	*cs.conns = append(*cs.conns, conn)

	return nil
}

func (cs *ConnectionService) handleRoomsChanges() {
	for msg := range cs.roomsChan {
		for _, conn := range *cs.conns {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("[ConnectionService_ListenRoomsChan] write err:", err)
			}
		}
	}
}

func (cs *ConnectionService) findPeer(voiceRoom *VoiceRoom, peerName string) *Peer {
	for _, peer := range voiceRoom.peers {
		if peer.name == peerName {
			return peer
		}
	}
	return nil
}
