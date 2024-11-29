package healthcheck_module

import (
	"fmt"
	"net/http"
	"os"

	api_module "github.com/pseudoelement/golang-utils/src/api"
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

	m.api.HandleFunc("/health/ip", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("called!!!")
		resp := struct {
			Remote_addr   string
			Forwarded_for string
			Real_ip       string
		}{
			Remote_addr:   r.RemoteAddr,
			Forwarded_for: r.Header.Get("X-Forwarded-For"),
			Real_ip:       r.Header.Get("X-Real-Ip"),
		}
		api_module.SuccessResponse(w, resp, 200)
	}).Methods(http.MethodGet)
}
