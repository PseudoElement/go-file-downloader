package games_module

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

func (m *GamesModule) _downloadGame(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	gameId, ok := pathParams["id"]
	gamePathName, exists := GAME_IDS[gameId]
	if !ok || !exists {
		msg := fmt.Sprintf("%s id is invalid.", gameId)
		api_module.FailResponse(w, msg, 400)
		return
	}

	dirName, _ := os.Getwd()
	path := fmt.Sprintf("%s/src/modules/games/static/%s", dirName, gamePathName)
	http.ServeFile(w, r, path)
}
