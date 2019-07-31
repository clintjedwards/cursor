package master

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/clintjedwards/cursor/api"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	golangBinaryName = "go"
)

// executeCmd wraps a context around the command and executes it.
func executeCmd(path string, args []string, env []string, dir string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = env
	cmd.Dir = dir

	// Execute command
	return cmd.CombinedOutput()
}

// cloneRepository is used to download plugin repositories
func (master *CursorMaster) cloneRepository(pluginID string, repo *api.GitRepo) error {

	repoPath := fmt.Sprintf("%s/%s", master.config.Master.RepoDirectoryPath, pluginID)

	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:           repo.Url,
		ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
	})

	return err
}

// buildPlugin attempts to compile and store a specified plugin
func (master *CursorMaster) buildPlugin(pluginID string) error {

	repoPath := fmt.Sprintf("%s/%s", master.config.Master.RepoDirectoryPath, pluginID)
	binarypath := fmt.Sprintf("%s/%s", master.config.Master.PluginDirectoryPath, pluginID)

	buildArgs := []string{"build", "-o", binarypath}

	golangBinaryPath, err := exec.LookPath(golangBinaryName)
	if err != nil {
		return err
	}

	_, err = executeCmd(golangBinaryPath, buildArgs, nil, repoPath)
	if err != nil {
		return err
	}

	return nil
}
