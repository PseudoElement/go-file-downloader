package seabattle

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type SeaBattleService struct {
	rooms   []Room
	queries seabattle_queries.SeaBattleQueries
}

func NewSeaBattleService(rooms []Room, queries seabattle_queries.SeaBattleQueries) SeaBattleService {
	return SeaBattleService{
		rooms:   rooms,
		queries: queries,
	}
}

func (srv *SeaBattleService) createNewRoom(roomName string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	roomId := uuid.New().String()
	players := make(map[string]*Player)

	newRoom := Room{
		id:        roomId,
		name:      roomName,
		players:   players,
		isPlaying: false,
		queries:   srv.queries,
	}
	player := NewPlayer(playerEmail, true, newRoom, w, req)

	if e := player.Connect(); e != nil {
		return e
	}

	srv.rooms = append(srv.rooms, newRoom)

	go player.Broadcast()

	return nil
}

func (srv *SeaBattleService) connectUserToToom(roomName string, roomId string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	room, e := srv.findAvailableRoom(roomName, roomId)
	if e != nil {
		return e
	}

	isOwner := len(room.players) == 0
	player := NewPlayer(playerEmail, isOwner, room, w, req)

	if e := player.Connect(); e != nil {
		return e
	}

	go player.Broadcast()

	return nil
}

func (srv *SeaBattleService) findAvailableRoom(name string, id string) (Room, error) {
	room, err := slice_utils_module.Find(srv.rooms, func(r Room) bool {
		return r.id == id || r.name == name
	})

	if err != nil {
		if name != "" {
			return Room{}, fmt.Errorf("Room with name %s not found.", name)
		} else {
			return Room{}, fmt.Errorf("Room with id %s not found.", id)
		}
	}

	if len(room.players) == 2 {
		return Room{}, fmt.Errorf("Room is full.")
	}

	return room, nil
}
