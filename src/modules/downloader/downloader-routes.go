package downloader_module

import "net/http"

func (m *DownloaderModule) SetRoutes() {
	m.api.HandleFunc("/download/txt-file", m._downloadTxtFileController).Methods(http.MethodPost)
	m.api.HandleFunc("/download/test-txt-file", m._testTextFileController).Methods(http.MethodGet)
	m.api.HandleFunc("/download/sql-file", m._downloadSqlFileController).Methods(http.MethodPost)
	m.api.HandleFunc("/download/test-sql-file", m._testSqlFileController).Methods(http.MethodGet)
}
