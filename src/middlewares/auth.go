package middlewares

import (
	"encoding/json"
	"net/http"
)

func AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "secret-key" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(401)
			res := struct {
				Message string `json:"message"`
			}{Message: "Authorization header required."}
			json.NewEncoder(w).Encode(res)
			return
		}
		next.ServeHTTP(w, r)
	}
}
