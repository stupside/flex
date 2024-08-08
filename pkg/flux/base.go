package flux

import (
	"flex.io/pkg/git"
)

// FluxController manages git rest related to tenants
type FluxController struct {
	author  git.GitAuthor
	service git.GitService
}

func NewFluxController(s git.GitService, a git.GitAuthor) *FluxController {
	return &FluxController{
		author:  a,
		service: s,
	}
}
