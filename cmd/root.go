package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the base of the command line interface all other commands build from here
var RootCmd = &cobra.Command{
	Use:   "cursor",
	Short: "cursor",
}

// Execute runs the command line tool
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
