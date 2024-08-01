package downloader_module

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/mux"
	mock_constants "github.com/pseudoelement/go-file-downloader/src/constants/mock"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
)

type DownloaderModule struct {
	api           *mux.Router
	downloaderSrv *DownloaderService
}

func NewModule(api *mux.Router) *DownloaderModule {
	return &DownloaderModule{
		api:           api,
		downloaderSrv: NewDownloaderService(),
	}
}

func (m *DownloaderModule) MockCreateTextFile() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateFileWithContentSync(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	})
	fmt.Println("MockCreateTextFile Sync took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFile() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateFileWithContentSync(types_module.DownloadSqlReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "sql",
			DocName:   "bimba-production",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		TableName:       "my_first_table",
		NeedCreateTable: true,
		ColumnsData:     mock_constants.MOCK_SQL_COLUMNS_DATA,
	})
	fmt.Println("MockCreateSqlFile Sync took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFileAsync() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateFileWithContentAsync(types_module.DownloadSqlReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "sql",
			DocName:   "bimba-production",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		TableName:       "my_first_table",
		NeedCreateTable: true,
		ColumnsData:     mock_constants.MOCK_SQL_COLUMNS_DATA,
	})
	fmt.Println("MockCreateSqlFile Async took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateTextFileAsync() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateFileWithContentAsync(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	})
	fmt.Println("MockCreateTextFile Async took ", time.Since(now))

	return f, err
}
