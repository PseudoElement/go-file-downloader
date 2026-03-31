package voicechat

import (
	"log"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
)

type PeerWsCommand interface {
	Send(senderPeer *Peer)
	UpdateRoomState(senderPeer *Peer)
}

func CreatePeerCommandsMap(voiceRoom *VoiceRoom) map[models.WsAction]PeerWsCommand {
	return map[models.WsAction]PeerWsCommand{
		models.CONNECT:    &OnPeerConnect{voiceRoom},
		models.DISCONNECT: &OnPeerDisconnect{voiceRoom},
	}
}

/*-------------------------------------------------------------------------------------------------------- */

type OnPeerConnect struct {
	voiceRoom *VoiceRoom
}

func (opc *OnPeerConnect) Send(senderPeer *Peer) {
	for _, peer := range opc.voiceRoom.peers {
		if peer.id != senderPeer.id {
			msg := models.WsConnectionMsg{
				Action: models.CONNECT,
				Data: models.ConnectionData{
					PeerName:       senderPeer.name,
					PeerDescriptor: senderPeer.descriptor,
				},
			}
			err := peer.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnPeerConnect_Send] err:", err.Error())
			}
		}
	}
}

func (opc *OnPeerConnect) UpdateRoomState(senderPeer *Peer) {
	opc.voiceRoom.AddPeer(senderPeer)
}

var _ PeerWsCommand = (*OnPeerConnect)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnPeerDisconnect struct {
	voiceRoom *VoiceRoom
}

func (opd *OnPeerDisconnect) Send(senderPeer *Peer) {
	if len(opd.voiceRoom.peers) == 0 {
		return
	}

	hostPeer := opd.voiceRoom.peers[0]
	for _, peer := range opd.voiceRoom.peers {
		if peer.isHost {
			hostPeer = peer
			break
		}
	}

	for _, peer := range opd.voiceRoom.peers {
		if peer.id != senderPeer.id {
			msg := models.WsDisconnectionMsg{
				Action: models.DISCONNECT,
				Data: models.DisconnectionData{
					DisconnectedPeerName: senderPeer.name,
					NewHostName:          hostPeer.name,
					NewHostId:            hostPeer.id,
				},
			}
			err := peer.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnPeerDisconnect_Send] err:", err.Error())
			}
		}
	}
}

/**
 * call before Send()
 */
func (opc *OnPeerDisconnect) UpdateRoomState(senderPeer *Peer) {
	opc.voiceRoom.RemovePeer(senderPeer.id)
}

var _ PeerWsCommand = (*OnPeerDisconnect)(nil)

/*-------------------------------------------------------------------------------------------------------- */
