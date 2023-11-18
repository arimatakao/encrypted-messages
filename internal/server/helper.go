package server

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, output any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(output)
}

func WriteStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func WriteJSON(w http.ResponseWriter, statusCode int, input any) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(input)
}
