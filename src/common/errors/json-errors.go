package errors_module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type JsonError struct {
	message string
	status  int
}

func (e *JsonError) Error() string {
	return e.message
}

func (e *JsonError) Status() int {
	return e.status
}

func UnknownFieldJson(errorMsg string) ErrorWithStatus {
	fieldName := strings.TrimPrefix(errorMsg, "json: unknown field ")
	msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
	return &JsonError{message: msg, status: http.StatusBadRequest}
}

func EmptyFieldInJson(fieldName string) ErrorWithStatus {
	fieldNameToLower := strings.ToLower(fieldName)
	msg := fmt.Sprintf("Request body has empty field - %s!", fieldNameToLower)
	return &JsonError{message: msg, status: http.StatusBadRequest}
}

func BadlyFormedJson() ErrorWithStatus {
	return &JsonError{message: "Request body contains badly-formed JSON", status: http.StatusBadRequest}
}

func InvalidValueJson(unmarshalTypeError *json.UnmarshalTypeError) ErrorWithStatus {
	msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	return &JsonError{message: msg, status: http.StatusBadRequest}
}

func UnmarshalError(unmarshalType string) ErrorWithStatus {
	return &JsonError{message: fmt.Sprintf("Cannot unmarshal type %v", unmarshalType), status: http.StatusBadRequest}
}
