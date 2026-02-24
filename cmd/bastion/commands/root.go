package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bastion",
	Short: "Bastion is a secure E2EE secrets vault CLI",
	Long: `A secure, self-hosted fortress to manage multiple client secrets 
with blind-backend architecture. For more info: https://github.com/dcdavidev/bastion`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be defined here
	rootCmd.PersistentFlags().StringP("url", "u", "http://localhost:8080", "Bastion server URL")
}
