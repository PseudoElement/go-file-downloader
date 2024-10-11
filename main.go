package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	downloader_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader"
	games_module "github.com/pseudoelement/go-file-downloader/src/modules/games"
	healthcheck_module "github.com/pseudoelement/go-file-downloader/src/modules/healthcheck"
	"github.com/pseudoelement/go-file-downloader/src/utils/logger"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	logger := logger.NewLogger()

	healthModule := healthcheck_module.NewModule(api)
	downloaderModule := downloader_module.NewModule(api, logger)
	gamesModule := games_module.NewModule(api)

	healthModule.SetRoutes()
	downloaderModule.SetRoutes()
	gamesModule.SetRoutes()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELTE"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		MaxAge:           10,
		Debug:            true,
	})
	handler := c.Handler(r)

	fmt.Println("Listening port :8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handler))
}
