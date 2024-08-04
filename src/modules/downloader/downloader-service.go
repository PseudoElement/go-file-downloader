package downloader_module

import (
	"os"
	"reflect"

	app_errors "github.com/pseudoelement/go-file-downloader/src/errors"
	common_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/common"
	sql_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/sql"
	downloader_errors "github.com/pseudoelement/go-file-downloader/src/modules/downloader/errors"
	content_creators "github.com/pseudoelement/go-file-downloader/src/modules/downloader/services/content-creators"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
	errors_module "github.com/pseudoelement/golang-utils/src/errors"
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

func (srv *DownloaderService) CreateFileWithContentSync(body interface{}) (*os.File, error) {
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

func (srv *DownloaderService) CreateFileWithContentAsync(body interface{}) (*os.File, error) {
	var file *os.File
	var err error
	switch body.(type) {
	case types_module.DownloadSqlReqBody:
		sqlBody, _ := body.(types_module.DownloadSqlReqBody)
		fileContent, er := srv.contentCreators[sql_constants.SQL].CreateFileContentAsync(sqlBody)
		if er != nil {
			return nil, er
		}
		file, err = custom_utils.CreateTempFile(sqlBody.DocName, sqlBody.DocType, fileContent)
	case types_module.DownloadTextReqBody:
		textBody, _ := body.(types_module.DownloadTextReqBody)
		fileContent, er := srv.contentCreators[sql_constants.RAW_TEXT].CreateFileContentAsync(textBody)
		if er != nil {
			return nil, er
		}
		file, err = custom_utils.CreateTempFile(textBody.DocName, textBody.DocType, fileContent)
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (srv *DownloaderService) ValidateColumnParams(body interface{}) errors_module.ErrorWithStatus {
	v := reflect.ValueOf(body)

	// Check if the passed body is a pointer and get the element it points to
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get the 'ColumnsData' field
	columnsDataField := v.FieldByName("ColumnsData")
	if !columnsDataField.IsValid() {
		return &app_errors.ApiError{Message: "invalid ColumnsData request body type"}
	}

	// Iterate over the columns
	for i := 0; i < columnsDataField.Len(); i++ {
		column := columnsDataField.Index(i)

		// Get the 'Type', 'Min', and 'Max' fields
		typeField := column.FieldByName("Type")
		minField := column.FieldByName("Min")
		maxField := column.FieldByName("Max")

		// Ensure the fields are valid
		if !typeField.IsValid() || !minField.IsValid() || !maxField.IsValid() {
			continue
		}

		// Get the actual values of the fields
		columnType := typeField.String()
		minValue := int(minField.Int())
		maxValue := int(maxField.Int())

		// Validate the values against the restrictions
		restrictions, ok := common_constants.RESTRICTIONS_BY_COLUMN_TYPE[columnType]
		if !ok {
			continue
		}
		if restrictions.MaximumMin < minValue {
			return downloader_errors.InvalidMinParam(columnType, restrictions.MaximumMin)
		}
		if restrictions.MinimalMax > maxValue {
			return downloader_errors.InvalidMaxParam(columnType, restrictions.MinimalMax)
		}
	}

	return nil
}
