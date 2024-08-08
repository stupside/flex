package file

import (
	"fmt"

	"flex.io/pkg/git"
)

// PushFileReq represents the request body for pushing a file to the repository.
type PushFileReq struct {
	Message   string        `json:"message"`
	Content   string        `json:"content"`
	Committer git.GitAuthor `json:"committer"`
}

type PushFileResp struct{}

func PushFileHandler(c git.GitHttpClient, path string) git.GitHttpHandler[PushFileReq, PushFileReq] {
	return git.NewHttpHandler[PushFileReq, PushFileReq](c, git.GitEndpoint{
		Path:   fmt.Sprintf("/contents/%s", path),
		Method: "PUT",
	})
}
