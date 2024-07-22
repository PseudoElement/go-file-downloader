package healthcheck_module

import (
	"fmt"
	"net/http"
	"os"
)

func (m *HealthcheckModule) SetRoutes() {
	m.api.HandleFunc("/health/test-json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pathDir, _ := os.Getwd()
		jsonDir := fmt.Sprintf("%s/src/modules/healthcheck/test.json", pathDir)

		http.ServeFile(w, r, jsonDir)
	}).Methods(http.MethodGet)

	m.api.HandleFunc("/health/test-txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pathDir, _ := os.Getwd()
		txtDir := fmt.Sprintf("%s/src/modules/healthcheck/test.txt", pathDir)

		http.ServeFile(w, r, txtDir)
	}).Methods(http.MethodGet)
}
