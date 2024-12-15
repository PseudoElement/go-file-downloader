package seabattle

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	types_module "github.com/pseudoelement/golang-utils/src/types"
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

func (srv *SeaBattleService) getRoomsList() (RoomsListResp, error) {
	dbResp, err := srv.queries.GetRoomsList()
	if err != nil {
		return RoomsListResp{}, err
	}

	var roomsList = RoomsListResp{Rooms: make(map[string]RoomsListRoomResp)}
	for _, dbRow := range dbResp {
		room, ok := roomsList.Rooms[dbRow.RoomId]
		if !ok {
			room = RoomsListRoomResp{
				RoomName:  dbRow.RoomName,
				CreatedAt: dbRow.CreatedAt,
				Players:   make([]RoomsListPlayerResp, 0, 2),
			}
		}

		player := RoomsListPlayerResp{
			PlayerId:    dbRow.PlayerId,
			PlayerEmail: dbRow.PlayerEmail,
			IsOwner:     dbRow.IsOwner,
		}
		room.Players = append(room.Players, player)
	}

	return roomsList, nil

}

func (srv *SeaBattleService) createNewRoom(roomName string, playerEmail string, w http.ResponseWriter, req *http.Request) (err error, msgWithCode *types_module.MessageWithCode) {
	if isExists, err := srv.queries.CheckRoomAlreadyExists(roomName); err != nil {
		return err, nil
	} else if isExists {
		return nil, &types_module.MessageWithCode{
			Message: fmt.Sprintf("Room with name %s elready exists.", roomName),
			Code:    ROOM_ALREADY_EXISTS,
		}
	}

	roomId := uuid.New().String()
	players := make(map[string]*Player)

	newRoom := Room{
		id:        roomId,
		name:      roomName,
		players:   players,
		isPlaying: false,
		queries:   srv.queries,
	}

	if err := srv.connectUserToToom(newRoom.name, newRoom.id, playerEmail, w, req); err != nil {
		return err, nil
	}

	srv.rooms = append(srv.rooms, newRoom)

	return nil, nil
}

func (srv *SeaBattleService) connectUserToToom(roomName string, roomId string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	room, e := srv.findRoom(roomName, roomId)
	if e != nil {
		return e
	}
	if len(room.players) == 2 {
		return fmt.Errorf("Room is full.")
	}

	player := NewPlayer(playerEmail, room, w, req)
	if e := player.Connect(); e != nil {
		return e
	}

	// add player to room
	room.players[player.info.id] = player

	go player.Broadcast()

	return nil
}

func (srv *SeaBattleService) disconnectUserFromRoom(roomId string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	room, e := srv.findRoom("", roomId)
	if e != nil {
		return e
	}

	playersOfRoom := srv.getPlayersFromRoom(playerEmail, room)

	err := playersOfRoom.CurrentPlayer.Disconnect()
	if err != nil {
		return err
	}
	delete(room.players, playersOfRoom.CurrentPlayer.info.id)

	if playersOfRoom.CurrentPlayer.info.isOwner && playersOfRoom.Enemy != nil {
		playersOfRoom.Enemy.MakeAsOwner()
	}

	return nil
}

func (srv *SeaBattleService) findRoom(name string, id string) (Room, error) {
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

	return room, nil
}

func (srv *SeaBattleService) getPlayersFromRoom(playerEmail string, room Room) RoomPlayers {
	playersOfRoom := RoomPlayers{}
	for _, player := range room.players {
		if player.info.email == playerEmail {
			playersOfRoom.CurrentPlayer = player
		} else {
			playersOfRoom.Enemy = player
		}
	}

	return playersOfRoom
}
