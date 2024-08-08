package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"flex.io/pkg/git"
	"flex.io/pkg/git/apis/github/blob"
	"flex.io/pkg/git/apis/github/commit"
	"flex.io/pkg/git/apis/github/file"
	"flex.io/pkg/git/apis/github/ref"
	"flex.io/pkg/git/apis/github/tree"
)

// Constants for file modes and types
const (
	blobType = "blob"
	fileMode = "100644"
)

type GithubFilePusher struct {
	client *git.GitHttpClient
}

// Push pushes the specified files to the repository using the provided Git controller and commit.
func (g *GithubFilePusher) PushFile(ctx context.Context, c git.GitCommit, files ...git.GitFile) error {

	if len(files) == 1 {
		_, err := file.PushFileHandler(*g.client, files[0].Path)(ctx, &file.PushFileReq{
			Committer: c.Author,
			Message:   c.Message,
			Content:   base64.StdEncoding.EncodeToString([]byte(files[0].Content)),
		})

		if err != nil {
			return fmt.Errorf("failed to push file: %w", err)
		}

		return nil
	}

	getTreeResp, err := tree.GetTreeHandler(*g.client, g.client.Source.Branch)(ctx, &tree.GetTreeReq{})

	if err != nil {
		return fmt.Errorf("failed to get tree: %w", err)
	}

	createTreeReq := &tree.CreateTreeReq{
		Base: getTreeResp.Sha,
	}

	wg := sync.WaitGroup{}

	work := make(chan error, 1)

	for _, file := range files {

		wg.Add(1)
		go func() {

			encodedContent := base64.StdEncoding.EncodeToString([]byte(file.Content))

			blobResp, err := blob.CreateBlobHandler(*g.client)(ctx, &blob.CreateBlobReq{
				Content:  encodedContent,
				Encoding: blob.DefaultEncoding,
			})

			if err != nil {
				work <- fmt.Errorf("failed to create blob for file %s: %w", file.Path, err)
			}

			createTreeReq.AddBlob(file.Path, git.GitBlob{
				Type: blobType,
				Mode: fileMode,
				Sha:  blobResp.Sha,
			})

			defer wg.Done()
		}()
	}

	select {
	case err := <-work:
		return err
	default:
	}

	wg.Wait()

	createTreeResp, err := tree.CreateTreeHandler(*g.client)(ctx, createTreeReq)

	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	createCommitResp, err := commit.CreateCommitHandler(*g.client)(ctx, &commit.CreateCommitReq{
		GitCommit: c,
		Tree:      createTreeResp.Sha,
		Parents:   []string{},
	})

	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	_, err = ref.UpdateRefHandler(*g.client)(ctx, &ref.UpdateRefReq{
		Sha:         createCommitResp.Sha,
		FastForward: true,
	})

	if err != nil {
		return fmt.Errorf("failed to update ref: %w", err)
	}

	return nil
}

type GithubFileGetter struct {
	client *git.GitHttpClient
}

// Get retrieves the file content from the repository using the provided Git controller and path.
func (g *GithubFileGetter) GetFile(ctx context.Context, path string) (git.GitFile, error) {
	resp, err := file.GetFileHandler(*g.client, path)(ctx, &file.GetFileReq{})

	if err != nil {
		return git.GitFile{}, fmt.Errorf("failed to get file: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(resp.Content)
	if err != nil {
		return git.GitFile{}, fmt.Errorf("failed to decode content: %w", err)
	}

	return git.GitFile{
		Path:    resp.Path,
		Content: string(decodedContent),
	}, nil
}

type GithubTreeGetter struct {
	client *git.GitHttpClient
}

// Get retrieves the tree content from the repository using the provided Git controller and SHA.
func (g *GithubTreeGetter) GetTree(ctx context.Context, sha string) ([]git.GitBlobWithPath, error) {
	resp, err := tree.GetTreeHandler(*g.client, sha)(ctx, &tree.GetTreeReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	files := make([]git.GitBlobWithPath, len(resp.Tree))

	for i, file := range resp.Tree {

		files[i] = git.GitBlobWithPath{
			Path: file.Path,
			GitBlob: git.GitBlob{
				Sha:  sha,
				Mode: file.Mode,
				Type: file.Type,
			},
		}
	}

	return files, nil
}

type GithubService struct {
	GithubFilePusher
	GithubFileGetter
	GithubTreeGetter
}

func NewGithubService(c *git.GitHttpClient) git.GitService {
	return &GithubService{
		GithubFilePusher: GithubFilePusher{
			client: c,
		},
		GithubFileGetter: GithubFileGetter{
			client: c,
		},
		GithubTreeGetter: GithubTreeGetter{
			client: c,
		},
	}
}
