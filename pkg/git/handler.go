package git

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// GitEndpoint represents an HTTP endpoint with a URL and method.
type GitEndpoint struct {
	Path   string
	Method string
}

// GitHttpHandler is a type alias for a function that handles HTTP requests and returns a response or an error.
type GitHttpHandler[TData any, TResp any] func(ctx context.Context, req *TData) (*TResp, error)

// NewHttpHandler creates a new HTTP handler function for the given client and endpoint.
func NewHttpHandler[TData any, TResp any](c GitHttpClient, e GitEndpoint) GitHttpHandler[TData, TResp] {
	return func(ctx context.Context, d *TData) (*TResp, error) {
		return call[TData, TResp](ctx, c, e.Method, fmt.Sprintf("%s/%s", c.Source, e.Path), d)
	}
}

// call sends an HTTP call and decodes the response.
func call[TData any, TResp any](ctx context.Context, c GitHttpClient, method string, url string, d *TData) (*TResp, error) {
	// Marshal the request body into JSON.
	marshaled, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP request with the context and marshaled body.
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(marshaled))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the HTTP request using the provided client.
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("failed to close response body: %v", cerr)
		}
	}()

	log.Printf("executed request: %d %s %s", resp.StatusCode, method, url)

	// Check for non-2xx HTTP status codes.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		payload, _ := json.Marshal(d)

		body, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to dump response: %w", err)
		}

		return nil, fmt.Errorf("request returned non-2xx status: %d\nRequest payload: %s\nResponse: %s", resp.StatusCode, payload, body)
	}

	// Decode the response body into the expected response type.
	var res TResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &res, nil
}
