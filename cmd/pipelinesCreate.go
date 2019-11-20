package cmd

import (
	"log"

	"github.com/clintjedwards/cursor/api"
	"github.com/spf13/cobra"
)

var cmdPipelinesCreate = &cobra.Command{
	Use:   "create <name> <repository_url>",
	Short: "Create a single Pipeline",
	Long: `A pipeline is a collection of tasks that accomplishes some end goal.
To create and start using a pipeline you must first write the pipeline the
cursor-sdk.

The repository URL should point to the URL where your plugin code is hosted.
For more information on what formats repository URL accepts see:
https://github.com/hashicorp/go-getter#supported-protocols-and-detectors

ex: cursor pipelines create cursor-test github.com/clintjedwards/cursor-test`,
	Args: cobra.MinimumNArgs(2),
	Run:  runPipelinesCreateCmd,
}

func runPipelinesCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]
	repositoryURL := args[1]
	description, _ := cmd.Flags().GetString("description")

	conn := initClientConn()
	client := api.NewCursorMasterClient(conn)
	ctx := generateClientContext()

	_, err := client.CreatePipeline(ctx, &api.CreatePipelineRequest{
		Name:          name,
		Description:   description,
		RepositoryUrl: repositoryURL,
	})
	if err != nil {
		log.Fatalf("could not create pipeline: %v", err)
	}
}

func init() {
	var description string
	cmdPipelinesCreate.Flags().StringVarP(&description,
		"description", "d", "", "description of pipeline's purpose and other details")

	cmdPipelines.AddCommand(cmdPipelinesCreate)
}
