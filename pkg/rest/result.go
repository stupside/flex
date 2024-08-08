package rest

import (
	"encoding/json"
	"net/http"
)

// Result represents a function that writes a typed HTTP response.
type Result[TRespSchema any] func(w http.ResponseWriter)

// ErrorResult represents an error message in JSON format.
type ErrorResult struct {
	Message string `json:"message"`
}

// Respond creates an HttpResult that writes a status code and JSON-encoded data of type TRespSchema.
func Respond[TRespSchema any](statusCode int, data TRespSchema) Result[TRespSchema] {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
