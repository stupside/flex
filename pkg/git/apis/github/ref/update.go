package ref

import (
	"fmt"

	"flex.io/pkg/git"
)

// UpdateRefReq represents the request body for updating a Git reference.
type UpdateRefReq struct {
	Sha         string `json:"sha"`
	FastForward bool   `json:"force"`
}

// UpdateRefResp represents the response structure for updating a Git reference.
type UpdateRefResp struct{}

// UpdateRefHandler creates a new handler for updating a Git reference in the repository.
func UpdateRefHandler(c git.GitHttpClient) git.GitHttpHandler[UpdateRefReq, UpdateRefResp] {
	return git.NewHttpHandler[UpdateRefReq, UpdateRefResp](c, git.GitEndpoint{
		Path:   fmt.Sprintf("/git/refs/heads/%s", c.Source.Branch),
		Method: "PATCH",
	})
}
