package downloader_module

import (
	"net/http"

	api_module "github.com/pseudoelement/golang-utils/src/api"
)

func (m *DownloaderModule) _downloadTxtFileController(w http.ResponseWriter, req *http.Request) {
	body, err := api_module.ParseReqBody[DownloadReqBody](w, req)
	if err != nil {
		api_module.FailResponse(w, err.Error(), 400)
	}
}
