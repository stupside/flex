package git

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"
)

// GitFile represents a file to be pushed to the repository.
type GitFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// GitBlobOut represents a blob in the Git repository.
type GitBlob struct {
	Sha  string `json:"sha"`
	Mode string `json:"mode"`
	Type string `json:"type"`
}

type GitBlobWithPath struct {
	GitBlob
	Path string `json:"path"`
}

type FilePusher interface {
	PushFile(ctx context.Context, com GitCommit, files ...GitFile) error
}

type FileGetter interface {
	GetFile(ctx context.Context, path string) (GitFile, error)
}

type TreeGetter interface {
	GetTree(ctx context.Context, sha string) ([]GitBlobWithPath, error)
}

// GitFileTemplate represents a template to be rendered and pushed.
type GitFileTemplate struct {
	Path        string
	Data        interface{}
	Destination string
}

// PushTmpl pushes the specified template files to the repository using the provided Git controller and commit.
func PushTmpl(ctx context.Context, p FilePusher, path string, commit GitCommit, tmplFiles ...GitFileTemplate) error {
	files := make([]GitFile, len(tmplFiles))

	var buffer bytes.Buffer
	for i, tmpl := range tmplFiles {
		tmplParsed, err := template.ParseFiles(tmpl.Path)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tmpl.Path, err)
		}

		if err := tmplParsed.Execute(&buffer, tmpl.Data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", tmpl.Path, err)
		}

		files[i] = GitFile{
			Path:    tmpl.Destination,
			Content: buffer.String(),
		}

		log.Default().Println(files[i].Content)

		buffer.Reset()
	}

	return p.PushFile(ctx, commit, files...)
}

type GitService interface {
	FilePusher
	FileGetter
	TreeGetter
}
