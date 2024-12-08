package seabattle

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type PlayerInfo struct {
	email   string
	id      string
	isOwner bool
}

type Player struct {
	info          PlayerInfo
	positions     string
	room          Room
	eventHandlers EventHandlers
	conn          *websocket.Conn
	w             http.ResponseWriter
	req           *http.Request
}

func NewPlayer(
	email string,
	isOwner bool,
	room Room,
	w http.ResponseWriter,
	req *http.Request,
) Player {
	id := uuid.New().String()

	return Player{
		info: PlayerInfo{
			email:   email,
			isOwner: isOwner,
			id:      id,
		},
		eventHandlers: EventHandlers{room: room},
		w:             w,
		req:           req,
		room:          room,
	}
}

func (p *Player) queries() seabattle_queries.SeaBattleQueries {
	return p.room.queries
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
	// add player in room.players map
	p.room.players[p.info.id] = p
	if err := p.eventHandlers.handleConnection(p.info.email); err != nil {
		return err
	}

	log.Printf("Client %s connected to %s.", p.info.email, p.room.name)

	return nil
}

func (p *Player) Conn() *websocket.Conn {
	return p.conn
}

func (p *Player) Disconnect() error {
	if err := p.conn.Close(); err != nil {
		return err
	}

	if err := p.queries().DisconnectPlayerFromRoom(p.info.email, p.room.name); err != nil {
		return err
	}

	return nil
}

func (p *Player) Broadcast() {
	defer p.Disconnect()

	for {
		_, bytesData, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("Broadcast_ReadMessage err: ", err)
			return
		}

		var msgBody SocketRequestMsg[any]
		if err := json.Unmarshal(bytesData, &msgBody); err != nil {
			log.Println("Broadcast_Unmarshal err: ", err.Error())
			return
		}
		if err = p.eventHandlers.HandleNewMsg(msgBody); err != nil {
			log.Println("Broadcast_p.eventHandlers.HandleNewMsg err: ", err.Error())
			return
		}

	}
}

var _ PlayerSocket = (*Player)(nil)
