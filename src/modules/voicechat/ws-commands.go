package voicechat

import (
	"encoding/json"
	"log"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/utils"
)

type UserWsCommand interface {
	/**
	 * call before Send()
	 */
	UpdateRoomState(senderPeer *User, msg models.WsMsgJson)
	Send(senderPeer *User, msg models.WsMsgJson)
}

func CreatePeerCommandsMap(
	voiceRoom *VoiceRoom,
	rooms map[string]*VoiceRoom,
) map[models.WsAction]UserWsCommand {
	return map[models.WsAction]UserWsCommand{
		models.CONNECT:                 &OnUserConnect{voiceRoom, rooms},
		models.DISCONNECT:              &OnUserDisconnect{voiceRoom, rooms},
		models.ANSWER:                  &OnAnswer{voiceRoom, rooms},
		models.OFFER:                   &OnOffer{voiceRoom, rooms},
		models.USER_TOGGLED_MIC:        &OnMicrophoneToggle{voiceRoom, rooms},
		models.USER_VOICE_CHANGED:      &OnVoiceChangedToggle{voiceRoom, rooms},
		models.ICE_CANDIDATE_TO_SERVER: &OnIceCandidate{voiceRoom, rooms},
	}
}

/*-------------------------------------------------------------------------------------------------------- */

type OnUserConnect struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnUserConnect) Send(senderUser *User, wsMsg models.WsMsgJson) {
	for _, user := range cmd.voiceRoom.users {
		if user.id == senderUser.id {
			msg := models.WsConnectionMsgToNewConnectedClient{
				Action: models.YOU_CONNECTED,
				Data: models.RoomDataToClient{
					Room: ApiRoomToClientRoom(cmd.voiceRoom),
				},
			}
			err := user.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnUserConnect_Send] YOU_CONNECTED err:", err.Error())
				_disconnectUserOnSendConnectionError(user, cmd, wsMsg)
			}
		} else {
			msg := models.WsConnectionMsgToOtherClient{
				Action: models.USER_CONNECTED,
				Data: models.ConnectionDataToClient{
					ConnectedUserName: senderUser.name,
					ConnectedUserId:   senderUser.id,
					RoomId:            cmd.voiceRoom.id,
				},
			}
			err := user.conn.WriteJSON(msg)
			// disconnect failed coonection and send message to sockets
			if err != nil {
				log.Println("[OnUserConnect_Send] USER_CONNECTED err:", err.Error())
				_disconnectUserOnSendConnectionError(user, cmd, wsMsg)
			}
		}
	}

	cmd.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_JOINED,
		Data: models.ConnectionDataToClient{
			ConnectedUserName: senderUser.name,
			ConnectedUserId:   senderUser.id,
			RoomId:            cmd.voiceRoom.id,
		},
	}
}

func (cmd *OnUserConnect) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {
	cmd.voiceRoom.AddUser(senderPeer)
}

var _ UserWsCommand = (*OnUserConnect)(nil)

func _disconnectUserOnSendConnectionError(user *User, cmd *OnUserConnect, wsMsg models.WsMsgJson) {
	user.conn.Close()
	onUserDisconnect := OnUserDisconnect{cmd.voiceRoom, cmd.rooms}
	dataBytes, err := json.Marshal(models.DisconnectionDataFromClient{
		DisconnectedUserName: user.name,
		DisconnectedUserId:   user.id,
	})
	if err != nil {
		log.Println("[OnUserConnect_Send] json.Marshal err:", err.Error())
		return
	}
	onUserDisconnect.Send(user, models.WsMsgJson{
		Action: models.DISCONNECT,
		Data:   dataBytes,
	})
	onUserDisconnect.UpdateRoomState(user, wsMsg)
}

/*-------------------------------------------------------------------------------------------------------- */

type OnUserDisconnect struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnUserDisconnect) Send(senderUser *User, msg models.WsMsgJson) {
	var hostUser *User = findHost(cmd.voiceRoom)
	var hostName string = "noHost"
	var hostId string = ""
	if hostUser != nil {
		hostLeft := hostUser.id == senderUser.id
		if hostLeft {
			if len(cmd.voiceRoom.users) > 0 {
				hostName = cmd.voiceRoom.users[0].name
				hostId = cmd.voiceRoom.users[0].id
			}
		} else {
			hostName = hostUser.name
			hostId = hostUser.id
		}
	}

	cmd.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_LEFT,
		Data: models.DisconnectionDataToClient{
			DisconnectedUserId:   senderUser.id,
			DisconnectedUserName: senderUser.name,
			RoomId:               cmd.voiceRoom.id,
			NewHostName:          hostName,
			NewHostId:            hostId,
		},
	}

	if len(cmd.voiceRoom.users) == 0 {
		return
	}

	for _, user := range cmd.voiceRoom.users {
		if user.id != senderUser.id {
			msg := models.WsDisconnectionMsgToClient{
				Action: models.USER_DISCONNECTED,
				Data: models.DisconnectionDataToClient{
					DisconnectedUserId:   senderUser.id,
					DisconnectedUserName: senderUser.name,
					RoomId:               cmd.voiceRoom.id,
					NewHostName:          hostName,
					NewHostId:            hostId,
				},
			}
			if user.conn != nil {
				err := user.conn.WriteJSON(msg)
				if err != nil {
					log.Println("[OnUserDisconnect_Send] err:", err.Error())
				}
			}
		}
	}
}

func (cmd *OnUserDisconnect) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {
	cmd.voiceRoom.RemoveUser(senderPeer.id)
	if len(cmd.voiceRoom.users) == 0 {
		go cmd.voiceRoom.SetDeletionTimer(cmd.rooms)
	}
}

var _ UserWsCommand = (*OnUserDisconnect)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnOffer struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnOffer) Send(senderUser *User, msg models.WsMsgJson) {
	var offerMsgData models.OfferDataFromClient
	err := utils.UnmarshalOmitEmpty(msg.Data, &offerMsgData)
	if err != nil {
		msg := models.WsErrorMsgToClient{
			Action: models.ERROR,
			Data: models.WsErrorMsg{
				Error: "invalid offer message",
			},
		}
		senderUser.conn.WriteJSON(msg)
		log.Println("[OnOffer_Send] invalid offer message")
		return
	}

	log.Printf("[OnOffer] offer from %s to %s", offerMsgData.OfferingUserId, offerMsgData.TargetUserId)

	for _, user := range cmd.voiceRoom.users {
		if user.id == offerMsgData.TargetUserId {
			msg := models.WsOfferMessageToClient{
				Action: models.INCOMING_OFFER,
				Data: models.OfferDataToClient{
					OfferingUserId:         offerMsgData.OfferingUserId,
					OfferingUserDescriptor: offerMsgData.OfferingUserDescriptor,
				},
			}
			err := user.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnOffer_Send] err:", err.Error())
			}
		}
	}
}

func (cmd *OnOffer) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnOffer)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnAnswer struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnAnswer) Send(senderUser *User, msg models.WsMsgJson) {
	var answerMsgData models.AnswerDataFromClient
	err := utils.UnmarshalOmitEmpty(msg.Data, &answerMsgData)
	if err != nil {
		msg := models.WsErrorMsgToClient{
			Action: models.ERROR,
			Data: models.WsErrorMsg{
				Error: "invalid answer message",
			},
		}
		senderUser.conn.WriteJSON(msg)
		log.Println("[OnOffer_Send] invalid offer message")
		return
	}

	for _, user := range cmd.voiceRoom.users {
		if user.id == answerMsgData.TargetUserId {
			msg := models.WsAnswerMessageToClient{
				Action: models.INCOMING_ANSWER,
				Data: models.AnswerDataToClient{
					AnsweringUserId:         answerMsgData.AnsweringUserId,
					AnsweringUserDescriptor: answerMsgData.AnsweringUserDescriptor,
				},
			}
			err := user.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnUserDisconnect_Send] err:", err.Error())
			}
		}
	}
}

func (cmd *OnAnswer) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnAnswer)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnMicrophoneToggle struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnMicrophoneToggle) Send(senderUser *User, msg models.WsMsgJson) {
	// unmarshal error already checked in UpdateRoomState()
	var micToggledData models.MicrophoneToggledDataFromClient
	json.Unmarshal(msg.Data, &micToggledData)

	for _, user := range cmd.voiceRoom.users {
		msg := models.WsMicrophoneToggledMessageToClient{
			Action: models.USER_TOGGLED_MIC,
			Data:   micToggledData,
		}
		err := user.conn.WriteJSON(msg)
		if err != nil {
			log.Println("[OnMicrophoneToggle_Send] err:", err.Error())
		}
	}
}

func (cmd *OnMicrophoneToggle) UpdateRoomState(senderUser *User, msg models.WsMsgJson) {
	var data models.MicrophoneToggledDataFromClient
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		msg := models.WsErrorMsgToClient{
			Action: models.ERROR,
			Data: models.WsErrorMsg{
				Error: "invalid message. " + err.Error(),
			},
		}
		senderUser.conn.WriteJSON(msg)
		log.Println("[OnMicrophoneToggle_UpdateRoomState] invalid message:", err.Error())
		return
	}
	senderUser.muted = !data.MicEnabled
}

var _ UserWsCommand = (*OnMicrophoneToggle)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnVoiceChangedToggle struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnVoiceChangedToggle) Send(senderUser *User, msg models.WsMsgJson) {
	// unmarshal error already checked in UpdateRoomState()
	var voiceChangedData models.UserVoiceChangedDataFromClient
	json.Unmarshal(msg.Data, &voiceChangedData)

	for _, user := range cmd.voiceRoom.users {
		msg := models.WsUserVoiceChangedMessageToClient{
			Action: models.USER_VOICE_CHANGED,
			Data:   voiceChangedData,
		}
		err := user.conn.WriteJSON(msg)
		if err != nil {
			log.Println("[OnVoiceChangedToggle_Send] err:", err.Error())
		}
	}
}

func (cmd *OnVoiceChangedToggle) UpdateRoomState(senderUser *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnVoiceChangedToggle)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnIceCandidate struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (cmd *OnIceCandidate) Send(senderUser *User, msg models.WsMsgJson) {
	// unmarshal error already checked in UpdateRoomState()
	var iceCanidateData models.UserIceCandidateDataFromClient
	json.Unmarshal(msg.Data, &iceCanidateData)

	for _, user := range cmd.voiceRoom.users {
		if user.id == iceCanidateData.TargetUserId {
			msg := models.WsUserIceCandidateMessageToClient{
				Action: models.ICE_CANDIDATE_FROM_SERVER,
				Data:   iceCanidateData,
			}
			err := user.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnIceCandidate_Send] err:", err.Error())
			}
		}
	}
}

func (cmd *OnIceCandidate) UpdateRoomState(senderUser *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnIceCandidate)(nil)
