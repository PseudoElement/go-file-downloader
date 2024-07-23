package downloader_module

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
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
	_, err := m.downloaderSrv.CreateTxtFileWithContentSync(DownloadTextReqBody{
		DocType:   "pdf",
		DocName:   "Test-Sintol-4",
		RowsCount: 30000,
		ColumnsData: []TextColumnInfo{
			TextColumnInfo{
				Name:              "Id",
				Type:              "NUMBER",
				NullValuesPercent: 0,
				Min:               0,
				Max:               1000,
			},
			TextColumnInfo{
				Name:              "Name",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               5,
				Max:               10,
			},
			TextColumnInfo{
				Name:              "Surname",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               10,
				Max:               15,
			},
			TextColumnInfo{
				Name: "IsMarried",
				Type: "BOOL",
			},
			TextColumnInfo{
				Name:              "Region",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               17,
				Max:               30,
			},
			TextColumnInfo{
				Name:              "Child",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               20,
				Max:               40,
			},
			TextColumnInfo{
				Name:              "WorkTitle",
				Type:              "STRING",
				NullValuesPercent: 0,
				Min:               5,
				Max:               10,
			},
		},
	})

	fmt.Println("Error ===> ", err)
	fmt.Println("It took ", time.Since(now), " ms!")
	// fmt.Println("File ===> ", f)
}
