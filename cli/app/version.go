package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	sparkCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "spark ai cli versions.",
	Long:  `Print the version number of spark`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Spark AI CLI v0.1 -- HEAD")
	},
}
