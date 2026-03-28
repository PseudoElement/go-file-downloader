package voicechat

import "net/http"

func (m *VoicechatModule) SetRoutes() {
	m.api.HandleFunc("/voicechat/create", m._createRoomHandler).Methods(http.MethodPost)
	m.api.HandleFunc("/voicechat/connect", m._createRoomHandler).Methods(http.MethodPost)
	m.api.HandleFunc("/voicechat/listen", m._createRoomHandler).Methods(http.MethodGet)
	m.api.HandleFunc("/voicechat/disconnect", m._createRoomHandler).Methods(http.MethodGet)
}
