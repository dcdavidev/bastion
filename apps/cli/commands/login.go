package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type loginResponse struct {
	Token string `json:"token"`
}

var loginEmail string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Bastion server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		if loginEmail == "" {
			loginEmail, _ = pterm.DefaultInteractiveTextInput.Show("Enter Email (leave empty for Admin fallback)")
		}
		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Password")

		spinner, _ := pterm.DefaultSpinner.Start("Authenticating...")

		payload, _ := json.Marshal(map[string]string{
			"email":    loginEmail,
			"password": password,
		})

		resp, err := http.Post(serverURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			spinner.Fail("Failed to connect to server")
			return fmt.Errorf("failed to connect to server: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			spinner.Fail("Authentication failed")
			return fmt.Errorf("authentication failed: %s", resp.Status)
		}

		var loginResp loginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			spinner.Fail("Failed to decode response")
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if err := saveToken(loginResp.Token); err != nil {
			spinner.Fail("Failed to save token")
			return fmt.Errorf("failed to save token: %w", err)
		}

		spinner.Success("Successfully authenticated!")
		return nil
	},
}

func saveToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".bastion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "token")
	return os.WriteFile(configPath, []byte(token), 0600)
}

func init() {
	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Email for login")
	rootCmd.AddCommand(loginCmd)
}
