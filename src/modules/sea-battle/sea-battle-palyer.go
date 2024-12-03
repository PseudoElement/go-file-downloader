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
	email   string
	isOwner bool
	room    Room
	conn    *websocket.Conn
	w       http.ResponseWriter
	req     *http.Request
	queries seabattle_queries.SeaBattleQueries
}

func NewPlayer(
	email string,
	isOwner bool,
	room Room,
	w http.ResponseWriter,
	req *http.Request,
	queries seabattle_queries.SeaBattleQueries,
) Player {
	return Player{
		email:   email,
		isOwner: isOwner,
		w:       w,
		req:     req,
		queries: queries,
		room:    room,
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
	log.Printf("Client %s connected to %s.", p.email, p.room.name)

	return nil
}

func (p *Player) Conn() *websocket.Conn {
	return p.conn
}

func (p *Player) Disconnect() error {
	err := p.conn.Close()
	if err != nil {
		return fmt.Errorf("Error in Player_Disconnect_Close. Err: %s", err.Error())
	}

	err = p.queries.DisconnectPlayerFromRoom(p.email, p.room.name)
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
		if err = p.handleNewMsg(msgBody); err != nil {
			log.Println("Listen_handlePlayerAction err: ", err)
			return
		}

		if err := p.conn.WriteMessage(messageType, bytesData); err != nil {
			log.Println("Listen_WriteMessage err: ", err)
			return
		}

	}
}

func (p *Player) handleNewMsg(msgBody any) error {
	switch val := msgBody.(type) {
	case ConnectPlayerMsg:
		if val.Type == CONNECT_PLAYER {
			return p.queries.ConnectPlayerToRoom(val.Email, val.RoomName)
		} else {
			return p.queries.DisconnectPlayerFromRoom(val.Email, val.RoomName)
		}
	case NewStepMsg:
		return p.updatePositions(val.Email, val)
	default:
		return fmt.Errorf("Unknown msgBody type.")
	}

}

func (p *Player) updatePositions(activePlayerEmail string, step NewStepMsg) error {
	enemy, _ := slice_utils_module.Find(p.room.players, func(player Player) bool {
		return player.email != activePlayerEmail
	})

	splitterBetweenPlayerNameAndPositions := fmt.Sprintf("%s - ", enemy.email)
	splitterBetweenPlayers := fmt.Sprintf("__")

	splitted := strings.Split(*p.room.positions, splitterBetweenPlayerNameAndPositions)[1]
	enemyPositionsStr := strings.Split(splitted, splitterBetweenPlayers)[0]
	enemyPositionsSlice := strings.Split(enemyPositionsStr, ",")

	var newEnemyPositionsStr string
	for _, position := range enemyPositionsSlice {
		updated := position
		if strings.HasPrefix(position, step.Cell) {
			if strings.HasSuffix(position, "+") {
				// strike cell with ship
				updated += "*"
			} else {
				updated += "."
			}
		}
		newEnemyPositionsStr += updated
	}

	err := p.queries.UpdatePositions(newEnemyPositionsStr)
	if err != nil {
		return err
	}

	for _, player := range p.room.players {
		msg := NewStepMsgResp{
			ActivePlayerEmail: activePlayerEmail,
			EnemyPlayerEmail:  enemy.email,
			Cell:              step.Cell,
		}
		err = player.Conn().WriteJSON(msg)
		if err != nil {
			log.Println("player.Conn().WriteJSON err: ", err.Error())
		}
	}

	return nil
}

var _ PlayerSocket = (*Player)(nil)
