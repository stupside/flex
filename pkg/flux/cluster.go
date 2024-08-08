package flux

import (
	"context"
	"fmt"

	"flex.io/pkg/flux/templates"
	"flex.io/pkg/git"
)

type ClusterController struct {
	*FluxController
}

// NewFluxController initializes a new FluxController instance
func NewClusterController(s git.GitService, a git.GitAuthor) *ClusterController {
	return &ClusterController{
		FluxController: NewFluxController(s, a),
	}
}

// TenantFileTemplate defines a function signature for creating Git templates
type TenantFileTemplate func(tenant string, destination string) git.GitFileTemplate

// rbac represents data needed for RBAC template generation
type rbac struct {
	Tenant string
}

// WithTenantRbac creates a function for generating RBAC templates
func (c *ClusterController) WithTenantRbac() TenantFileTemplate {
	return func(tenant, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: rbac{
				Tenant: tenant,
			},
			Path:        templates.GetPathToTmpl(templates.TenantRbacTmplPath),
			Destination: fmt.Sprintf("%s/rbac.yml", destination),
		}
	}
}

// withTenantRepository represents data needed for Repository template generation
type repository struct {
	Tenant        string
	RepositoryUrl string
}

// WithTenantRepository creates a function for generating Repository templates
func (c *ClusterController) WithTenantRepository(url string) TenantFileTemplate {
	return func(tenant, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: repository{
				Tenant:        tenant,
				RepositoryUrl: url,
			},
			Path:        templates.GetPathToTmpl(templates.TenantRepositoryTmplPath),
			Destination: fmt.Sprintf("%s/repository.yml", destination),
		}
	}
}

// withTenantRelease represents data needed for Release template generation
type release struct {
	Tenant string
}

// WithTenantRelease creates a function for generating Release templates
func (c *ClusterController) WithTenantRelease() TenantFileTemplate {
	return func(tenant, destination string) git.GitFileTemplate {
		return git.GitFileTemplate{
			Data: release{
				Tenant: tenant,
			},
			Path:        templates.GetPathToTmpl(templates.TenantReleaseTmplPath),
			Destination: fmt.Sprintf("%s/release.yml", destination),
		}
	}
}

// CreateTenant creates a new tenant by applying a set of template functions
func (c *ClusterController) CreateTenant(ctx context.Context, cluster string, tenant string, fns []TenantFileTemplate) error {
	destination := fmt.Sprintf("/environments/%s/tenants/%s", cluster, tenant)

	opts := make([]git.GitFileTemplate, len(fns))
	for i, fn := range fns {
		opts[i] = fn(tenant, destination)
	}

	com := git.GitCommit{
		Message: "feat: create tenant",
		Author: git.GitAuthor{
			Name:  c.author.Name,
			Email: c.author.Email,
		},
	}

	if err := git.PushTmpl(ctx, c.service, destination, com, opts...); err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	return nil
}
