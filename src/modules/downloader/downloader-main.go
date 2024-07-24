package downloader_module

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
)

type DownloaderModule struct {
	api           *mux.Router
	downloaderSrv DownloaderService
}

func NewModule(api *mux.Router) *DownloaderModule {
	return &DownloaderModule{
		api:           api,
		downloaderSrv: *NewDownloaderService(),
	}
}

func (m *DownloaderModule) CreateFile() {
	now := time.Now()
	_, err := m.downloaderSrv.CreateTxtFileWithContentSync(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: 1_000,
		},
		ColumnsData: []types_module.TextColumnInfo{
			types_module.TextColumnInfo{
				Name:              "Id",
				Type:              "NUMBER",
				NullValuesPercent: 0,
				Min:               0,
				Max:               1000,
			},
			types_module.TextColumnInfo{
				Name:              "Name",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               5,
				Max:               10,
			},
			types_module.TextColumnInfo{
				Name:              "Surname",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               10,
				Max:               15,
			},
			types_module.TextColumnInfo{
				Name: "IsMarried",
				Type: "BOOL",
			},
			types_module.TextColumnInfo{
				Name:              "Region",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               17,
				Max:               30,
			},
			types_module.TextColumnInfo{
				Name:              "Child",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               20,
				Max:               40,
			},
			types_module.TextColumnInfo{
				Name:              "WorkTitle",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               5,
				Max:               10,
			},
		},
	})

	fmt.Println("Error ===> ", err)
	fmt.Println("It took ", time.Since(now))
}
