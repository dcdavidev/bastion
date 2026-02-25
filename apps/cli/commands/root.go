package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bastion",
	Short: "Bastion is a secure E2EE secrets vault CLI",
	Long: `A secure, self-hosted fortress to manage multiple client secrets 
with blind-backend architecture. For more info: https://github.com/dcdavidev/bastion`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle("Bastion", pterm.NewStyle(pterm.FgCyan)),
		).Render()
		pterm.Info.Println("The Secure E2EE Secrets Vault")
	},
}

// createCmd represents the create command group
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources (superuser, collaborator, project, client, jwtsecret, masterkey)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	config, _ := LoadConfig()
	serverURL := "http://localhost:8081"
	if config != nil && config.ServerURL != "" {
		serverURL = config.ServerURL
	}

	rootCmd.PersistentFlags().StringP("url", "u", serverURL, "Bastion server URL")

	// Add create command group
	rootCmd.AddCommand(createCmd)
}
