package master

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	getter "github.com/hashicorp/go-getter"
)

const (
	golangBinaryName = "go"
)

// executeCmd wraps a context around a given command and executes it.
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

// getRepository is used to download plugin repositories. Uses hashicorp's go-getter
// to provide the ability to download from most string URLs without writing a custom
// implementation to do so.
// should be able to download from most common sources. (eg: git, http, mercurial)
// See (https://github.com/hashicorp/go-getter#url-format) for more information on how to form input
func (master *CursorMaster) getRepository(pluginID string, repoURL string) error {

	cursorRepoPath := fmt.Sprintf("%s/%s", master.config.Master.RepoDirectoryPath, pluginID)

	err := getter.GetAny(cursorRepoPath, repoURL)
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
