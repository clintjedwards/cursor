package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/clintjedwards/cursor/api"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

var cmdPipelinesGet = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details about a single pipeline",
	Args:  cobra.MinimumNArgs(1),
	Run:   runPipelinesGetCmd,
}

func runPipelinesGetCmd(cmd *cobra.Command, args []string) {
	id := args[0]

	conn := initClientConn()
	client := api.NewCursorMasterClient(conn)
	ctx := generateClientContext()

	pipelineProto, err := client.GetPipeline(ctx, &api.GetPipelineRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not get pipeline: %v", err)
	}

	pipeline := pipelineProto.Pipeline

	fmt.Printf("ID: %s\n", pipeline.Id)
	fmt.Printf("Name: %s\n", pipeline.Name)
	fmt.Printf("Description: %s\n", pipeline.Description)
	fmt.Printf("Repository URL: %s\n", pipeline.RepositoryUrl)
	fmt.Printf("Created: %s\n", humanize.Time(time.Unix(pipeline.Created, 0)))
	fmt.Printf("Modified: %s\n", humanize.Time(time.Unix(pipeline.Modified, 0)))
	fmt.Printf("Last Compiled: %s\n", humanize.Time(time.Unix(pipeline.LastCompiled, 0)))
	fmt.Printf("Root Task: [%s] %s\n", pipeline.RootTaskId, pipeline.Tasks[pipeline.RootTaskId].Name)

	taskGraph := treeprint.New()
	taskGraph = populateTree(pipeline.RootTaskId, pipeline.Tasks, taskGraph)
	fmt.Printf("\nTask Graph:\n")
	fmt.Println(taskGraph.String())
}

func init() {
	cmdPipelines.AddCommand(cmdPipelinesGet)
}

// populateTree is a recursive function that parses a graph map and turns it into
// tree graph
func populateTree(nodeID string, taskMap map[string]*api.Task, taskTree treeprint.Tree) treeprint.Tree {
	if _, ok := taskMap[nodeID]; !ok {
		return taskTree
	}

	if len(taskMap[nodeID].Children) == 0 {
		return taskTree.AddNode(taskMap[nodeID].Name)
	}

	taskTree = taskTree.AddBranch(taskMap[nodeID].Name)

	for _, child := range taskMap[nodeID].Children {
		taskTree = populateTree(child, taskMap, taskTree)
	}

	return taskTree
}
