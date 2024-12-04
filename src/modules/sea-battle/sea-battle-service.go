package seabattle

type SeaBattleService struct {
	rooms *[]Room
}

func NewSeaBattleService(rooms *[]Room) *SeaBattleService {
	return &SeaBattleService{
		rooms: rooms,
	}
}

func (srv *SeaBattleService) setPlayerPositions(email string, positions string, roomName string) {

}
