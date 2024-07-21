package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pseudoelement/go-file-downloader/src/utils"
)

func createNewFileAsync() (*os.File, error) {
	var content string
	for i := range 1000 {
		content += utils.CreateRandomWord(15, false)
		if i < 999 {
			content += "\n"
		} else {
			content += "."
		}
	}
	f, err := utils.CreateTempFile("temp-pdf-file", "pdf", content)

	return f, err
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/pdf-file", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "./src/temp/csv-file.pdf")
	}).Methods("GET")

	methods := handlers.AllowedMethods([]string{"POST", "GET"})
	ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"*"})

	fmt.Println("Listening port :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(methods, ttl, origins)(api)))
}
