package templates

import (
	"fmt"
	"os"
)

const (
	TenantRbacTmplPath       = "templates/tenant/rbac.yml.tmpl"
	TenantReleaseTmplPath    = "templates/tenant/release.yml.tmpl"
	TenantRepositoryTmplPath = "templates/tenant/repository.yml.tmpl"
)

const (
	TenantAppReleasePath        = "templates/app/release.yml.tmpl"
	TenantAppRepositoryPath     = "templates/app/repository.yml.tmpl"
	TenantAppNetworkIngressPath = "templates/app/network/ingress.yml.tmpl"
)

func GetPathToTmpl(path string) string {
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/pkg/api/%s", dir, path)
}
