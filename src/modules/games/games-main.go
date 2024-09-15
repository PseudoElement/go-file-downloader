package games_module

import (
	"github.com/gorilla/mux"
)

type GamesModule struct {
	api *mux.Router
}

func NewModule(api *mux.Router) *GamesModule {
	return &GamesModule{
		api: api,
	}
}
