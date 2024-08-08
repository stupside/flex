package apps

import (
	"net/http"

	"flex.io/pkg/flux"
	"flex.io/pkg/rest"
)

type BodySchema struct {
	Rules      []flux.AppIngressRule `json:"rules"`
	AppName    string                `json:"application_name"`
	ChartURL   string                `json:"chart_url"`
	ChartName  string                `json:"chart_name"`
	TenantName string                `json:"tenant_name"`
	AppVersion string                `json:"application_version"`
}

type RespSchema struct {
}

func CreateAppCommand(c *flux.TenantController) rest.HttpCommand[BodySchema, RespSchema] {
	return func(r *http.Request, b *BodySchema) (rest.Result[RespSchema], rest.Result[rest.ErrorResult]) {
		fns := make([]flux.WithTenantAppFileTemplateFn, 0)

		for _, rule := range b.Rules {
			fns = append(fns, c.WithAppNetworkIngress(flux.DefaultIngressClass, flux.DefaultClusterIssuer, rule))
		}

		fns = append(fns, c.WithAppRelease(b.AppVersion))
		fns = append(fns, c.WithAppRepository(b.ChartURL))

		err := c.CreateApp(r.Context(), b.TenantName, b.AppName, fns)

		if err != nil {
			return nil, rest.Respond(http.StatusCreated, rest.ErrorResult{
				Message: err.Error(),
			})
		}

		return rest.Respond(http.StatusCreated, RespSchema{}), nil
	}
}
