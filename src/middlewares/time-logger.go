package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pseudoelement/go-file-downloader/src/utils/common"
)

func TimeLoggerCommonMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		now := time.Now()
		ipAddr := common.GetClientIP(req, false)
		log.Printf("Start request from IP - %s", ipAddr)

		bytesBody, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bytesBody))

		if req.Method == "POST" {
			log.Printf("Request body: %s", string(bytesBody))
		}
		if len(req.URL.RawQuery) > 0 {
			keyValuePairs := strings.Split(req.URL.RawQuery, "&")
			if len(keyValuePairs) > 0 {
				params := make(map[string]string, len(keyValuePairs))
				for _, keyValue := range keyValuePairs {
					splitted := strings.Split(keyValue, "=")
					if len(splitted) < 2 {
						continue
					}
					params[splitted[0]] = splitted[1]
				}
				log.Printf("Request params: %v", params)
			}
		}

		next.ServeHTTP(w, req)

		duration := time.Since(now)

		log.Printf("End request. Duration: %v.\n", duration)
	})
}

func TimeLoggerLocalMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		now := time.Now()
		bytesBody, _ := io.ReadAll(req.Body)

		log.Printf("Start request from IP - %s", req.RemoteAddr)
		log.Printf("Body: %s", string(bytesBody))

		next(w, req)

		duration := time.Since(now)

		log.Printf("End request. Duration: %v.\n", duration)
	}
}
