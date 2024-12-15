package seabattle

import "net/http"

func (m *SeaBattleModule) SetRoutes() {
	m.api.HandleFunc("/seabattle/create", m._createRoomController).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/connect", m._connectToRoomWsController).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/disconnect", m._disconenctFromRoom).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/get-rooms", m._getRoomsListController).Methods(http.MethodGet)
}
