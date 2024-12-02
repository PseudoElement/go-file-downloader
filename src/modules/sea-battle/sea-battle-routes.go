package seabattle

import "net/http"

func (m *SeaBattleModule) SetRoutes() {
	m.api.HandleFunc("/seabattle/create", m._createRoomController).Methods(http.MethodGet)
}
