package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "development"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of cert-manager-selfservice",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
