package voicechat

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/utils"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type Peer struct {
	descriptor string
	name       string
	isHost     bool
	id         string
	commands   map[models.WsAction]PeerWsCommand
	conn       *websocket.Conn
}

func NewPeer(name, descriptor string, isHost bool, commands map[models.WsAction]PeerWsCommand) *Peer {
	return &Peer{
		name:       name,
		isHost:     isHost,
		descriptor: descriptor,
		commands:   commands,
		id:         common.RandomString(),
	}
}

func (p *Peer) Connect(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
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
	p.conn = conn

	go p.broadcast(ctx)

	return nil
}

func (p *Peer) broadcast(ctx context.Context) {
	defer p.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var wsMsg models.WsMsgJson
			err := p.conn.ReadJSON(&wsMsg)
			log.Println("[Peer_Broadcast] msgBytes: ", wsMsg)
			if err != nil {
				msg := models.WsErrorMsg{Error: "invalid message"}
				p.conn.WriteJSON(msg)
				continue
			}

			wsAction := wsMsg.Action
			command, ok := p.commands[wsAction]
			if !ok {
				msg := models.WsErrorMsg{Error: "unknown action type"}
				p.conn.WriteJSON(msg)
				log.Println("[Peer_Broadcast] unknown action type: ", wsAction)
				continue
			}

			if wsAction == models.CONNECT {
				var connData models.ConnectionData
				err := utils.UnmarshalOmitEmpty(wsMsg.Data, &connData)
				log.Println("[Peer_Broadcast] connData:", connData)
				if err != nil {
					msg := models.WsErrorMsg{Error: err.Error()}
					p.conn.WriteJSON(msg)
					log.Println("[Peer_Broadcast] unmarshal err: ", err)
					continue
				}
				p.setDescriptor(connData.PeerDescriptor)
			}

			command.UpdateRoomState(p)
			command.Send(p)
		}
	}

}

func (p *Peer) setDescriptor(descriptor string) {
	p.descriptor = descriptor
}
