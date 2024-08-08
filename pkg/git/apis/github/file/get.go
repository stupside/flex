package file

import (
	"fmt"

	"flex.io/pkg/git"
)

type GetFileReq struct{}

type GetFileResp struct {
	Sha      string `json:"sha"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GetFileHandler(c git.GitHttpClient, path string) git.GitHttpHandler[GetFileReq, GetFileResp] {
	return git.NewHttpHandler[GetFileReq, GetFileResp](c, git.GitEndpoint{
		Path:   fmt.Sprintf("/contents/%s", path),
		Method: "GET",
	})
}
