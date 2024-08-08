package flux

import (
	"context"
	"fmt"

	"flex.io/pkg/flux/templates"
	"flex.io/pkg/git"
)

const (
	DefaultIngressClass  = "nginx"
	DefaultClusterIssuer = "letsencrypt-prod"
)

// TenantController manages tenant applications using Git.
type TenantController struct {
	*FluxController
}

// NewTenantController creates a new FluxTenantController with the given Git controller.
func NewTenantController(s git.GitService, a git.GitAuthor) *TenantController {
	return &TenantController{
		FluxController: NewFluxController(s, a),
	}
}

// WithAppRepository defines the data for app repository templates.
type WithAppRepository struct {
	ChartURL  string
	ChartName string
}

// WithAppRepository returns a function that creates a Git template for app repository configuration.
func (c *TenantController) WithAppRepository(chartUrl string) WithTenantAppFileTemplateFn {
	return func(tenant, name, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: WithAppRepository{
				ChartName: name,
				ChartURL:  chartUrl,
			},
			Path:        templates.GetPathToTmpl(templates.TenantAppRepositoryPath),
			Destination: fmt.Sprintf("%s/repository.yml", destination),
		}
	}
}

// WithAppRelease defines the data for app release templates.
type WithAppRelease struct {
	Tenant       string
	Application  string
	ChartName    string
	ChartVersion string
}

// WithAppRelease returns a function that creates a Git template for app release configuration.
func (c *TenantController) WithAppRelease(version string) WithTenantAppFileTemplateFn {
	return func(tenant, name, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: WithAppRelease{
				Application:  name,
				ChartName:    name,
				Tenant:       tenant,
				ChartVersion: version,
			},
			Path:        templates.GetPathToTmpl(templates.TenantAppReleasePath),
			Destination: fmt.Sprintf("%s/release.yml", destination),
		}
	}
}

// AppIngressRule defines an ingress rule with domain, subdomain, and paths.
type AppIngressRule struct {
	Domain    string
	Subdomain string
	Paths     []AppIngressRulePath
}

func NewAppIngreeRule(domain, subdomain string) AppIngressRule {
	return AppIngressRule{
		Domain:    domain,
		Subdomain: subdomain,
	}
}

// AppIngressRulePath represents the path configuration for an ingress rule.
type AppIngressRulePath struct {
	Port    int16
	Service string
}

func (r *AppIngressRule) AddPath(service string, port int16) {
	r.Paths = append(r.Paths, AppIngressRulePath{
		Port:    port,
		Service: service,
	})
}

// WithAppNetworkIngress defines the data for app network ingress templates.
type WithAppNetworkIngress struct {
	Application   string
	IngressClass  string
	ClusterIssuer string
	Rules         []AppIngressRule
}

// WithAppNetworkIngress returns a function that creates a Git template for network ingress configuration.
func (c *TenantController) WithAppNetworkIngress(ingressClass, clusterIssuer string, rules ...AppIngressRule) WithTenantAppFileTemplateFn {
	return func(tenant, name, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: WithAppNetworkIngress{
				Application:   name,
				Rules:         rules,
				IngressClass:  ingressClass,
				ClusterIssuer: clusterIssuer,
			},
			Path:        templates.GetPathToTmpl(templates.TenantAppNetworkIngressPath),
			Destination: fmt.Sprintf("%s/network/ingress.yml", destination),
		}
	}
}

// WithTenantAppFileTemplateFn is a type for functions that create Git templates for applications.
type WithTenantAppFileTemplateFn func(tenant, name, namespace string) git.GitFileTemplate

// CreateApp creates a new tenant application by applying multiple templates.
func (c *TenantController) CreateApp(ctx context.Context, tenant string, name string, fns []WithTenantAppFileTemplateFn) error {
	destination := fmt.Sprintf("apps/%s", name)

	opts := make([]git.GitFileTemplate, len(fns))
	for i, fn := range fns {
		opts[i] = fn(tenant, name, destination)
	}

	com := git.GitCommit{
		Author:  c.author,
		Message: "feat: create app",
	}

	if err := git.PushTmpl(ctx, c.service, destination, com, opts...); err != nil {
		return fmt.Errorf("failed to create app '%s': %w", name, err)
	}

	return nil
}
