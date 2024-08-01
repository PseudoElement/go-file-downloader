package downloader_module

import "net/http"

func (m *DownloaderModule) SetRoutes() {
	m.api.HandleFunc("/download/txt-file", m._downloadTxtFileController).Methods(http.MethodPost)
	m.api.HandleFunc("/download/sql-file", m._downloadSqlFileController).Methods(http.MethodPost)
	m.api.HandleFunc("/download/sync/test-txt-file", m._testTextFileController).Methods(http.MethodGet)
	m.api.HandleFunc("/download/sync/test-sql-file", m._testSqlFileController).Methods(http.MethodGet)
	m.api.HandleFunc("/download/async/test-txt-file", m._testTextFileAsyncController).Methods(http.MethodGet)
	m.api.HandleFunc("/download/async/test-sql-file", m._testSqlFileAsyncController).Methods(http.MethodGet)
}
