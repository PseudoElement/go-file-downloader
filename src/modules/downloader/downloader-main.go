package downloader_module

import "github.com/gorilla/mux"

type DownloaderModule struct {
	api *mux.Router
}

func NewModule(api *mux.Router) *DownloaderModule {
	return &DownloaderModule{
		api: api,
	}
}
