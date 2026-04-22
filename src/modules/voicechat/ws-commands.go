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
		models.CONNECT:          &OnUserConnect{voiceRoom, rooms},
		models.DISCONNECT:       &OnUserDisconnect{voiceRoom, rooms},
		models.ANSWER:           &OnAnswer{voiceRoom, rooms},
		models.OFFER:            &OnOffer{voiceRoom, rooms},
		models.USER_TOGGLED_MIC: &OnMicrophoneToggle{voiceRoom, rooms},
	}
}

/*-------------------------------------------------------------------------------------------------------- */

type OnUserConnect struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opc *OnUserConnect) Send(senderUser *User, wsMsg models.WsMsgJson) {
	for _, user := range opc.voiceRoom.users {
		if user.id == senderUser.id {
			msg := models.WsConnectionMsgToNewConnectedClient{
				Action: models.YOU_CONNECTED,
				Data: models.RoomDataToClient{
					Room: ApiRoomToClientRoom(opc.voiceRoom),
				},
			}
			err := user.conn.WriteJSON(msg)
			if err != nil {
				log.Println("[OnUserConnect_Send] YOU_CONNECTED err:", err.Error())
				_disconnectUserOnSendConnectionError(user, opc, wsMsg)
			}
		} else {
			msg := models.WsConnectionMsgToOtherClient{
				Action: models.USER_CONNECTED,
				Data: models.ConnectionDataToClient{
					ConnectedUserName: senderUser.name,
					ConnectedUserId:   senderUser.id,
					RoomId:            opc.voiceRoom.id,
				},
			}
			err := user.conn.WriteJSON(msg)
			// disconnect failed coonection and send message to sockets
			if err != nil {
				log.Println("[OnUserConnect_Send] USER_CONNECTED err:", err.Error())
				_disconnectUserOnSendConnectionError(user, opc, wsMsg)
			}
		}
	}

	opc.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_JOINED,
		Data: models.ConnectionDataToClient{
			ConnectedUserName: senderUser.name,
			ConnectedUserId:   senderUser.id,
			RoomId:            opc.voiceRoom.id,
		},
	}
}

func (opc *OnUserConnect) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {
	opc.voiceRoom.AddUser(senderPeer)
}

var _ UserWsCommand = (*OnUserConnect)(nil)

func _disconnectUserOnSendConnectionError(user *User, opc *OnUserConnect, wsMsg models.WsMsgJson) {
	user.conn.Close()
	onUserDisconnect := OnUserDisconnect{opc.voiceRoom, opc.rooms}
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

func (opd *OnUserDisconnect) Send(senderUser *User, msg models.WsMsgJson) {
	var hostUser *User = findHost(opd.voiceRoom)
	var hostName string = "noHost"
	var hostId string = ""
	if hostUser != nil {
		hostLeft := hostUser.id == senderUser.id
		if hostLeft {
			if len(opd.voiceRoom.users) > 0 {
				hostName = opd.voiceRoom.users[0].name
				hostId = opd.voiceRoom.users[0].id
			}
		} else {
			hostName = hostUser.name
			hostId = hostUser.id
		}
	}

	opd.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_LEFT,
		Data: models.DisconnectionDataToClient{
			DisconnectedUserId:   senderUser.id,
			DisconnectedUserName: senderUser.name,
			RoomId:               opd.voiceRoom.id,
			NewHostName:          hostName,
			NewHostId:            hostId,
		},
	}

	if len(opd.voiceRoom.users) == 0 {
		return
	}

	for _, user := range opd.voiceRoom.users {
		if user.id != senderUser.id {
			msg := models.WsDisconnectionMsgToClient{
				Action: models.USER_DISCONNECTED,
				Data: models.DisconnectionDataToClient{
					DisconnectedUserId:   senderUser.id,
					DisconnectedUserName: senderUser.name,
					RoomId:               opd.voiceRoom.id,
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

func (opc *OnUserDisconnect) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {
	opc.voiceRoom.RemoveUser(senderPeer.id)
	if len(opc.voiceRoom.users) == 0 {
		go opc.voiceRoom.SetDeletionTimer(opc.rooms)
	}
}

var _ UserWsCommand = (*OnUserDisconnect)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnOffer struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opd *OnOffer) Send(senderUser *User, msg models.WsMsgJson) {
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

	for _, user := range opd.voiceRoom.users {
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

func (opc *OnOffer) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnOffer)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnAnswer struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opd *OnAnswer) Send(senderUser *User, msg models.WsMsgJson) {
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

	for _, user := range opd.voiceRoom.users {
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

func (opc *OnAnswer) UpdateRoomState(senderPeer *User, msg models.WsMsgJson) {}

var _ UserWsCommand = (*OnAnswer)(nil)

/*-------------------------------------------------------------------------------------------------------- */

type OnMicrophoneToggle struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opd *OnMicrophoneToggle) Send(senderUser *User, msg models.WsMsgJson) {
	// unmarshal error already checked in UpdateRoomState()
	var micToggledData models.MicrophoneToggledDataFromClient
	utils.UnmarshalOmitEmpty(msg.Data, &micToggledData)

	for _, user := range opd.voiceRoom.users {
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

func (opc *OnMicrophoneToggle) UpdateRoomState(senderUser *User, msg models.WsMsgJson) {
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
