package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"

	"github.com/gkwa/halfsl500/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of halfsl500",
	Long:  `All software has versions. This is halfsl500's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
