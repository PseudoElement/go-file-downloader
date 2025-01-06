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

	if roomInfo, err := m.srv.createNewRoom(params["room_name"], params["player_email"]); err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	} else {
		api_module.SuccessResponse(w, roomInfo, 200)
	}
}

func (m *SeaBattleModule) _disconnectFromRoom(w http.ResponseWriter, req *http.Request) {
	params, err := api_module.MapQueryParams(req, "player_email", "room_name")
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	if err := m.srv.disconnectUserFromRoom(params["player_email"], params["room_name"]); err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	}

	msg := types_module.MessageJson{Message: "You disconnected from room."}
	api_module.SuccessResponse(w, msg, 200)
}

func (m *SeaBattleModule) _getRoomsListController(w http.ResponseWriter, req *http.Request) {
	if rooms, err := m.srv.getRoomsList(); err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	} else {
		api_module.SuccessResponse(w, rooms, 200)
	}
}

func (m *SeaBattleModule) _connectToRoomWsController(w http.ResponseWriter, req *http.Request) {
	params := api_module.MapQueryParamsSafe(req, "player_email", "room_name", "room_id")

	if e := m.srv.connectUserToToom(params["room_name"], params["room_id"], params["player_email"], w, req); e != nil {
		api_module.FailResponse(w, e.Error(), 400)
		return
	} else {
		msg := types_module.MessageJson{Message: "You connected from room."}
		api_module.SuccessResponse(w, msg, 200)
	}
}

func (m *SeaBattleModule) _getRoomInfoController(w http.ResponseWriter, req *http.Request) {
	params := api_module.MapQueryParamsSafe(req, "room_name", "player_email")
	if room, err := m.srv.getRoomInfo(params["room_name"], params["player_email"]); err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	} else {
		api_module.SuccessResponse(w, room, 200)
	}
}
