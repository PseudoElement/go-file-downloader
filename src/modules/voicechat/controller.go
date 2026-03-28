package voicechat

import (
	"net/http"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

func (m *VoicechatModule) _createRoomHandler(w http.ResponseWriter, req *http.Request) {
	body, err := api_module.ParseReqBody[models.CreateRoomReqBody](w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	r := m.connectionSrv.CreateRoom(body)
	data := struct {
		Room *VoiceRoom `json:"room"`
	}{
		Room: r,
	}

	resp := models.Response{Message: "Room created.", Data: data}

	api_module.SuccessResponse(w, resp, 200)
}
