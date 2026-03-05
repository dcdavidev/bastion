package commands

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dcdavidev/bastion/packages/auth"
	"github.com/dcdavidev/bastion/packages/config"
	"github.com/dcdavidev/bastion/packages/crypto"
	"github.com/dcdavidev/bastion/packages/db"
	"github.com/dcdavidev/bastion/packages/models"
	"github.com/joho/godotenv"
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
	Short: "Authenticate and connect with a Bastion server",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")

		if serverURL == "" {
			var err error
			serverURL, err = pterm.DefaultInteractiveTextInput.Show("Enter Bastion Server URL (Leave empty for Local login)")
			if err != nil {
				return err
			}
		}

		if serverURL == "" {
			pterm.Info.Println("No URL provided. Attempting local authentication...")
			return handleLocalLogin()
		}

		// Remote Login Logic
		if !strings.HasPrefix(serverURL, "http") {
			serverURL = "http://" + serverURL
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Connecting to %s...", serverURL))
		status, err := checkServerStatus(serverURL)
		if err != nil {
			spinner.Fail("Failed to connect to server: " + err.Error())
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

		spinner, _ = pterm.DefaultSpinner.Start("Authenticating with remote server...")
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
		return saveLoginToConfig(serverURL, loginResp.Token)
	},
}

func handleLocalLogin() error {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("BASTION_JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("BASTION_JWT_SECRET not found in local environment or .env")
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

	spinner, _ := pterm.DefaultSpinner.Start("Authenticating locally...")

	var role string
	var userID string
	var username string

	if loginEmail != "" {
		database, err := db.NewConnection()
		if err != nil {
			spinner.Fail("Failed to connect to local database: " + err.Error())
			return err
		}
		defer database.Close()

		user, storedHashHex, saltHex, err := database.GetUserByEmail(context.Background(), loginEmail)
		if err != nil {
			spinner.Fail("User not found in local database")
			return err
		}

		salt, _ := hex.DecodeString(saltHex)
		storedHash, _ := hex.DecodeString(storedHashHex)
		computedHash := crypto.DeriveKey([]byte(password), salt)

		if subtle.ConstantTimeCompare(computedHash, storedHash) != 1 {
			spinner.Fail("Invalid local password")
			return fmt.Errorf("unauthorized")
		}

		role = user.Role
		userID = user.ID.String()
		username = user.Username
	} else {
		// Admin Fallback
		if !auth.VerifyAdmin(password) {
			spinner.Fail("Invalid admin credentials")
			return fmt.Errorf("unauthorized")
		}
		role = "ADMIN"
		userID = "00000000-0000-0000-0000-000000000000"
		username = "admin"
	}

	// Generate JWT locally
	uid, _ := models.ParseUUID(userID)
	token, err := auth.GenerateToken(uid, username, role == "ADMIN")
	if err != nil {
		spinner.Fail("Failed to generate local token")
		return err
	}

	spinner.Success("Local authentication successful!")

	port := os.Getenv("BASTION_PORT")
	if port == "" {
		port = "8287"
	}
	localURL := fmt.Sprintf("http://localhost:%s", port)

	return saveLoginToConfig(localURL, token)
}

func saveLoginToConfig(serverURL, token string) error {
	parsedURL, err := url.Parse(serverURL)
	profileName := "local"
	if err == nil && parsedURL.Host != "" && !strings.Contains(parsedURL.Host, "localhost") {
		profileName = parsedURL.Host
	}

	cfg, _ := config.LoadConfig()
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]config.Profile)
	}

	cfg.Profiles[profileName] = config.Profile{
		Name:     profileName,
		URL:      serverURL,
		Token:    token,
		IsActive: true,
	}
	cfg.ActiveProfile = profileName

	if err := cfg.Save(); err != nil {
		pterm.Error.Printf("Failed to save profile configuration: %v\n", err)
		return err
	}

	pterm.Success.Printf("Connection saved as profile '%s' in ~/.bastion/config.yaml\n", profileName)
	return nil
}

func checkServerStatus(serverURL string) (*StatusResponse, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(serverURL + "/api/v1/status")
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
