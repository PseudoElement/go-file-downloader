package voicechat

import "net/http"

func (m *VoicechatModule) SetRoutes() {
	m.api.HandleFunc("/voicechat/create", m._createRoomHandler).Methods(http.MethodPost)
	m.api.HandleFunc("/voicechat/room", m._getRoomByIdHandler).Methods(http.MethodGet)
	m.api.HandleFunc("/voicechat/rooms", m._getRoomsListHandler).Methods(http.MethodGet)
	m.api.HandleFunc("/voicechat/ws/connect", m._connectToRoomWsHandler).Methods(http.MethodGet)
	m.api.HandleFunc("/voicechat/ws/rooms", m._getRoomsListWsHandler).Methods(http.MethodGet)
}
