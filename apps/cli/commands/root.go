package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/dcdavidev/bastion/packages/config"
	"github.com/dcdavidev/bastion/packages/version"
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
	Long:    `A secure, self-hosted fortress to manage multiple client secrets.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		if selectedProfile != "" {
			if p, ok := cfg.Profiles[selectedProfile]; ok {
				p_copy := p
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
	},
}

func runRootInteractive(cmd *cobra.Command, args []string) error {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Bastion", pterm.NewStyle(pterm.FgCyan)),
	).Render()
	pterm.Info.Printf("CLI Version: %s\n", version.Version)

	if activeProfile != nil && activeProfile.URL != "" {
		pterm.Success.Printf("Connected to: %s (Profile: %s)\n", activeProfile.URL, activeProfile.Name)
	} else {
		pterm.Warning.Println("No active connection. Please login to connect to a Bastion server.")
	}

	fmt.Println()

	options := []string{
		"Login - Connect to a Bastion server",
		"Init - Initialize local Bastion (database & admin)",
		"Create - Create resources",
		"Reset - Reset resources (credentials, etc.)",
		"Exit",
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("What do you want to do?")
	if err != nil {
		return err
	}

	switch {
	case strings.HasPrefix(selected, "Login"):
		return loginCmd.RunE(loginCmd, []string{})
	case strings.HasPrefix(selected, "Init"):
		return initCmd.RunE(initCmd, []string{})
	case strings.HasPrefix(selected, "Create"):
		return createInteractive()
	case strings.HasPrefix(selected, "Reset"):
		return resetInteractive()
	case selected == "Exit":
		return nil
	}
	return nil
}

// createCmd represents the create command group
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
}

func createInteractive() error {
	options := []string{
		"secretkey - Generate a new secure JWT secret",
		"masterkey - Initialize the vault with a new Master Key",
		"Back",
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("What do you want to create?")
	if err != nil {
		return err
	}

	if selected == "Back" {
		return runRootInteractive(rootCmd, []string{})
	}

	cmdStr := strings.Split(selected, " ")[0]
	for _, c := range createCmd.Commands() {
		if c.Use == cmdStr {
			return c.RunE(c, []string{})
		}
	}

	pterm.Error.Println("Command not implemented interactively yet")
	return nil
}

func resetInteractive() error {
	options := []string{
		"credentials - Reset user password",
		"Back",
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("What do you want to reset?")
	if err != nil {
		return err
	}

	if selected == "Back" {
		return runRootInteractive(rootCmd, []string{})
	}

	cmdStr := strings.Split(selected, " ")[0]
	for _, c := range resetCmd.Commands() {
		if c.Use == cmdStr {
			return c.RunE(c, []string{})
		}
	}

	pterm.Error.Println("Command not implemented interactively yet")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.RunE = runRootInteractive
	createCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createInteractive()
	}

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
}
