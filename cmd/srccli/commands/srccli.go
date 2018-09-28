package commands

import (
	"os"
	"github.com/spf13/cobra"
)

var SrcCliCmd = &cobra.Command{
	Use:   "srccli",
	Short: "Srccli is a command line client fot silkroad chain core",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func AddCommands(){
	TxFeedCmd()
}

func Execute() {
	AddCommands()
	//AddTemplateFunc()

	if _, err := SrcCliCmd.ExecuteC(); err != nil {
		os.Exit(1)
	}
}