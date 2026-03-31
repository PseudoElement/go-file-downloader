package voicechat

import "github.com/gorilla/mux"

type VoicechatModule struct {
	api           *mux.Router
	connectionSrv *ConnectionService
}

func NewModule(api *mux.Router) *VoicechatModule {
	return &VoicechatModule{
		api:           api,
		connectionSrv: NewConnectionService(),
	}
}
