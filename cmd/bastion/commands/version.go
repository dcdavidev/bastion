package commands

import (
	"github.com/dcdavidev/bastion/internal/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Bastion",
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Printf("Bastion CLI version: %s
", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
