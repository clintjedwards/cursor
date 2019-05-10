package cmd

import (
	"log"
	"strings"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/client"
	"github.com/clintjedwards/cursor/config"
	"github.com/spf13/cobra"
)

var cmdPipelinesCreate = &cobra.Command{
	Use:   "create <name> <git_url>",
	Short: "Create a single pipeline",
	Long: `A pipeline is a grouping of tasks to acheive an end goal.
Pipelines are written according to cursor documentation and then compiled and added as a module.

Pipelines are first downloaded from git repositories and require a git url to work`,
	Args: cobra.MinimumNArgs(1),
	Run:  runPipelinesCreateCmd,
}

func runPipelinesCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]
	gitURL := args[1]
	description, _ := cmd.Flags().GetString("description")
	gitBranch, _ := cmd.Flags().GetString("gitBranch")

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

	_, err = cursorClient.CreatePipeline(&api.CreatePipelineRequest{
		Name:        name,
		Description: description,
		GitRepo: &api.GitRepo{
			Url:    gitURL,
			Branch: gitBranch,
		},
	})
	if err != nil {
		log.Fatalf("could not create pipeline: %v", err)
	}
}

func init() {
	var description string
	cmdPipelinesCreate.Flags().StringVarP(&description, "description", "d", "", "long form description of pipeline")

	var gitBranch string
	cmdPipelinesCreate.Flags().StringVarP(&gitBranch, "git_branch", "b", "master", "git branch name")

	cmdPipeline.AddCommand(cmdPipelinesCreate)
}
