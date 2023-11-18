package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CtxKey string

func NewRequestContext(reqCtx context.Context, userId string) (context.Context, error) {
	ctxKey := CtxKey("user_id")
	ctx := context.WithValue(reqCtx, ctxKey, userId)
	return ctx, nil
}

func GetUserId(r *http.Request) (string, error) {
	key := CtxKey("user_id")

	val := r.Context().Value(key)
	if val == "" {
		return "", fmt.Errorf("key not found in ctx")
	}

	userId, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("can't convert value of ctx to string")
	}

	return userId, nil
}

func ReadJSON(r *http.Request, output any) error {
	decoder := json.NewDecoder(r.Body)
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
