package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	downloader_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader"
	healthcheck_module "github.com/pseudoelement/go-file-downloader/src/modules/healthcheck"
)

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	healthModule := healthcheck_module.NewModule(api)
	downloaderModule := downloader_module.NewModule(api)

	healthModule.SetRoutes()
	downloaderModule.SetRoutes()

	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})
	ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"*"})

	fmt.Println("Listening port :8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handlers.CORS(methods, ttl, origins)(api)))
}
