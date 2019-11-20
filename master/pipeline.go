package master

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	cursorPlugin "github.com/clintjedwards/cursor/plugin"
	proto "github.com/clintjedwards/cursor/plugin/proto"
	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-plugin"
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
// TODO: should clean up after itself if it fails
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

func (master *CursorMaster) getPipelineInfo(pipelineID string) (proto.GetPipelineInfoResponse, error) {

	pluginPath := fmt.Sprintf("%s/%s", master.config.Master.PluginDirectoryPath, pipelineID)

	// We create a client so that we can communicate with plugins aka pipelines
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  cursorPlugin.Handshake,
		Plugins:          master.pluginMap,
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return proto.GetPipelineInfoResponse{}, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pipelineID)
	if err != nil {
		return proto.GetPipelineInfoResponse{}, err
	}

	pipeline := raw.(cursorPlugin.PipelineDefinition)
	response, err := pipeline.GetPipelineInfo(&proto.GetPipelineInfoRequest{})
	if err != nil {
		return proto.GetPipelineInfoResponse{}, err
	}

	return *response, nil
}

// runTasks checks the plugins loaded for the specific task to run
func (master *CursorMaster) runTasks(pipelineID, taskID string) error {

	pluginPath := fmt.Sprintf("%s/%s", master.config.Master.PluginDirectoryPath, pipelineID)

	// We create a client so that we can communicate with plugins aka pipelines
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  cursorPlugin.Handshake,
		Plugins:          master.pluginMap,
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pipelineID)
	if err != nil {
		return err
	}

	pipeline := raw.(cursorPlugin.PipelineDefinition)
	message, err := pipeline.ExecuteTask(&proto.ExecuteTaskRequest{
		Id: taskID,
	})
	if err != nil {
		return err
	}

	// TODO: recusively Executetask and paralize with waitgroups intelligently
	// Possibly also make it so that if one child fails the entire pipeline is
	// stopped and marked as a failure while still recording the other task.
	fmt.Println(message)

	return nil
}

func (master *CursorMaster) addPlugin(id string) {
	master.pluginMapMutex.Lock()
	defer master.pluginMapMutex.Unlock()
	master.pluginMap[id] = &cursorPlugin.CursorPlugin{}
}

func (master *CursorMaster) deletePlugin(id string) {
	master.pluginMapMutex.Lock()
	defer master.pluginMapMutex.Unlock()
	delete(master.pluginMap, id)
}
