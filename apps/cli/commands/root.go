package commands

import (
	"fmt"
	"os"

	"github.com/dcdavidev/bastion/packages/core/config"
	"github.com/dcdavidev/bastion/packages/core/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	selectedProfile string
	activeProfile   *config.Profile
)

var rootCmd = &cobra.Command{
	Use:     "bastion",
	Version: version.Version,
	Short:   "Bastion is a secure E2EE secrets vault CLI",
	Long: `A secure, self-hosted fortress to manage multiple client secrets 
with blind-backend architecture. For more info: https://github.com/dcdavidev/bastion`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		cfg, _ := config.LoadConfig()
		if selectedProfile != "" {
			if p, ok := cfg.Profiles[selectedProfile]; ok {
				p_copy := p // create copy to avoid pointer issues
				activeProfile = &p_copy
			} else {
				pterm.Warning.Printf("Profile '%s' not found, using active profile '%s'\n", selectedProfile, cfg.ActiveProfile)
				activeProfile = cfg.GetActiveProfile()
			}
		} else {
			activeProfile = cfg.GetActiveProfile()
		}

		if os.Getenv("BASTION_TEST") == "true" {
			return
		}
		CheckForUpdates()
		pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle("Bastion", pterm.NewStyle(pterm.FgCyan)),
		).Render()
		pterm.Info.Println("The Secure E2EE Secrets Vault")
		if activeProfile != nil {
			pterm.Info.Printf("Active Profile: %s (%s)\n", pterm.Bold.Sprint(activeProfile.Name), activeProfile.URL)
		}
		fmt.Println()
	},
}

// createCmd represents the create command group
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources (superuser, collaborator, project, client, jwtsecret, masterkey)",
}

// listCmd represents the list command group
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources (clients, projects)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	cfg, _ := config.LoadConfig()
	defaultURL := ""
	if cfg != nil && cfg.GetActiveProfile() != nil {
		defaultURL = cfg.GetActiveProfile().URL
	}

	rootCmd.PersistentFlags().StringP("url", "u", defaultURL, "Bastion server URL")
	rootCmd.PersistentFlags().StringVarP(&selectedProfile, "profile", "P", "", "Profile to use")

	// Add command groups
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
}
