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
	rooms       map[string]*VoiceRoom
	roomsChan   chan models.WsMsgToClientJson
	globalConns *[]*websocket.Conn
}

func NewConnectionService() *ConnectionService {
	globalConns := make([]*websocket.Conn, 0)
	return &ConnectionService{
		rooms:       make(map[string]*VoiceRoom, 10),
		roomsChan:   make(chan models.WsMsgToClientJson),
		globalConns: &globalConns,
	}
}

func (cs *ConnectionService) CreateRoom(reqBody models.CreateRoomReqBody) *VoiceRoom {
	voiceRoom := NewVoiceRoom(reqBody.RoomName, reqBody.MaxUsers, reqBody.HostName, cs.roomsChan)
	cs.rooms[voiceRoom.id] = voiceRoom
	log.Printf("[ConnectionService_CreateRoom] room created: %+v", *voiceRoom)

	roomModel := ApiRoomToClientRoom(voiceRoom)
	cs.roomsChan <- models.WsMsgToClientJson{
		Action: models.ROOM_CREATED,
		Data: models.RoomDataToClient{
			Room: roomModel,
		},
	}

	go voiceRoom.SetDeletionTimer(cs.rooms)

	return voiceRoom
}

/**
 * В этом методе проверяется возможность подключить нвоого юезра к сокету в комнате
 * Если возможно - то в комнату добалвяется новый юзер и для него открывается сокет соединение
 */
func (cs *ConnectionService) ConnectToRoom(w http.ResponseWriter, req *http.Request) error {
	params, e := api_module.MapQueryParams(req, "room_id", "user_name")
	if e != nil {
		return e
	}

	roomId := params["room_id"]
	userName := params["user_name"]
	voiceRoom, ok := cs.rooms[roomId]
	if !ok {
		return fmt.Errorf("Room with id %s not found.", roomId)
	}

	if len(voiceRoom.users) >= voiceRoom.maxUsers {
		return fmt.Errorf("Room is full.")
	}

	found := cs.findUserByName(voiceRoom, userName)
	if found != nil {
		found.conn.Close()
		err := found.Connect(context.TODO(), w, req)
		return err
	}

	wsCommands := CreatePeerCommandsMap(voiceRoom, cs.rooms)
	isHost := userName == voiceRoom.hostName

	user := NewUser(userName, isHost, wsCommands)
	err := user.Connect(context.TODO(), w, req)

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
	*cs.globalConns = append(*cs.globalConns, conn)

	return nil
}

func (cs *ConnectionService) handleRoomsChanges() {
	for msg := range cs.roomsChan {
		for _, conn := range *cs.globalConns {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("[ConnectionService_ListenRoomsChan] write err:", err)
				conn.Close()
			}
		}
	}
}

func (cs *ConnectionService) findUserByName(voiceRoom *VoiceRoom, peerName string) *User {
	for _, peer := range voiceRoom.users {
		if peer.name == peerName {
			return peer
		}
	}
	return nil
}
