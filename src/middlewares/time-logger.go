package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func TimeLoggerCommonMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		now := time.Now()

		bytesBody, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bytesBody))

		log.Printf("Start request from IP - %s", req.RemoteAddr)
		log.Printf("Body: %s", string(bytesBody))

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
