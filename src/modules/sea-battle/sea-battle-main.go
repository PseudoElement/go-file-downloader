package seabattle

import (
	"database/sql"

	"github.com/gorilla/mux"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
)

type Room struct {
	id         string
	name       string
	created_at string
	isPlaying  bool
	queries    seabattle_queries.SeaBattleQueries
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
	queries := seabattle_queries.New(db)
	srv := NewSeaBattleService(queries)

	if err := queries.CreateTables(); err != nil {
		panic(err)
	}

	return SeaBattleModule{
		db:      db,
		api:     api,
		queries: queries,
		srv:     srv,
		rooms:   srv.rooms,
	}
}
