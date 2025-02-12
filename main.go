package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pseudoelement/go-file-downloader/src/db/postgres"
	"github.com/pseudoelement/go-file-downloader/src/middlewares"
	downloader_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader"
	games_module "github.com/pseudoelement/go-file-downloader/src/modules/games"
	healthcheck_module "github.com/pseudoelement/go-file-downloader/src/modules/healthcheck"
	seabattle "github.com/pseudoelement/go-file-downloader/src/modules/sea-battle"
	"github.com/pseudoelement/go-file-downloader/src/utils/logger"
	"github.com/rs/cors"
)

func getAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		panic("Add ALLOWED_ORIGINS var in .env file!")
	}

	return strings.Split(origins, "__")
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	api.Use(middlewares.TimeLoggerCommonMW)

	logger := logger.New()
	db := postgres.New()

	db.Connect()

	healthModule := healthcheck_module.NewModule(api)
	downloaderModule := downloader_module.NewModule(api, logger)
	gamesModule := games_module.NewModule(api)
	seabattleModule := seabattle.NewModule(db.Conn(), api)

	healthModule.SetRoutes()
	downloaderModule.SetRoutes()
	gamesModule.SetRoutes()
	seabattleModule.SetRoutes()

	c := cors.New(cors.Options{
		AllowedOrigins:     getAllowedOrigins(),
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Accept", "Authorization"},
		OptionsPassthrough: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return true
		// },
		AllowCredentials: true,
		MaxAge:           10,
		Debug:            true,
	})
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	api.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := c.Handler(api)

	fmt.Println("Listening port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
