package api_module

import (
	"encoding/json"
	"net/http"

	types_module "github.com/pseudoelement/go-file-downloader/src/common/types"
)

func SuccessResponse[T any](w http.ResponseWriter, body T, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func FailResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	var res = &types_module.MessageJson{Message: message}
	json.NewEncoder(w).Encode(res)
}

func SuccessBytesResponse[T any](value T) []byte {
	valueBytes, _ := json.Marshal(value)
	return valueBytes
}

func FailBytesResponse(message string) []byte {
	var res = &types_module.MessageJson{Message: message}
	valueBytes, _ := json.Marshal(res)
	return valueBytes
}
