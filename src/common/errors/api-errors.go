package errors_module

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	message string
	status  int
}

func (e *ApiError) Error() string {
	return e.message
}

func (e *ApiError) Status() int {
	return e.status
}

func IncorrectQueryParams() ErrorWithStatus {
	return &ApiError{message: "Incorrect query parameters!", status: http.StatusBadRequest}
}

func IncorrectQueryParamValue(key string) ErrorWithStatus {
	msg := fmt.Sprintf("Incorrect value in parameter `%v`", key)
	return &ApiError{message: msg, status: http.StatusBadRequest}
}

func Unauthorized() ErrorWithStatus {
	return &ApiError{message: "Unauthorized!", status: http.StatusUnauthorized}
}

func IncorrectBody() ErrorWithStatus {
	return &ApiError{message: "Incorrect request body!", status: http.StatusBadRequest}
}

func EmptyBody() ErrorWithStatus {
	return &ApiError{message: "Empty request body!", status: http.StatusBadRequest}
}

func TooLargeBody() ErrorWithStatus {
	return &ApiError{message: "Request body must not be larger than 1MB", status: http.StatusBadRequest}
}
