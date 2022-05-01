package cmd

import (
	"fmt"

	"github.com/Mario-F/cert-manager-selfservice/internal/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of cert-manager-selfservice",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
