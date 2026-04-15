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
	UpdateRoomState(senderPeer *User)
	Send(senderPeer *User, msg models.WsMsgJson)
}

func CreatePeerCommandsMap(
	voiceRoom *VoiceRoom,
	rooms map[string]*VoiceRoom,
) map[models.WsAction]UserWsCommand {
	return map[models.WsAction]UserWsCommand{
		models.CONNECT:    &OnUserConnect{voiceRoom, rooms},
		models.DISCONNECT: &OnUserDisconnect{voiceRoom, rooms},
		models.ANSWER:     &OnAnswer{voiceRoom, rooms},
		models.OFFER:      &OnOffer{voiceRoom, rooms},
	}
}

/*-------------------------------------------------------------------------------------------------------- */

type OnUserConnect struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opc *OnUserConnect) Send(senderUser *User, msg models.WsMsgJson) {
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
				_disconnectUserOnSendConnectionError(user, opc)
			}
		} else {
			msg := models.WsConnectionMsgToOtherClient{
				Action: models.USER_CONNECTED,
				Data: models.ConnectionDataToClient{
					ConnectedUserName: senderUser.name,
					ConnectedUserId:   senderUser.id,
				},
			}
			err := user.conn.WriteJSON(msg)
			// disconnect failed coonection and send message to sockets
			if err != nil {
				log.Println("[OnUserConnect_Send] USER_CONNECTED err:", err.Error())
				_disconnectUserOnSendConnectionError(user, opc)
			}
		}
	}

	opc.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_JOINED,
		Data: models.ConnectionDataToClient{
			ConnectedUserName: senderUser.name,
			ConnectedUserId:   senderUser.id,
		},
	}
}

func (opc *OnUserConnect) UpdateRoomState(senderPeer *User) {
	opc.voiceRoom.AddUser(senderPeer)
}

var _ UserWsCommand = (*OnUserConnect)(nil)

func _disconnectUserOnSendConnectionError(user *User, opc *OnUserConnect) {
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
	onUserDisconnect.UpdateRoomState(user)
}

/*-------------------------------------------------------------------------------------------------------- */

type OnUserDisconnect struct {
	voiceRoom *VoiceRoom
	rooms     map[string]*VoiceRoom
}

func (opd *OnUserDisconnect) Send(senderUser *User, msg models.WsMsgJson) {
	if len(opd.voiceRoom.users) == 0 {
		return
	}

	var hostUser *User = findHost(opd.voiceRoom)
	if hostUser.id == senderUser.id {
		hostUser = opd.voiceRoom.users[0]
	}

	for _, user := range opd.voiceRoom.users {
		if user.id != senderUser.id {
			msg := models.WsDisconnectionMsgToClient{
				Action: models.USER_DISCONNECTED,
				Data: models.DisconnectionDataToClient{
					DisconnectedUserId:   senderUser.id,
					DisconnectedUserName: senderUser.name,
					NewHostName:          hostUser.name,
					NewHostId:            hostUser.id,
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

	opd.voiceRoom.roomsChan <- models.WsMsgToClientJson{
		Action: models.USER_LEFT,
		Data: models.DisconnectionDataToClient{
			DisconnectedUserId:   senderUser.id,
			DisconnectedUserName: senderUser.name,
			NewHostName:          hostUser.name,
			NewHostId:            hostUser.id,
		},
	}
}

func (opc *OnUserDisconnect) UpdateRoomState(senderPeer *User) {
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
				Action: models.OFFER_CREATED,
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

func (opc *OnOffer) UpdateRoomState(senderPeer *User) {}

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
				Action: models.ANSWER_CREATED,
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

func (opc *OnAnswer) UpdateRoomState(senderPeer *User) {}

var _ UserWsCommand = (*OnAnswer)(nil)

/*-------------------------------------------------------------------------------------------------------- */
