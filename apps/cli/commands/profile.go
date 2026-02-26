package commands

import (
	"fmt"

	"github.com/dcdavidev/bastion/packages/core/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Bastion profiles (environments)",
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		tableData := pterm.TableData{
			{"Active", "Name", "URL"},
		}

		for name, p := range cfg.Profiles {
			active := ""
			if name == cfg.ActiveProfile {
				active = pterm.FgGreen.Sprint("●")
			}
			tableData = append(tableData, []string{
				active,
				name,
				p.URL,
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		return nil
	},
}

var profileAddCmd = &cobra.Command{
	Use:   "add [NAME] [URL]",
	Short: "Add a new profile",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		cfg.Profiles[name] = config.Profile{
			Name: name,
			URL:  url,
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		pterm.Success.Printf("Profile '%s' added successfully!\n", name)
		return nil
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use [NAME]",
	Short: "Set the active profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		if _, ok := cfg.Profiles[name]; !ok {
			return fmt.Errorf("profile '%s' not found", name)
		}

		cfg.ActiveProfile = name
		if err := cfg.Save(); err != nil {
			return err
		}

		pterm.Success.Printf("Now using profile '%s'\n", name)
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileUseCmd)
	rootCmd.AddCommand(profileCmd)
}
