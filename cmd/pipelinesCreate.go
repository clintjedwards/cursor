package cmd

import (
	"log"

	"github.com/clintjedwards/cursor/config"
	"github.com/spf13/cobra"
)

var cmdPipelinesCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a single pipeline",
	Long: `A pipeline is a grouping of tasks to acheive an end goal.
Pipelines are written according to cursor documentation and then compiled and added as a module.

Pipelines are first downloaded from git repositories and require a git url to work`,
	Args: cobra.MinimumNArgs(1),
	Run:  runPipelinesCreateCmd,
}

func runPipelinesCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]
	description, _ := cmd.Flags().GetString("description")
	gitURL, _ := cmd.Flags().GetString("gitURL")
	gitBranch, _ := cmd.Flags().GetString("gitBranch")

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration")
		
	}
}

func init() {
	var description string
	cmdPipelinesCreate.Flags().StringVarP(&description, "description", "d", "", "long form description of pipeline")

	var gitURL string
	cmdPipelinesCreate.Flags().StringVarP(&gitURL, "git_url", "u", "", "url of git repository")

	var gitBranch string
	cmdPipelinesCreate.Flags().StringVarP(&gitBranch, "git_branch", "b", "master", "git branch name")

	cmdPipeline.AddCommand(cmdPipelinesCreate)
}
