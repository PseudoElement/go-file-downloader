package downloader_module

import (
	"fmt"
	"os"

	value_types "github.com/pseudoelement/go-file-downloader/src/constants/value-types"
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
			value_types.RAW_TEXT: content_creators.NewTextContentCreator(),
		},
	}
	return srv
}

func (srv *DownloaderService) CreateTxtFileWithContentSync(body interface{}) (*os.File, error) {
	var file *os.File
	var err error
	switch body.(type) {
	case types_module.DownloadSqlReqBody:
		// sqlBody := body.(types_module.DownloadSqlReqBody)
		return nil, fmt.Errorf("Method not implemented!")
	case types_module.DownloadTextReqBody:
		textBody, _ := body.(types_module.DownloadTextReqBody)
		fileContent, _ := srv.contentCreators[value_types.RAW_TEXT].CreateFileContent(textBody)
		file, err = custom_utils.CreateTempFile(textBody.DocName, textBody.DocType, fileContent)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}
