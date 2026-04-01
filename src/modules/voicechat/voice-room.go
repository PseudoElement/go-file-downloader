package voicechat

import (
	"log"
	"time"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type VoiceRoom struct {
	peers               []*Peer
	name                string
	id                  string
	maxPeers            int
	hostName            string
	deletionTimerActive bool
	roomsChan           chan<- models.WsMsgToClientJson
}

func NewVoiceRoom(name string, maxPeers int, hostName string, roomsChan chan<- models.WsMsgToClientJson) *VoiceRoom {
	return &VoiceRoom{
		peers:               make([]*Peer, 0),
		name:                name,
		maxPeers:            maxPeers,
		hostName:            hostName,
		deletionTimerActive: false,
		id:                  common.RandomString(),
	}
}

func (vr *VoiceRoom) SetDeletionTimer(rooms map[string]*VoiceRoom) {
	if vr.deletionTimerActive {
		return
	}

	vr.deletionTimerActive = true
	time.Sleep(1 * time.Minute)
	if len(vr.peers) == 0 {
		roomModel := ApiRoomToClientRoom(vr)
		vr.roomsChan <- models.WsMsgToClientJson{
			Action: models.ROOM_REMOVED,
			Data: models.RoomData{
				Room: roomModel,
			},
		}
		delete(rooms, vr.id)
		log.Println("[ConnectionService_CreateRoom] incative room deleted. id:", vr.id)
	}
	vr.deletionTimerActive = false
}

func (vr *VoiceRoom) AddPeer(peer *Peer) {
	vr.peers = append(vr.peers, peer)
}

func (vr *VoiceRoom) RemovePeer(id string) {
	var filteredPeers []*Peer
	var removedPeer *Peer
	for _, peer := range vr.peers {
		if peer.id == id {
			removedPeer = peer
		} else {
			filteredPeers = append(filteredPeers, peer)
		}
	}

	if removedPeer.isHost && len(filteredPeers) > 0 {
		filteredPeers[0].isHost = true
	}

	vr.peers = filteredPeers
}
