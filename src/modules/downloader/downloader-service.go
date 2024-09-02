package downloader_module

import (
	"os"
	"reflect"

	app_errors "github.com/pseudoelement/go-file-downloader/src/errors"
	common_constants "github.com/pseudoelement/go-file-downloader/src/modules/downloader/constants/common"
	downloader_errors "github.com/pseudoelement/go-file-downloader/src/modules/downloader/errors"
	types_module "github.com/pseudoelement/go-file-downloader/src/modules/downloader/types"
	custom_utils "github.com/pseudoelement/go-file-downloader/src/utils"
	errors_module "github.com/pseudoelement/golang-utils/src/errors"
)

type DownloaderService struct{}

func NewDownloaderService() *DownloaderService {
	return &DownloaderService{}
}

func (srv *DownloaderService) CreateTempFileWithContent(body interface{}, contentCreator types_module.FileContentCreator, isAsync bool) (*os.File, error) {
	var fileContent string
	var err error
	if isAsync {
		fileContent, err = contentCreator.CreateFileContentAsync(body)
	} else {
		fileContent, err = contentCreator.CreateFileContent(body)
	}
	if err != nil {
		return nil, err
	}

	reflectedBody := reflect.ValueOf(body)
	docName := reflectedBody.FieldByName("DocName").String()
	docType := reflectedBody.FieldByName("DocType").String()
	file, err := custom_utils.CreateTempFile(docName, docType, fileContent)

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
		return &app_errors.ApiError{Message: "invalid columns_data request body prop"}
	}

	for i := 0; i < columnsDataField.Len(); i++ {
		column := columnsDataField.Index(i)

		typeField := column.FieldByName("Type")
		minField := column.FieldByName("Min")
		maxField := column.FieldByName("Max")

		if !minField.IsValid() && !maxField.IsValid() {
			continue
		}

		columnType := typeField.String()
		minValue := int(minField.Int())
		maxValue := int(maxField.Int())

		if minValue > maxValue {
			return &app_errors.ApiError{Message: "min value should be more than max value"}
		}
		if !maxField.IsZero() && maxValue < 5 {
			return &app_errors.ApiError{Message: "too short max value"}
		}
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
