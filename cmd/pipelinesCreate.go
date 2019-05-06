package cmd

import (
	"github.com/spf13/cobra"
)

var cmdPipelinesCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a single pipeline",
	Args:  cobra.MinimumNArgs(1),
	Run:   runPipelinesCreateCmd,
}

func runPipelinesCreateCmd(cmd *cobra.Command, args []string) {

}

func init() {
	var description string
	cmdPipelinesCreate.Flags().StringVarP(&description, "description", "d", "", "long form description of pipeline")

	var gitURL string
	cmdPipelinesCreate.Flags().St

	cmdPipeline.AddCommand(cmdPipelinesCreate)
}
