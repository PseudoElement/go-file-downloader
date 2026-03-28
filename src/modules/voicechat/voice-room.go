package voicechat

import (
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type VoiceRoom struct {
	peers    []Peer
	name     string
	id       string
	maxPeers int
}

func NewVoiceRoom(name string, maxPeers int) *VoiceRoom {
	return &VoiceRoom{
		peers:    make([]Peer, 0),
		name:     name,
		maxPeers: maxPeers,
		id:       common.RandomString(),
	}
}

func (vr *VoiceRoom) AddPeer(peer Peer) {
	vr.peers = append(vr.peers, peer)
}

func (vr *VoiceRoom) RemovePeer(id string) {
	var filteredPeers []Peer
	var removedPeer Peer
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
