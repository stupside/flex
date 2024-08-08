package tree

import (
	"fmt"

	"flex.io/pkg/git"
)

type GetTreeReq struct{}

type getTreeReqTreeItem struct {
	CreateTreeReqItem
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

// GetTreeResp represents the response structure for getting a Git tree.
type GetTreeResp struct {
	Sha  string               `json:"sha"`
	Tree []getTreeReqTreeItem `json:"tree"`
}

// GetTreeHandler returns a GitActionHandler for getting the current Git tree.
func GetTreeHandler(c git.GitHttpClient, sha string) git.GitHttpHandler[GetTreeReq, GetTreeResp] {
	return git.NewHttpHandler[GetTreeReq, GetTreeResp](c, git.GitEndpoint{
		Path:   fmt.Sprintf("/git/trees/%s", sha),
		Method: "GET",
	})
}
