package seabattle

import (
	"fmt"
	"net/http"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

type SeabattleMW struct {
	queries seabattle_queries.SeaBattleQueries
}

func NewSeabattleMW(queries seabattle_queries.SeaBattleQueries) SeabattleMW {
	return SeabattleMW{queries: queries}
}

func (this *SeabattleMW) onlyUniqueEmail(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := api_module.MapQueryParams(req, "player_email")
		if err != nil {
			api_module.FailResponse(w, err.Error(), err.Status())
			return
		}

		isEmailTaken := this.queries.IsPlayerEmailAlreadyTaken(params["player_email"])
		if isEmailTaken {
			msg := fmt.Sprintf("Player name %s already taken. Choose another one!", params["player_email"])
			api_module.FailResponse(w, msg, 400)
			return
		}

		next(w, req)
	}
}

func (this *SeabattleMW) onlyUniqueRoomName(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := api_module.MapQueryParams(req, "room_name")
		if err != nil {
			api_module.FailResponse(w, err.Error(), err.Status())
			return
		}

		isRoomExists, _ := this.queries.IsRoomAlreadyExists(params["room_name"])
		if isRoomExists {
			msg := fmt.Sprintf("Room wit name %s already exists. Choose another one!", params["room_name"])
			api_module.FailResponse(w, msg, 400)
			return
		}

		next(w, req)
	}
}
