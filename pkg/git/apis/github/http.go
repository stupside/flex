package github

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"flex.io/pkg/git"
)

// Constants for HTTP requests to the GitHub API
const (
	GithubAPIBaseURL            = "https://api.github.com"
	GithubAcceptHeader          = "application/vnd.github.v3+json"
	GithubContentTypeHeader     = "application/json"
	GithubDefaultRequestTimeout = 10 * time.Second
)

// GithubHttpClientRoundTripper is a custom HTTP round tripper for the Git HTTP client
type GithubHttpClientRoundTripper struct {
	http.RoundTripper
	Bearer string
}

// RoundTrip adds the necessary headers to an HTTP request for the GitHub API
func (rt *GithubHttpClientRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !strings.HasPrefix(req.URL.String(), GithubAPIBaseURL) {
		return nil, fmt.Errorf("request not directed to Git API: %s", req.URL.String())
	}

	req.Header.Set("Accept", GithubAcceptHeader)
	req.Header.Set("Content-Type", GithubContentTypeHeader)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.Bearer))

	return rt.RoundTripper.RoundTrip(req)
}

// NewGitHttpClient creates a new Git controller with Bearer authentication
func NewGitHttpClient(bearer string, source git.GitSource, timeout time.Duration) *git.GitHttpClient {
	return &git.GitHttpClient{
		Source: source,
		HttpClient: &http.Client{
			Timeout: timeout,
			Transport: &GithubHttpClientRoundTripper{
				Bearer:       bearer,
				RoundTripper: http.DefaultTransport,
			},
		},
	}
}
