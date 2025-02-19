package seabattle

import (
	"fmt"
	"log"
	"net/http"

	seabattle_queries "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle/db"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

type SeaBattleService struct {
	rooms   []*Room
	queries seabattle_queries.SeaBattleQueries
}

func NewSeaBattleService(queries seabattle_queries.SeaBattleQueries) SeaBattleService {
	srv := SeaBattleService{queries: queries}
	roomsFromDB := srv.loadExistingRoomsFromDB()
	srv.rooms = roomsFromDB
	fmt.Println("ROOMS_FROM_DB ===> ", srv.rooms)

	return srv
}

func (this *SeaBattleService) loadExistingRoomsFromDB() []*Room {
	rooms := make([]*Room, 0, 100)
	roomsList, err := this.getRoomsList()
	if err != nil {
		log.Println("loadExistingRoomsFromDB_getRoomsList_ERROR ===> ", err.Error())
		return rooms
	}

	for _, room := range roomsList.Rooms {
		newRoom := &Room{
			id:         room.RoomId,
			name:       room.RoomName,
			created_at: room.CreatedAt,
			isPlaying:  false,
			queries:    this.queries,
			players:    map[string]*Player{},
		}

		for _, p := range *room.Players {
			newRoom.players[p.PlayerId] = NewPlayer(p.PlayerEmail, p.PlayerId, newRoom, this.rooms, MockRespWriter(), MockHttpReq())
		}

		rooms = append(rooms, newRoom)
	}

	return rooms
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
				Players:   new([]RoomsListPlayerResp),
			}
			roomsList.Rooms[room.RoomId] = room
		}

		player := RoomsListPlayerResp{
			PlayerId:    string(dbRow.PlayerId.Int64),
			PlayerEmail: dbRow.PlayerEmail.String,
			IsOwner:     dbRow.IsOwner.Bool,
		}
		*room.Players = append(*room.Players, player)
	}

	return roomsList, nil
}

func (this *SeaBattleService) createNewRoom(roomName string, playerEmail string) (ConnectPlayerResp, error) {
	room, err := this.queries.CreateRoom(roomName)
	if err != nil {
		return ConnectPlayerResp{}, nil
	}

	players := make(map[string]*Player)

	newRoom := &Room{
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
	if IsPlayerAlreadyConnectedToRoom(room, playerEmail) {
		return fmt.Errorf("You've already connected to room %s.", roomName)
	}
	// pass user with email that exists in room but has not player.Conn()
	if len(room.players) == 2 && !IsPlayerAlreadyAddedToRoomFromDB(room, playerEmail) {
		return fmt.Errorf("Room is full.")
	}

	player := NewPlayer(playerEmail, "", room, this.rooms, w, req)
	if e := player.Connect(); e != nil {
		return e
	}
	go player.Broadcast()

	return nil
}

func (this *SeaBattleService) disconnectUserFromRoom(roomId string, playerEmail string) error {
	room, e := this.findRoom("", roomId)
	if e != nil {
		return e
	}

	playersOfRoom, isEmpty := GetPlayersFromRoom(playerEmail, room)
	if isEmpty {
		return fmt.Errorf("Room is empty! You can't disconnect.")
	}

	err := playersOfRoom.CurrentPlayer.Disconnect(nil)
	if err != nil {
		return err
	}

	if playersOfRoom.CurrentPlayer.info.isOwner && playersOfRoom.Enemy != nil {
		if err := this.queries.ChangeOwnerStatus(playersOfRoom.Enemy.info.id, true); err != nil {
			this.queries.SaveNewError(room.name, err.Error())
		}
		playersOfRoom.Enemy.MakeAsOwner()
	}

	return nil
}

func (this *SeaBattleService) findRoom(name string, id string) (*Room, error) {
	room, err := slice_utils_module.Find(this.rooms, func(r *Room) bool {
		return r.id == id || r.name == name
	})

	if err != nil {
		if name != "" {
			return nil, fmt.Errorf("Room with name %s not found.", name)
		} else {
			return nil, fmt.Errorf("Room with id %s not found.", id)
		}
	}

	return room, nil
}

func (this *SeaBattleService) getRoomInfo(roomName string, playerEmail string) (ConnectPlayerResp, error) {
	room, err := this.findRoom(roomName, "")
	if err != nil {
		return ConnectPlayerResp{}, err
	}
	playersOfRoom, isEmpty := GetPlayersFromRoom(playerEmail, room)
	if err != nil {
		return ConnectPlayerResp{}, err
	}

	if isEmpty {
		return ConnectPlayerResp{
			RoomId:    room.id,
			RoomName:  room.name,
			CreatedAt: room.created_at,
			Player1:   PlayerInfoForClientOnConnection{},
			Player2:   PlayerInfoForClientOnConnection{},
		}, nil
	}

	var yourData PlayerInfoForClientOnConnection
	var enemyData PlayerInfoForClientOnConnection
	if playersOfRoom.CurrentPlayer != nil {
		yourData = PlayerInfoForClientOnConnection{
			PlayerId:    playersOfRoom.CurrentPlayer.info.id,
			PlayerEmail: playersOfRoom.CurrentPlayer.info.email,
			IsOwner:     playersOfRoom.CurrentPlayer.info.isOwner,
		}
	}
	if playersOfRoom.Enemy != nil {
		enemyData = PlayerInfoForClientOnConnection{
			PlayerId:    playersOfRoom.Enemy.info.id,
			PlayerEmail: playersOfRoom.Enemy.info.email,
			IsOwner:     playersOfRoom.Enemy.info.isOwner,
		}
	}

	return ConnectPlayerResp{
		RoomId:              room.id,
		RoomName:            room.name,
		CreatedAt:           room.created_at,
		SteppingPlayerEmail: "",
		Player1: PlayerInfoForClientOnConnection{
			PlayerId:    yourData.PlayerId,
			PlayerEmail: yourData.PlayerEmail,
			IsOwner:     yourData.IsOwner,
		},
		Player2: PlayerInfoForClientOnConnection{
			PlayerId:    enemyData.PlayerId,
			PlayerEmail: enemyData.PlayerEmail,
			IsOwner:     enemyData.IsOwner,
		},
	}, nil
}
