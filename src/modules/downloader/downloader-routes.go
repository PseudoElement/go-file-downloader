package downloader_module

import "net/http"

func (m *DownloaderModule) SetRoutes() {
	m.api.HandleFunc("/download/txt-file", m._downloadTxtFileController).Methods(http.MethodPost)
	m.api.HandleFunc("/download/test", m._testController).Methods(http.MethodGet)
}
