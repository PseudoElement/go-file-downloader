package voicechat

import (
	"log"
	"time"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

type VoiceRoom struct {
	users               []*User
	name                string
	id                  string
	maxUsers            int
	hostName            string
	deletionTimerActive bool
	roomsChan           chan<- models.WsMsgToClientJson
}

func NewVoiceRoom(name string, maxUsers int, hostName string, roomsChan chan<- models.WsMsgToClientJson) *VoiceRoom {
	return &VoiceRoom{
		users:               make([]*User, 0),
		name:                name,
		maxUsers:            maxUsers,
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
	if len(vr.users) == 0 {
		roomModel := ApiRoomToClientRoom(vr)
		vr.roomsChan <- models.WsMsgToClientJson{
			Action: models.ROOM_REMOVED,
			Data: models.RoomDataToClient{
				Room: roomModel,
			},
		}
		delete(rooms, vr.id)
		log.Println("[ConnectionService_CreateRoom] incative room deleted. id:", vr.id)
	}
	vr.deletionTimerActive = false
}

func (vr *VoiceRoom) AddUser(peer *User) {
	vr.users = append(vr.users, peer)
}

func (vr *VoiceRoom) RemoveUser(id string) {
	var filteredUsers []*User
	var removedUser *User
	for _, user := range vr.users {
		if user.id == id {
			removedUser = user
		} else {
			filteredUsers = append(filteredUsers, user)
		}
	}

	if removedUser.isHost && len(filteredUsers) > 0 {
		filteredUsers[0].isHost = true
	}

	vr.users = filteredUsers
}
