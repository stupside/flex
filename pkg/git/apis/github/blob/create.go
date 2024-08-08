package blob

import (
	"flex.io/pkg/git"
)

const (
	DefaultEncoding = "base64"
)

// CreateBlobReq represents the request body for creating a blob in a repository.
type CreateBlobReq struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// CreateBlobResp represents the response body for creating a blob in a repository.
type CreateBlobResp struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

// CreateBlob creates a new blob in the specified repository with the provided content.
func CreateBlobHandler(c git.GitHttpClient) git.GitHttpHandler[CreateBlobReq, CreateBlobResp] {
	return git.NewHttpHandler[CreateBlobReq, CreateBlobResp](c, git.GitEndpoint{
		Path:   "/git/blobs",
		Method: "POST",
	})
}
