package seabattle

import (
	"net/http"

	api_module "github.com/pseudoelement/golang-utils/src/api"
	types_module "github.com/pseudoelement/golang-utils/src/types"
)

func (m *SeaBattleModule) _createRoomController(w http.ResponseWriter, req *http.Request) {
	params, err := api_module.MapQueryParams(req, "player_email", "room_name")
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	if e := m.srv.createNewRoom(params["room_name"], params["player_email"], w, req); e != nil {
		api_module.FailResponse(w, e.Error(), 400)
		return
	}

	msg := types_module.MessageJson{Message: "Room created successfully."}
	api_module.SuccessResponse(w, msg, 200)
}

func (m *SeaBattleModule) _connectToRoomWsController(w http.ResponseWriter, req *http.Request) {
	params := api_module.MapQueryParamsSafe(req, "player_email", "room_name", "room_id")

	if e := m.srv.connectUserToToom(params["room_name"], params["room_id"], params["player_email"], w, req); e != nil {
		api_module.FailResponse(w, e.Error(), 400)
		return
	}

	msg := types_module.MessageJson{Message: "You connected to room."}
	api_module.SuccessResponse(w, msg, 200)
}
