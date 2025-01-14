package seabattle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type PlayerInfo struct {
	email   string
	id      string
	isOwner bool
	isReady bool
}

type Player struct {
	info          PlayerInfo
	positions     string
	rooms         []*Room
	room          *Room
	eventHandlers EventHandlers
	conn          *websocket.Conn
	w             http.ResponseWriter
	req           *http.Request
}

func NewPlayer(
	email string,
	id string,
	room *Room,
	rooms []*Room,
	w http.ResponseWriter,
	req *http.Request,
) *Player {
	return &Player{
		info: PlayerInfo{
			email:   email,
			isOwner: len(room.players) == 0,
			id:      id,
		},
		eventHandlers: NewEventHandlers(room, rooms),
		w:             w,
		req:           req,
		room:          room,
	}
}

func (p *Player) queries() seabattle_queries.SeaBattleQueries {
	return p.room.queries
}

func (p *Player) setIdFromDB(playerId string) {
	p.info.id = playerId
}

func (p *Player) setReadyStatus(isReady bool) {
	p.info.isReady = isReady
}

func (p *Player) MakeAsOwner() {
	p.info.isOwner = true
}

func (p *Player) Connect() error {
	allowedOrigins, ok := os.LookupEnv("ALLOWED_ORIGINS")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if !ok {
				return true
			}

			originsSlice := strings.Split(allowedOrigins, "__")
			origin := p.req.Header.Get("Origin")

			return slice_utils_module.Contains(originsSlice, origin)
		},
	}

	conn, err := upgrader.Upgrade(p.w, p.req, nil)
	if err != nil {
		return err
	}

	p.conn = conn

	playerId, err := p.queries().ConnectPlayerToRoom(p.info.email, p.room.name, p.info.isOwner)
	if err != nil {
		return err
	}

	p.setIdFromDB(playerId)
	// add player in room.players map
	p.room.players[playerId] = p

	if err := p.eventHandlers.handleConnection(p.info.email); err != nil {
		return err
	}

	if err := p.sendMsgToClientOnConnection(); err != nil {
		return err
	}

	log.Printf("Client %s connected to room `%s.`", p.info.email, p.room.name)

	return nil
}

func (p *Player) sendMsgToClientOnConnection() error {
	playersOfRoom, _ := GetPlayersFromRoom(p.info.email, p.room)

	var enemyData PlayerInfoForClientOnConnection
	if playersOfRoom.Enemy != nil {
		enemyData = PlayerInfoForClientOnConnection{
			PlayerId:    playersOfRoom.Enemy.info.id,
			PlayerEmail: playersOfRoom.Enemy.info.email,
			IsOwner:     playersOfRoom.Enemy.info.isOwner,
		}
	}

	msg := SocketRespMsg[ConnectPlayerResp]{
		Message:    fmt.Sprintf("Player %s connected to room %s.", p.info.email, p.room.name),
		ActionType: CONNECT_PLAYER,
		Data: ConnectPlayerResp{
			RoomId:    p.room.id,
			RoomName:  p.room.name,
			CreatedAt: p.room.created_at,
			YourData: PlayerInfoForClientOnConnection{
				PlayerId:    p.info.id,
				PlayerEmail: p.info.email,
				IsOwner:     p.info.isOwner,
			},
			EnemyData: enemyData,
		},
	}

	err := p.Conn().WriteJSON(msg)

	return err
}

func (p *Player) Conn() *websocket.Conn {
	return p.conn
}

func (p *Player) Disconnect() error {
	if err := p.Conn().Close(); err != nil {
		return err
	}
	if err := p.eventHandlers.handleDisconnection(p.info.email); err != nil {
		return err
	}

	return nil
}

func (p *Player) Broadcast() {
	defer p.Disconnect()

	for {
		_, bytesData, err := p.Conn().ReadMessage()

		if err != nil {
			log.Println("Broadcast_ReadMessage err =====> ", err)
			return
		}

		var msgBody SocketRequestMsg[any]
		if err := json.Unmarshal(bytesData, &msgBody); err != nil {
			log.Println("Broadcast_Unmarshal =====> ", err.Error())
			return
		}
		if err = p.eventHandlers.HandleNewMsg(msgBody); err != nil {
			log.Println("Broadcast_p.eventHandlers.HandleNewMsg  =====> ", err.Error())
			return
		}

	}
}

var _ PlayerSocket = (*Player)(nil)
