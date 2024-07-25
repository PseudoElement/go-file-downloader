package downloader_module

import (
	"encoding/json"
	"net/http"

	types_module "github.com/pseudoelement/go-file-downloader/src/types"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

func (m *DownloaderModule) _downloadTxtFileController(w http.ResponseWriter, req *http.Request) {
	body, err := api_module.ParseReqBody[types_module.DownloadTextReqBody](w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	file, e := m.downloaderSrv.CreateTxtFileWithContentSync(body)
	if e != nil {
		api_module.FailResponse(w, e.Error(), 400)
	}

	w.WriteHeader(200)
	http.ServeFile(w, req, file.Name())
}

func (m *DownloaderModule) _testController(w http.ResponseWriter, req *http.Request) {
	f, err := m.MockCreateFile()
	if err != nil {
		api_module.FailResponse(w, err.Error(), 400)
		return
	}

	w.WriteHeader(200)
	http.ServeFile(w, req, f.Name())
}

func (m *DownloaderModule) _downloadSqlFileController(w http.ResponseWriter, req *http.Request) {
	body, err := api_module.ParseReqBody[types_module.DownloadSqlReqBody](w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), err.Status())
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(body)
}
