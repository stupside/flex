package tree

import (
	"flex.io/pkg/git"
)

type CreateTreeReqItem struct {
	git.GitBlob
	Path string `json:"path"`
}

// CreateTreeReq represents the request body for creating a Git tree.
type CreateTreeReq struct {
	Base string              `json:"base_tree"`
	Tree []CreateTreeReqItem `json:"tree"`
}

func (r *CreateTreeReq) AddBlob(path string, blob git.GitBlob) {
	r.Tree = append(r.Tree, CreateTreeReqItem{
		Path:    path,
		GitBlob: blob,
	})
}

// CreateTreeResp represents the response structure for creating a Git tree.
type CreateTreeResp struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

// CreateTree creates a new tree in the repository using the provided blobs.
func CreateTreeHandler(c git.GitHttpClient) git.GitHttpHandler[CreateTreeReq, CreateTreeResp] {
	return git.NewHttpHandler[CreateTreeReq, CreateTreeResp](c, git.GitEndpoint{
		Path:   "git/trees",
		Method: "POST",
	})
}
