package git

import (
	"net/http"
)

// GitAuthor represents the information of a commit author
type GitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GitCommit represents a Git commit with an author and a message
type GitCommit struct {
	Author  GitAuthor `json:"author"`
	Message string    `json:"message"`
}

// GitSource represents the Git source including the repository, branch, and owner
type GitSource struct {
	Repo   string `json:"repo"`
	Owner  string `json:"owner"`
	Branch string `json:"branch"`
}

// HttpClient is an interface for making HTTP requests
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitHttpClient struct {
	HttpClient
	Source GitSource
}
