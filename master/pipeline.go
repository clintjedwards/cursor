package master

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	cursorPlugin "github.com/clintjedwards/cursor/plugin"
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

func (master *CursorMaster) runPipeline(pipelineID string) error {

	pluginPath := fmt.Sprintf("%s/%s", master.config.Master.PluginDirectoryPath, pipelineID)

	// We create a client so that we can communicate with plugins aka pipelines
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  cursorPlugin.Handshake,
		Plugins:          master.pluginMap,
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// // Connect via RPC
	// rpcClient, err := client.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // Request the plugin
	// raw, err := rpcClient.Dispense("kv_grpc")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }
	// // We should have a KV store now! This feels like a normal interface
	// // implementation but is in fact over an RPC connection.
	// kv := raw.(shared.KV)
	// os.Args = os.Args[1:]
	// switch os.Args[0] {
	// case "get":
	// 	result, err := kv.Get(os.Args[1])
	// 	if err != nil {
	// 		fmt.Println("Error:", err.Error())
	// 		os.Exit(1)
	// 	}

	// 	fmt.Println(string(result))

	// case "put":
	// 	err := kv.Put(os.Args[1], []byte(os.Args[2]))
	// 	if err != nil {
	// 		fmt.Println("Error:", err.Error())
	// 		os.Exit(1)
	// 	}

	// default:
	// 	fmt.Printf("Please only use 'get' or 'put', given: %q", os.Args[0])
	// 	os.Exit(1)
	// }
	// os.Exit(0)
	return nil
}
