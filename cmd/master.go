package cmd

import (
	"github.com/spf13/cobra"
)

var cmdMaster = &cobra.Command{
	Use:   "master",
	Short: "Starts a Cursor master server instance. Runs until there is an interrupt",
	Run:   runMasterCmd,
}

func runMasterCmd(cme *cobra.Command, args []string) {

}

func init() {
	RootCmd.AddCommand(cmdMaster)
}
