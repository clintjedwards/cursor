package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/client"
	"github.com/clintjedwards/cursor/config"
	"github.com/spf13/cobra"
)

var cmdPipelinesList = &cobra.Command{
	Use:   "list",
	Short: "list all available pipelines",
	Long:  ``,
	Run:   runPipelinesListCmd,
}

func runPipelinesListCmd(cmd *cobra.Command, args []string) {

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration")
	}

	hostPortTuple := strings.Split(config.Master.URL, ":")

	cursorClient := client.CursorClient{}
	err = cursorClient.Connect(hostPortTuple[0], hostPortTuple[1])
	if err != nil {
		log.Fatalf("could not connect to host: %v", err)
	}
	defer cursorClient.Close()

	pipelines, err := cursorClient.ListPipelines(&api.ListPipelinesRequest{})
	if err != nil {
		log.Fatalf("could not list pipelines: %v", err)
	}

	for key, value := range pipelines.Pipelines {
		fmt.Printf("%s :: %s\n", key, value.String())
	}
}

func init() {
	cmdPipeline.AddCommand(cmdPipelinesList)
}
