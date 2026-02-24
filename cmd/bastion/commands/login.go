package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type loginResponse struct {
	Token string `json:"token"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Bastion server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		fmt.Print("Enter Admin Password: ")
		var password string
		fmt.Scanln(&password)

		payload, _ := json.Marshal(map[string]string{
			"password": password,
		})

		resp, err := http.Post(serverURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to connect to server: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("authentication failed: %s", resp.Status)
		}

		var loginResp loginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if err := saveToken(loginResp.Token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Successfully authenticated!")
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
	rootCmd.AddCommand(loginCmd)
}
