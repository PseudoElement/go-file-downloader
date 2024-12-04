package seabattle

import (
	"database/sql"

	"github.com/gorilla/mux"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
)

type Room struct {
	id        string
	name      string
	isPlaying bool
	positions *string
	players   map[string]Player
}

type SeaBattleModule struct {
	db      *sql.DB
	api     *mux.Router
	queries seabattle_queries.SeaBattleQueries
	rooms   []Room
}

func NewModule(db *sql.DB, api *mux.Router) SeaBattleModule {
	return SeaBattleModule{
		db:      db,
		api:     api,
		queries: seabattle_queries.New(db),
		rooms:   make([]Room, 0, 1000),
	}
}
