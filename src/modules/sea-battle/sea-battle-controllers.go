package seabattle

import (
	"net/http"

	"github.com/google/uuid"
	api_module "github.com/pseudoelement/golang-utils/src/api"
	types_module "github.com/pseudoelement/golang-utils/src/types"
)

func (m *SeaBattleModule) _createRoomController(w http.ResponseWriter, req *http.Request) {
	params, err := api_module.MapQueryParams(req, "player_email", "room_name")
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	roomId := uuid.New().String()
	players := make(map[string]*Player)

	newRoom := Room{
		id:        roomId,
		name:      params["room_name"],
		players:   players,
		isPlaying: false,
		queries:   m.queries,
	}
	player := NewPlayer(params["player_email"], true, newRoom, w, req)

	if e := player.Connect(); e != nil {
		api_module.FailResponse(w, e.Error(), 400)
		return
	}

	m.rooms = append(m.rooms, newRoom)

	go player.Broadcast()

	msg := types_module.MessageJson{""}
	api_module.SuccessResponse(w, msg, 200)
}

func (m *SeaBattleModule) _connectToRoomWsController(w http.ResponseWriter, req *http.Request) {
	params, err := api_module.MapQueryParams(req, "player_email", "room_name", "room_id")
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	room, e := m.srv.findAvailableRoom(params["room_name"], params["room_name"])
	if e != nil {
		msg := types_module.MessageJson{
			Message: e.Error(),
		}
		api_module.SuccessResponse(w, msg, 200)
	}

	isOwner := len(room.players) == 0
	player := NewPlayer(params["player_email"], isOwner, room, w, req)

	if e := player.Connect(); e != nil {
		api_module.FailResponse(w, e.Error(), 400)
		return
	}

	go player.Broadcast()

	msg := types_module.MessageJson{""}
	api_module.SuccessResponse(w, msg, 200)
}
