package downloader_module

import (
	"os"

	sql_constants "github.com/pseudoelement/go-file-downloader/src/constants/sql"
	content_creators "github.com/pseudoelement/go-file-downloader/src/services/content-creators"
	types_module "github.com/pseudoelement/go-file-downloader/src/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
)

type DownloaderService struct {
	contentCreators map[string]types_module.FileContentCreator
}

func NewDownloaderService() *DownloaderService {
	srv := &DownloaderService{
		contentCreators: map[string]types_module.FileContentCreator{
			sql_constants.RAW_TEXT: content_creators.NewTextContentCreator(),
			sql_constants.SQL:      content_creators.NewSqlContentCreator(),
		},
	}
	return srv
}

func (srv *DownloaderService) CreateTxtFileWithContentSync(body interface{}) (*os.File, error) {
	var file *os.File
	var err error
	switch body.(type) {
	case types_module.DownloadSqlReqBody:
		sqlBody, _ := body.(types_module.DownloadSqlReqBody)
		fileContent, er := srv.contentCreators[sql_constants.SQL].CreateFileContent(sqlBody)
		if er != nil {
			return nil, er
		}
		file, err = custom_utils.CreateTempFile(sqlBody.DocName, sqlBody.DocType, fileContent)
	case types_module.DownloadTextReqBody:
		textBody, _ := body.(types_module.DownloadTextReqBody)
		fileContent, _ := srv.contentCreators[sql_constants.RAW_TEXT].CreateFileContent(textBody)
		file, err = custom_utils.CreateTempFile(textBody.DocName, textBody.DocType, fileContent)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}
