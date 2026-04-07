package voicechat

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type User struct {
	name     string
	isHost   bool
	id       string
	commands map[models.WsAction]UserWsCommand
	conn     *websocket.Conn
}

func NewUser(name string, isHost bool, commands map[models.WsAction]UserWsCommand) *User {
	return &User{
		name:     name,
		isHost:   isHost,
		commands: commands,
		id:       common.RandomString(),
	}
}

func (u *User) Connect(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
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
	u.conn = conn

	go u.broadcast(ctx)

	return nil
}

func (u *User) broadcast(ctx context.Context) {
	defer u.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var wsMsg models.WsMsgJson
			err := u.conn.ReadJSON(&wsMsg)
			log.Println("[Peer_Broadcast] msgBytes: ", wsMsg)
			if err != nil {
				msg := models.WsErrorMsgToClient{
					Action: models.ERROR,
					Data:   models.WsErrorMsg{Error: "invalid message"},
				}
				u.conn.WriteJSON(msg)
				continue
			}

			wsAction := wsMsg.Action
			command, ok := u.commands[wsAction]
			if !ok {
				msg := models.WsErrorMsgToClient{
					Action: models.ERROR,
					Data:   models.WsErrorMsg{Error: "unknown action type"},
				}
				u.conn.WriteJSON(msg)
				log.Println("[Peer_Broadcast] unknown action type: ", wsAction)
				continue
			}

			command.UpdateRoomState(u)
			command.Send(u, wsMsg)
		}
	}

}
