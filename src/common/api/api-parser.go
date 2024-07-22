package api_module

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	errors_module "github.com/pseudoelement/go-file-downloader/src/common/errors"
	validators_module "github.com/pseudoelement/go-file-downloader/src/common/utils/validators"
)

func MapQueryParams(req *http.Request, queryParamsKeys ...string) (map[string]string, errors_module.ErrorWithStatus) {
	mapppedParams := make(map[string]string)
	query := req.URL.Query()

	for _, key := range queryParamsKeys {
		param := query.Get(key)
		if param == "" {
			return mapppedParams, errors_module.IncorrectQueryParams()
		}
		mapppedParams[key] = param
	}

	return mapppedParams, nil
}

/* Parses body to provided generic type and restricts body size to 1MB */
func ParseReqBody[T any](w http.ResponseWriter, req *http.Request) (T, errors_module.ErrorWithStatus) {
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	body := new(T)
	err := decoder.Decode(&body)

	if hasEmpty, emptyField := validators_module.HasEmptyField(*body); hasEmpty {
		return *body, errors_module.EmptyFieldInJson(emptyField)
	}

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return *body, errors_module.BadlyFormedJson()

		case errors.Is(err, io.ErrUnexpectedEOF):
			return *body, errors_module.BadlyFormedJson()

		case errors.As(err, &unmarshalTypeError):
			return *body, errors_module.InvalidValueJson(unmarshalTypeError)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return *body, errors_module.UnknownFieldJson(err.Error())

		case errors.Is(err, io.EOF):
			return *body, errors_module.EmptyBody()

		case err.Error() == "http: request body too large":
			return *body, errors_module.TooLargeBody()

		default:
			return *body, errors_module.IncorrectBody()
		}
	}

	return *body, nil
}
