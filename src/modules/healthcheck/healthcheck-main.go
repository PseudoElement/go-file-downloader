package healthcheck_module

import "github.com/gorilla/mux"

type HealthcheckModule struct {
	api *mux.Router
}

func NewModule(api *mux.Router) *HealthcheckModule {
	return &HealthcheckModule{
		api: api,
	}
}
