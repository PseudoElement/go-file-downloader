package games_module

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	api_module "github.com/pseudoelement/golang-utils/src/api"
	slice_utils_module "github.com/pseudoelement/golang-utils/src/utils/slices"
)

func (m *GamesModule) _downloadGame(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	gameId, ok := pathParams["id"]
	if !ok || !slice_utils_module.Contains(GAME_IDS, gameId) {
		msg := fmt.Sprintf("%s id is invalid.", gameId)
		api_module.FailResponse(w, msg, 400)
		return
	}

	dirName, _ := os.Getwd()
	path := fmt.Sprintf("%s/src/modules/games/static/%s.exe", dirName, gameId)
	http.ServeFile(w, r, path)
}
