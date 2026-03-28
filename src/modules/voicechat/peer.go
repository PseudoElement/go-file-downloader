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
	conn       *websocket.Conn
}

func NewPeer(name, descriptor string, isHost bool) Peer {
	return Peer{
		descriptor: descriptor,
		name:       name,
		isHost:     isHost,
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
	go p._broadcast(context.TODO())

	return err
}

func (p *Peer) _broadcast(ctx context.Context) {
	defer p.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var wsAction models.WsActionJson
			err := p.conn.ReadJSON(&wsAction)
			log.Println("wsAction ==>", wsAction)
			if err != nil {
				log.Println("[Peer_broadcast] read err:", err.Error())
			}
			// @TODO
			// 1. отправлять всем пирам в комнате свой дескриптор(+имя и айди) на подключение()
			// 2. дескрипторы всех уже подключенных пиров брать из текущей VoiceRoom и сразу отправлять в сокете на клиент
		}
	}
}
