package commit

import (
	"flex.io/pkg/git"
)

// CreateCommitReq represents the request body for creating a Git commit.
type CreateCommitReq struct {
	git.GitCommit
	Tree    string   `json:"tree"`
	Parents []string `json:"parents"`
}

// CreateCommitResp represents the response body when a Git commit is created.
type CreateCommitResp struct {
	Sha string `json:"sha"`
}

func CreateCommitHandler(c git.GitHttpClient) git.GitHttpHandler[CreateCommitReq, CreateCommitResp] {
	return git.NewHttpHandler[CreateCommitReq, CreateCommitResp](c, git.GitEndpoint{
		Path:   "/git/commits",
		Method: "POST",
	})
}
