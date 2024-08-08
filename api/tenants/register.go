package tenants

import (
	"net/http"

	"flex.io/pkg/flux"
	"flex.io/pkg/rest"
)

type BodySchema struct {
	TenantName    string `json:"tenantName"`
	ClusterName   string `json:"clusterName"`
	RepositoryUrl string `json:"repositoryUrl"`
}

type RespSchema struct {
}

func RegisterTenantCommand(c *flux.ClusterController) rest.HttpCommand[BodySchema, RespSchema] {
	return func(r *http.Request, b *BodySchema) (rest.Result[RespSchema], rest.Result[rest.ErrorResult]) {

		fns := []flux.TenantFileTemplate{
			c.WithTenantRbac(),
			c.WithTenantRelease(),
			c.WithTenantRepository(b.RepositoryUrl),
		}

		err := c.CreateTenant(r.Context(), b.ClusterName, b.TenantName, fns)

		if err != nil {
			return nil, rest.Result[rest.ErrorResult](rest.Respond(http.StatusInternalServerError, rest.ErrorResult{Message: err.Error()}))
		}

		return rest.Result[RespSchema](rest.Respond(http.StatusCreated, RespSchema{})), nil
	}
}
