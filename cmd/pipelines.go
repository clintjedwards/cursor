package cmd

import (
	"github.com/spf13/cobra"
)

var cmdPipeline = &cobra.Command{
	Use:   "pipelines",
	Short: "Controls operations that can be performed on pipelines",
}

func init() {
	RootCmd.AddCommand(cmdPipeline)
}
