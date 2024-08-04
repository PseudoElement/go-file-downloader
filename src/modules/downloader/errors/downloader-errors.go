package downloader_errors

import (
	"fmt"

	app_errors "github.com/pseudoelement/go-file-downloader/src/errors"
	errors_module "github.com/pseudoelement/golang-utils/src/errors"
)

func InvalidMinParam(columnType string, maxMin int) errors_module.ErrorWithStatus {
	return &app_errors.ApiError{Message: fmt.Sprintf("Min param for %s column type should be less than or equal to %v!", columnType, maxMin)}
}

func InvalidMaxParam(columnType string, minMax int) errors_module.ErrorWithStatus {
	return &app_errors.ApiError{Message: fmt.Sprintf("Max param for %s column type should be more than or equal to %v!", columnType, minMax)}
}
