package games_module

import "net/http"

func (m *GamesModule) SetRoutes() {
	m.api.HandleFunc("/games/{id}", m._downloadGame).Methods(http.MethodGet)
}
