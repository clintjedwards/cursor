package master

import (
	"github.com/clintjedwards/cursor/api"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func cloneRepository(repo *api.GitRepo) error {

	_, err := git.PlainClone("/tmp/test", false, &git.CloneOptions{
		URL:           repo.Url,
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
	})

	return err
}

func executeBuild() error {

	return nil
}
