package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dcdavidev/bastion/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type StatusResponse struct {
	ConnectedToDB   bool     `json:"connected_to_db"`
	MissingEnvVars  []string `json:"missing_env_vars"`
	JwtSecretStatus string   `json:"jwt_secret_status"`
	Migrations      struct {
		CurrentVersion uint `json:"current_version"`
		HasPending     bool `json:"has_pending"`
		IsDirty        bool `json:"is_dirty"`
	} `json:"migrations"`
	HasAdmin bool   `json:"has_admin"`
	Version  string `json:"version"`
}

type loginResponse struct {
	Token string `json:"token"`
}

var loginEmail string
var loginPassword string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and connect with the Bastion server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")

		if serverURL == "" {
			var err error
			serverURL, err = pterm.DefaultInteractiveTextInput.WithDefaultText("http://localhost:8287").Show("Enter Bastion Server URL")
			if err != nil {
				return err
			}
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Connecting to %s...", serverURL))
		status, err := checkServerStatus(serverURL)
		if err != nil {
			spinner.Fail("Failed to connect to server")
			return err
		}
		spinner.Success("Connected to Bastion server!")

		pterm.DefaultSection.Println("Server Status")
		pterm.Info.Printf("Version: %s\n", status.Version)
		if status.ConnectedToDB {
			pterm.Success.Println("Database: Connected")
		} else {
			pterm.Error.Println("Database: Disconnected")
		}

		if status.JwtSecretStatus == "missing" || status.JwtSecretStatus == "weak" {
			pterm.Warning.Printf("JWT Secret Status: %s\n", status.JwtSecretStatus)
			pterm.Info.Println("We strongly recommend generating a new secure secret. Run 'bastion create secretkey'.")
		} else {
			pterm.Success.Println("JWT Secret Status: strong")
		}

		if loginEmail == "" {
			var err error
			loginEmail, err = pterm.DefaultInteractiveTextInput.Show("Enter Email (leave empty for Admin fallback)")
			if err != nil {
				return err
			}
		}

		password := loginPassword
		if password == "" {
			var err error
			password, err = pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Password")
			if err != nil {
				return err
			}
		}

		spinner, _ = pterm.DefaultSpinner.Start("Authenticating...")
		payload, _ := json.Marshal(map[string]string{
			"email":    loginEmail,
			"password": password,
		})

		resp, err := http.Post(serverURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			spinner.Fail("Failed to connect to server for authentication")
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

		spinner.Success("Successfully authenticated!")

		// Save to config
		parsedURL, err := url.Parse(serverURL)
		profileName := "default"
		if err == nil && parsedURL.Host != "" {
			profileName = parsedURL.Host
		}

		cfg, _ := config.LoadConfig()
		if cfg.Profiles == nil {
			cfg.Profiles = make(map[string]config.Profile)
		}

		cfg.Profiles[profileName] = config.Profile{
			Name:     profileName,
			URL:      serverURL,
			Token:    loginResp.Token,
			IsActive: true,
		}
		cfg.ActiveProfile = profileName

		if err := cfg.Save(); err != nil {
			pterm.Error.Printf("Failed to save profile configuration: %v\n", err)
		} else {
			pterm.Success.Printf("Connection saved as profile '%s' in ~/.bastion/config.yaml\n", profileName)
		}

		return nil
	},
}

func checkServerStatus(serverURL string) (*StatusResponse, error) {
	resp, err := http.Get(serverURL + "/api/v1/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

func init() {
	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "Email for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	rootCmd.AddCommand(loginCmd)
}
