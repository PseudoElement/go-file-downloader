package seabattle

import (
	"fmt"

	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type SeaBattleService struct {
	rooms []Room
}

func NewSeaBattleService(rooms []Room) *SeaBattleService {
	return &SeaBattleService{
		rooms: rooms,
	}
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
