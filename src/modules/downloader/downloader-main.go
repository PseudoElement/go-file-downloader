package downloader_module

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/mux"
	mock_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/mock"
	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	content_creators "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
)

type DownloaderModule struct {
	api             *mux.Router
	downloaderSrv   *DownloaderService
	contentCreators map[string]types_module.FileContentCreator
}

func NewModule(api *mux.Router) *DownloaderModule {
	return &DownloaderModule{
		api:           api,
		downloaderSrv: NewDownloaderService(),
		contentCreators: map[string]types_module.FileContentCreator{
			sql_constants.RAW_TEXT: content_creators.NewTextContentCreator(),
			sql_constants.SQL:      content_creators.NewSqlContentCreator(),
		},
	}
}

func (m *DownloaderModule) MockCreateTextFile() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.RAW_TEXT], false)
	fmt.Println("MockCreateTextFile Sync took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFile() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadSqlReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "sql",
			DocName:   "bimba-production",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		TableName:       "my_first_table",
		NeedCreateTable: true,
		ColumnsData:     mock_constants.MOCK_SQL_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.SQL], false)
	fmt.Println("MockCreateSqlFile Sync took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFileAsync() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadSqlReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "sql",
			DocName:   "bimba-production",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		TableName:       "my_first_table",
		NeedCreateTable: true,
		ColumnsData:     mock_constants.MOCK_SQL_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.SQL], true)
	fmt.Println("MockCreateSqlFile Async took ", time.Since(now))

	return f, err
}

func (m *DownloaderModule) MockCreateTextFileAsync() (*os.File, error) {
	now := time.Now()
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.RAW_TEXT], false)
	fmt.Println("MockCreateTextFile Async took ", time.Since(now))

	return f, err
}
