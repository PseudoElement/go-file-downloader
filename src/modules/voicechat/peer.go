package voicechat

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
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
	p.conn = conn

	go p.broadcast(ctx)

	return err
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
			log.Println("wsMsg ==>", wsMsg)
			if err != nil {
				log.Println("[Peer_Broadcast] read err:", err.Error())
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
				connData, ok := wsMsg.Data.(models.ConnectionData)
				if !ok {
					msg := models.WsErrorMsg{Error: "invalid \"data\" field"}
					p.conn.WriteJSON(msg)
					log.Println("[Peer_Broadcast] invalid \"data\" field ")
					continue
				}
				p.setDescriptor(connData.PeerDescriptor)
			}

			command.UpdateRoomState(p)
			command.Send(p)

			// @TODO
			// 1. отправлять всем пирам в комнате свой дескриптор(+имя и айди) на подключение()
			// 2. дескрипторы всех уже подключенных пиров брать из текущей VoiceRoom и сразу отправлять в сокете на клиент
		}
	}

}

func (p *Peer) setDescriptor(descriptor string) {
	p.descriptor = descriptor
}
