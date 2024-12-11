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
	queries   seabattle_queries.SeaBattleQueries
	// key - player_id
	players map[string]*Player
}

type SeaBattleModule struct {
	db      *sql.DB
	api     *mux.Router
	queries seabattle_queries.SeaBattleQueries
	srv     SeaBattleService
	rooms   []Room
}

func NewModule(db *sql.DB, api *mux.Router) SeaBattleModule {
	rooms := make([]Room, 0, 1000)
	queries := seabattle_queries.New(db)

	return SeaBattleModule{
		db:      db,
		api:     api,
		queries: queries,
		srv:     NewSeaBattleService(rooms, queries),
		rooms:   rooms,
	}
}
