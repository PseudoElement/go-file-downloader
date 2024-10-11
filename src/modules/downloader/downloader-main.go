package downloader_module

import (
	"os"

	"github.com/gorilla/mux"
	mock_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/mock"
	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	content_creators "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
	"github.com/pseudoelement/go-file-downloader/src/utils/logger"
)

type DownloaderModule struct {
	api             *mux.Router
	downloaderSrv   *DownloaderService
	contentCreators map[string]types_module.FileContentCreator
	logger          logger.Logger
}

func NewModule(api *mux.Router, logger logger.Logger) *DownloaderModule {
	return &DownloaderModule{
		api:           api,
		downloaderSrv: NewDownloaderService(),
		contentCreators: map[string]types_module.FileContentCreator{
			sql_constants.RAW_TEXT: content_creators.NewTextContentCreator(logger),
			sql_constants.SQL:      content_creators.NewSqlContentCreator(logger),
		},
		logger: logger,
	}
}

func (m *DownloaderModule) MockCreateTextFile() (*os.File, error) {
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.RAW_TEXT], false)

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFile() (*os.File, error) {
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

	return f, err
}

func (m *DownloaderModule) MockCreateSqlFileAsync() (*os.File, error) {
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

	return f, err
}

func (m *DownloaderModule) MockCreateTextFileAsync() (*os.File, error) {
	f, err := m.downloaderSrv.CreateTempFileWithContent(types_module.DownloadTextReqBody{
		CommonReqBody: types_module.CommonReqBody{
			DocType:   "pdf",
			DocName:   "Borrow",
			RowsCount: mock_constants.ROWS_COUNT,
		},
		ColumnsData: mock_constants.MOCK_TEXT_COLUMNS_DATA,
	}, m.contentCreators[sql_constants.RAW_TEXT], true)

	return f, err
}
