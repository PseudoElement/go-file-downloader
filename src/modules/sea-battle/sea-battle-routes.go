package seabattle

import "net/http"

func (m *SeaBattleModule) SetRoutes() {
	m.api.HandleFunc("/seabattle/create", m.mw.onlyUniqueEmail(m.mw.onlyUniqueRoomName(m._createRoomController))).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/connect", m.mw.onlyUniqueEmail(m._connectToRoomWsController)).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/disconnect", m._disconnectFromRoom).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/get-rooms", m._getRoomsListController).Methods(http.MethodGet)
	m.api.HandleFunc("/seabattle/get-room-info", m._getRoomInfoController).Methods(http.MethodGet)
}
