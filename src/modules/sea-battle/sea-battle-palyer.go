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

type Player struct {
	email    string
	isOwner  bool
	roomName string
	conn     *websocket.Conn
	w        http.ResponseWriter
	req      *http.Request
	queries  seabattle_queries.SeaBattleQueries
}

func NewPlayer(
	email string,
	isOwner bool,
	roomName string,
	w http.ResponseWriter,
	req *http.Request,
	queries seabattle_queries.SeaBattleQueries,
) Player {
	return Player{
		email:    email,
		isOwner:  isOwner,
		w:        w,
		req:      req,
		queries:  queries,
		roomName: roomName,
	}
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
		return fmt.Errorf("Error in Player_Connect. Err: %s", err.Error())
	}

	p.conn = conn
	log.Printf("Client %s connected to %s.", p.email, p.roomName)

	return nil
}

func (p *Player) Disconnect() error {
	err := p.conn.Close()
	if err != nil {
		return fmt.Errorf("Error in Player_Disconnect_Close. Err: %s", err.Error())
	}

	err = p.queries.DisconnectPlayerFromRoom(p.email, p.roomName)
	if err != nil {
		return fmt.Errorf("Error in Player_Disconnect_DisconnectPlayerFromRoom. Err: %s", err.Error())
	}

	return nil
}

func (p *Player) Listen() {
	defer p.Disconnect()

	for {
		messageType, bytesData, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("Listen_ReadMessage err: ", err)
			return
		}

		var msgBody any
		if err := json.Unmarshal(bytesData, &msgBody); err != nil {
			log.Println("Listen_Unmarshal err: ", err)
			return
		}
		p.handlePlayerAction(msgBody)

		if err := p.conn.WriteMessage(messageType, bytesData); err != nil {
			log.Println("Listen_WriteMessage err: ", err)
			return
		}

	}
}

func (p *Player) handlePlayerAction(msgBody any) error {
	switch val := msgBody.(type) {
	case ConnectPlayerMsg:
		if val.Type == CONNECT_PLAYER {
			return p.queries.ConnectPlayerToRoom(val.Email, val.RoomName)
		} else {
			return p.queries.DisconnectPlayerFromRoom(val.Email, val.RoomName)
		}
	default:
		return fmt.Errorf("Unknown msgBody type.")
	}

}

var _ PlayerSocket = (*Player)(nil)
