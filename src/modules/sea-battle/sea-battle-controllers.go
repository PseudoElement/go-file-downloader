package seabattle

import (
	"net/http"

	"github.com/google/uuid"
	api_module "github.com/pseudoelement/golang-utils/src/api"
	types_module "github.com/pseudoelement/golang-utils/src/types"
)

func (m *SeaBattleModule) _createRoomController(w http.ResponseWriter, req *http.Request) {
	params, err := api_module.MapQueryParams(req, "client_email", "room_name")
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	roomId := uuid.New().String()
	players := make([]Player, 0, 2)
	players = append(players, Player{email: params["client_email"], isOwner: true})

	newRoom := Room{
		id:      roomId,
		name:    params["room_name"],
		players: players,
	}
	m.rooms = append(m.rooms, newRoom)

	msg := types_module.MessageJson{
		Message: "Room created.",
	}

	api_module.SuccessResponse(w, msg, 200)
}