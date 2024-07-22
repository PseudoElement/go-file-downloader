package api_module

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Get[T any](url string, params map[string]string, headers map[string]string) (response T, e error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	queryParams := req.URL.Query()

	for key, value := range params {
		queryParams.Add(key, value)
	}
	req.URL.RawQuery = queryParams.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	resBytes, _ := io.ReadAll(res.Body)
	res_struct := new(T)
	if err := json.Unmarshal(resBytes, &res_struct); err != nil {
		return *res_struct, err
	}

	defer res.Body.Close()

	return *res_struct, nil
}

func Post[T any](url string, body interface{}, headers map[string]string) (responseBytes T, e error) {
	client := &http.Client{}
	jsonValue, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBytes, _ := io.ReadAll(res.Body)
	res_struct := new(T)
	if err := json.Unmarshal(resBytes, &res_struct); err != nil {
		return *res_struct, err
	}

	defer res.Body.Close()
	return *res_struct, nil
}
