package voicechat

import "net/http"

func (m *VoicechatModule) SetRoutes() {
	m.api.HandleFunc("/voicechat/create", m._createRoomHandler).Methods(http.MethodPost)
	m.api.HandleFunc("/voicechat/connect", m._connectToRoomHandler).Methods(http.MethodGet)
	m.api.HandleFunc("/voicechat/rooms", m._getRoomsListHandler).Methods(http.MethodGet)
}
