package server

import (
	"encoding/json"
	"net/http"
)

func WriteStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func WriteJSON(w http.ResponseWriter, statusCode int, body any) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(body)
}
