package ci

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v3"
)

func NewWorkspaceFromGit(root, url, branch string) (*workspaceImpl, error) {
	dir, err := createTempDir(root)
	if err != nil {
		return nil, err
	}

	ref, err := cloneRepoAndGetRef(url, dir, branch)
	if err != nil {
		return nil, err
	}

	return &workspaceImpl{
		branch: branch,
		commit: ref.Hash().String(),
		dir:    dir,
		env:    []string{},
	}, nil
}

func createTempDir(root string) (string, error) {
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return "", err
	}

	dir, err := os.MkdirTemp(root, "workspace")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func cloneRepoAndGetRef(url, dir, branch string) (*plumbing.Reference, error) {
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             1,
	})

	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return ref, nil
}

type workspaceImpl struct {
	branch string
	commit string
	dir    string
	env    []string
}

func (w *workspaceImpl) Branch() string {
	return w.branch
}

func (w *workspaceImpl) Commit() string {
	return w.commit
}

func (w *workspaceImpl) Dir() string {
	return w.dir
}

func (w *workspaceImpl) Env() []string {
	return w.env
}

func (w *workspaceImpl) LoadPipeline() (*Pipeline, error) {
	data, err := os.ReadFile(filepath.Join(w.dir, "build", "flow-ci.yaml"))
	if err != nil {
		return nil, err
	}

	var pipeline Pipeline
	err = yaml.Unmarshal(data, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}
