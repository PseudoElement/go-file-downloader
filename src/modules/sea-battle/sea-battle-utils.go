package seabattle

import (
	"fmt"
	"net/http"
)

func GetPlayersFromRoom(playerEmail string, room *Room) (RoomPlayers, bool) {
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

func MockHttpReq() *http.Request {
	return new(http.Request)
}

func MockRespWriter() http.ResponseWriter {
	return struct{ http.ResponseWriter }{}
}

func MockPlayer() *Player {
	return &Player{
		info:      PlayerInfo{},
		positions: "",
		room:      new(Room),
		rooms:     []*Room{},
	}
}
