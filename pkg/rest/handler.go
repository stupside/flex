package rest

import (
	"encoding/json"
	"net/http"
)

// HttpCommand represents a handler function that returns an HttpResult.
type HttpCommand[TReqSchema any, TRespSchema any] func(r *http.Request, payload *TReqSchema) (Result[TRespSchema], Result[ErrorResult])

// RegisterCommand registers an HTTP handler with specified middlewares.
func RegisterCommand[TReqSchema any, TRespSchema any](
	mux *http.ServeMux, methods []string, path string,
	command HttpCommand[TReqSchema, TRespSchema],
	middlewares ...Middleware[TReqSchema, TRespSchema],
) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is allowed
		if !isMethodAllowed(r.Method, methods) {
			Respond(http.StatusMethodNotAllowed, ErrorResult{Message: "method not allowed"})(w)
			return
		}

		// Decode the request body
		payload, err := decodeJSON[TReqSchema](r)

		if err != nil {
			Respond(http.StatusBadRequest, ErrorResult{Message: "invalid params"})(w)
			return
		}
		defer r.Body.Close()

		// Apply middlewares in reverse order
		final := command
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](r, final)
		}

		// Execute the final handler
		succRes, errRes := final(r, payload)
		if errRes != nil {
			errRes(w)
		} else {
			succRes(w)
		}
	})
}

// isMethodAllowed checks if a request method is allowed.
func isMethodAllowed(method string, allowedMethods []string) bool {
	for _, m := range allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// decodeJSON is a helper function for decoding JSON with better error handling.
func decodeJSON[T any](r *http.Request) (*T, error) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
