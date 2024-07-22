package errors_module

import "net/http"

func DbDefaultError(msg string) ErrorWithStatus {
	if msg == "" {
		msg = "Error occured trying send request to database."
	}
	return &ApiError{message: msg, status: http.StatusBadRequest}
}
