package seabattle

import (
	"fmt"
	"net/http"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type SeaBattleService struct {
	rooms   []Room
	queries seabattle_queries.SeaBattleQueries
}

func NewSeaBattleService(queries seabattle_queries.SeaBattleQueries) SeaBattleService {
	srv := SeaBattleService{
		rooms:   make([]Room, 0, 1000),
		queries: queries,
	}

	return srv
}

func (this *SeaBattleService) getRoomsList() (RoomsListResp, error) {
	dbResp, err := this.queries.GetRoomsList()
	if err != nil {
		return RoomsListResp{}, err
	}

	var roomsList = RoomsListResp{Rooms: make(map[string]RoomsListRoomResp)}
	for _, dbRow := range dbResp {
		room, ok := roomsList.Rooms[dbRow.RoomId]
		if !ok {
			room = RoomsListRoomResp{
				RoomId:    dbRow.RoomId,
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

func (this *SeaBattleService) createNewRoom(roomName string, playerEmail string, w http.ResponseWriter, req *http.Request) (ConnectPlayerResp, error) {
	if isExists, err := this.queries.CheckRoomAlreadyExists(roomName); err != nil {
		return ConnectPlayerResp{}, err
	} else if isExists {
		return ConnectPlayerResp{}, fmt.Errorf("Room with name %s already exists.", roomName)
	}

	room, err := this.queries.CreateRoom(roomName)
	if err != nil {
		return ConnectPlayerResp{}, nil
	}

	players := make(map[string]*Player)

	newRoom := Room{
		id:         room.RoomId,
		name:       room.RoomName,
		created_at: room.CreatedAt,
		players:    players,
		isPlaying:  false,
		queries:    this.queries,
	}
	this.rooms = append(this.rooms, newRoom)

	roomInfo, err := this.getRoomInfo(roomName, playerEmail)

	return roomInfo, err
}

func (this *SeaBattleService) connectUserToToom(roomName string, roomId string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	room, e := this.findRoom(roomName, roomId)
	if e != nil {
		return e
	}
	if len(room.players) == 2 {
		return fmt.Errorf("Room is full.")
	}

	player := NewPlayer(playerEmail, room, w, req)
	// @TODO somehow handle user connection on client
	if e := player.Connect(); e != nil {
		return e
	}
	go player.Broadcast()

	// add player to room
	room.players[player.info.id] = player

	return nil
}

func (this *SeaBattleService) disconnectUserFromRoom(roomId string, playerEmail string, w http.ResponseWriter, req *http.Request) error {
	room, e := this.findRoom("", roomId)
	if e != nil {
		return e
	}

	playersOfRoom, isEmpty := this.getPlayersFromRoom(playerEmail, room)
	if isEmpty {
		return fmt.Errorf("Room is empty! You can't disconnect.")
	}

	err := playersOfRoom.CurrentPlayer.Disconnect()
	if err != nil {
		return err
	}

	if playersOfRoom.CurrentPlayer.info.isOwner && playersOfRoom.Enemy != nil {
		playersOfRoom.Enemy.MakeAsOwner()
	}

	return nil
}

func (this *SeaBattleService) findRoom(name string, id string) (Room, error) {
	room, err := slice_utils_module.Find(this.rooms, func(r Room) bool {
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

func (this *SeaBattleService) getRoomInfo(roomName string, playerEmail string) (ConnectPlayerResp, error) {
	room, err := this.findRoom(roomName, "")
	playersOfRoom, isEmpty := this.getPlayersFromRoom(playerEmail, room)
	if err != nil {
		return ConnectPlayerResp{}, err
	}
	if isEmpty {
		return ConnectPlayerResp{
			RoomId:    room.id,
			RoomName:  room.name,
			CreatedAt: room.created_at,
			YourData:  PlayerInfoForClientOnConnection{},
			EnemyData: PlayerInfoForClientOnConnection{},
		}, nil
	}

	return ConnectPlayerResp{
		RoomId:    room.id,
		RoomName:  room.name,
		CreatedAt: room.created_at,
		YourData: PlayerInfoForClientOnConnection{
			PlayerId:    playersOfRoom.CurrentPlayer.info.id,
			PlayerEmail: playersOfRoom.CurrentPlayer.info.email,
			IsOwner:     playersOfRoom.CurrentPlayer.info.isOwner,
		},
		EnemyData: PlayerInfoForClientOnConnection{
			PlayerId:    playersOfRoom.Enemy.info.id,
			PlayerEmail: playersOfRoom.Enemy.info.email,
			IsOwner:     playersOfRoom.Enemy.info.isOwner,
		},
	}, nil
}

func (this *SeaBattleService) getPlayersFromRoom(playerEmail string, room Room) (RoomPlayers, bool) {
	playersOfRoom := RoomPlayers{}
	if len(room.players) == 0 {
		return playersOfRoom, true
	}

	for _, player := range room.players {
		if player.info.email == playerEmail {
			playersOfRoom.CurrentPlayer = player
		} else {
			playersOfRoom.Enemy = player
		}
	}
	fmt.Println("playersOfRoom ==> ", playersOfRoom)

	return playersOfRoom, false
}
