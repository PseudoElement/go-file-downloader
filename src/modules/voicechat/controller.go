package voicechat

import (
	"net/http"

	"github.com/pseudoelement/go-file-downloader/src/modules/voicechat/models"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

// @Summary      Create room handler
// @Description  create new room
// @Tags         voicechat
// @Accept       json
// @Produce      json
// @Param        request body models.CreateRoomReqBody true "Request body"
// @Success      200  {object}  models.MessageJson
// @Failure      400  {object}  models.MessageJson
// @Router       /voicechat/create [post]
func (m *VoicechatModule) _createRoomHandler(w http.ResponseWriter, req *http.Request) {
	body, err := api_module.ParseReqBody[models.CreateRoomReqBody](w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	r := m.connectionSrv.CreateRoom(body)
	data := models.CreateRoomRespBody{
		Message: "Room created.",
		Data: models.MinimalRoomData{
			RoomId:   r.id,
			RoomName: r.name,
		},
	}

	api_module.SuccessResponse(w, data, 200)
}

// @Summary      Connect to room handler
// @Description  connect to existing room
// @Tags         voicechat
// @Accept       json
// @Produce      json
// @Param        room_id query string true "id of room"
// @Param        peer_name query string true "username"
// @Success      200  {object}  models.MessageJson
// @Failure      400  {object}  models.MessageJson
// @Router       /voicechat/ws/connect [get]
func (m *VoicechatModule) _connectToRoomWsHandler(w http.ResponseWriter, req *http.Request) {
	err := m.connectionSrv.ConnectToRoom(w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	}

	resp := models.MessageJson{
		Message: "Connected to room.",
	}

	api_module.SuccessResponse(w, resp, 200)
}

// @Summary      Rooms list handler
// @Description  get rooms list
// @Tags         voicechat
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.GetRoomsListRespBody
// @Router       /voicechat/rooms [get]
func (m *VoicechatModule) _getRoomsListHandler(w http.ResponseWriter, req *http.Request) {
	rooms := ApiRoomsToClientRooms(m.connectionSrv.rooms)
	resp := models.GetRoomsListRespBody{
		Message: "Rooms list.",
		Data:    rooms,
	}

	api_module.SuccessResponse(w, resp, 200)
}

// @Summary      Rooms list ws handler
// @Description  connect to socket to get rooms changes
// @Tags         voicechat
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.MessageJson
// @Failure      400  {object}  models.MessageJson
// @Router       /voicechat/ws/rooms [get]
func (m *VoicechatModule) _getRoomsListWsHandler(w http.ResponseWriter, req *http.Request) {
	err := m.connectionSrv.ListenToRoomsChanges(w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	}
	resp := models.MessageJson{
		Message: "Listening to rooms changes.",
	}

	api_module.SuccessResponse(w, resp, 200)
}
