package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"flex.io/api/tenants/apps"
	"flex.io/api/tenants/tenants"
	"flex.io/pkg/flux"
	"flex.io/pkg/git"
	"flex.io/pkg/git/apis/github"
	"flex.io/pkg/rest"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	clusterToken, found := os.LookupEnv("GITHUB_CLUSTER_TOKEN")

	if !found {
		panic("GITHUB_CLUSTER_TOKEN not found")
	}

	tenantToken, found := os.LookupEnv("GITHUB_TENANT_TOKEN")

	if !found {
		panic("GITHUB_TENANT_TOKEN not found")
	}

	gitAuthor := git.GitAuthor{
		Name:  "flex",
		Email: "flex@flex.io",
	}

	mux := http.NewServeMux()

	{
		gitClient := github.NewGitHttpClient(clusterToken, git.GitSource{
			Repo:   "devops",
			Owner:  "stupside",
			Branch: "main",
		}, github.GithubDefaultRequestTimeout)

		gitService := github.NewGithubService(gitClient)

		clusterController := flux.NewClusterController(gitService, gitAuthor)

		rest.RegisterCommand(mux, []string{
			http.MethodPost}, "/cluster/tenant", tenants.RegisterTenantCommand(clusterController))
	}

	{
		gitClient := github.NewGitHttpClient(tenantToken, git.GitSource{
			Repo:   "devops-dev-tenant",
			Owner:  "stupside",
			Branch: "main",
		}, github.GithubDefaultRequestTimeout)

		gitService := github.NewGithubService(gitClient)

		tenantController := flux.NewTenantController(gitService, gitAuthor)

		rest.RegisterCommand(mux, []string{
			http.MethodPost,
		}, "/tenants/app", apps.CreateAppCommand(tenantController))
	}

	err := http.ListenAndServe(":8080", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server: %s", err)
	}
}
