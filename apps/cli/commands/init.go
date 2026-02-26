package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/dcdavidev/bastion/packages/core/db"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type BastionStatus struct {
	ConnectedToDB  bool     `json:"connected_to_db"`
	MissingEnvVars []string `json:"missing_env_vars"`
	Migrations     struct {
		CurrentVersion uint `json:"current_version"`
		HasPending     bool `json:"has_pending"`
		IsDirty        bool `json:"is_dirty"`
	} `json:"migrations"`
	HasAdmin bool   `json:"has_admin"`
	Version  string `json:"version"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup wizard for Bastion",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("BASTION SETUP WIZARD")

		serverURL, _ := cmd.Flags().GetString("url")
		if serverURL == "" {
			serverURL = "http://localhost:8287"
		}

		// 1. Check Status
		spinner, _ := pterm.DefaultSpinner.Start("Checking Bastion Server status at " + serverURL + "...")
		status, err := checkBastionStatus(serverURL)
		if err != nil {
			spinner.Warning("Could not connect to Bastion Server. Proceeding with local configuration.")
			status = &BastionStatus{}
		} else {
			spinner.Success("Server connected!")
		}

		envVars := make(map[string]string)

		// 1b. Server Connectivity Info
		if os.Getenv("BASTION_HOST") == "" {
			envVars["BASTION_HOST"] = serverURL
		}
		if os.Getenv("BASTION_PORT") == "" {
			// Extract port from URL if possible, or default to 8287
			port := "8287"
			if strings.Contains(serverURL, ":") {
				parts := strings.Split(serverURL, ":")
				last := parts[len(parts)-1]
				if !strings.Contains(last, "/") {
					port = last
				}
			}
			envVars["BASTION_PORT"] = port
		}

		// 2. Database Configuration
		dbURL := os.Getenv("BASTION_DATABASE_URL")
		if dbURL == "" {
			dbURL = os.Getenv("DATABASE_URL")
		}
		if dbURL == "" || !status.ConnectedToDB {
			pterm.DefaultSection.Println("Database Configuration")
			dbHost, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("localhost").Show("Database Host")
			dbPort, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("5432").Show("Database Port")
			dbUser, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("bastion").Show("Database User")
			dbPass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Database Password")
			dbName, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("bastion").Show("Database Name")

			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
			envVars["BASTION_DATABASE_URL"] = dbURL
			os.Setenv("BASTION_DATABASE_URL", dbURL)
		} else {
			pterm.Success.Println("Database already connected!")
		}

		// 3. JWT Secret
		jwtSecret := os.Getenv("BASTION_JWT_SECRET")
		if jwtSecret == "" {
			secret := make([]byte, 32)
			crypto.GenerateRandomKeyInto(secret)
			jwtSecret = hex.EncodeToString(secret)
			envVars["BASTION_JWT_SECRET"] = jwtSecret
		}

		// 4. Migrations
		if status.Migrations.HasPending || status.Migrations.IsDirty || !status.ConnectedToDB {
			spinner, _ := pterm.DefaultSpinner.Start("Running migrations...")
			database, err := db.NewConnection()
			if err != nil {
				return fmt.Errorf("failed to connect to database for migrations: %w", err)
			}
			if err := database.RunMigrations(); err != nil {
				spinner.Fail("Migration failed: " + err.Error())
				return err
			}
			spinner.Success("Migrations applied!")
			database.Close()
		}

		// 5. Admin Configuration
		if !status.HasAdmin {
			pterm.DefaultSection.Println("Admin Configuration")
			adminUsername, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("admin").Show("Admin Username")
			adminEmail, _ := pterm.DefaultInteractiveTextInput.Show("Admin Email")
			adminPassword, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Admin Password")

			spinner, _ := pterm.DefaultSpinner.Start("Creating admin user...")
			database, _ := db.NewConnection()
			
			salt, _ := crypto.GenerateSalt()
			hash := crypto.DeriveKey([]byte(adminPassword), salt)
			
			saltHex := hex.EncodeToString(salt)
			hashHex := hex.EncodeToString(hash)

			_, err = database.CreateUser(cmd.Context(), adminUsername, adminEmail, hashHex, saltHex, "ADMIN")
			if err != nil {
				spinner.Fail("Failed to create admin: " + err.Error())
				return err
			}

			// Add to environment variables for fallback/sync
			envVars["BASTION_ADMIN_PASSWORD_HASH"] = hashHex
			envVars["BASTION_ADMIN_PASSWORD_SALT"] = saltHex

			// Vault Initialization
			masterKey, _ := crypto.GenerateRandomKey()
			kek := crypto.DeriveKey([]byte(adminPassword), salt)
			wrappedMK, _ := crypto.WrapKey(kek, masterKey)
			database.InitializeVault(cmd.Context(), hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))

			spinner.Success("Admin created and vault initialized!")
			database.Close()
		}

		// 6. Environment Storage
		if len(envVars) > 0 {
			pterm.DefaultSection.Println("Environment Configuration")
			options := []string{".env file", ".bashrc", ".zshrc", "Password Store (pass)", "Display only"}
			selected, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Where do you want to save the environment variables?")

			switch selected {
			case ".env file":
				saveToEnvFile(envVars)
			case ".bashrc":
				saveToShellProfile(".bashrc", envVars)
			case ".zshrc":
				saveToShellProfile(".zshrc", envVars)
			case "Password Store (pass)":
				saveToPass(envVars)
			default:
				displayEnv(envVars)
			}
		}

		pterm.Success.Println("Bastion setup complete!")
		return nil
	},
}

func checkBastionStatus(url string) (*BastionStatus, error) {
	resp, err := http.Get(url + "/api/v1/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status BastionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

func saveToEnvFile(vars map[string]string) {
	f, _ := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	fmt.Fprintln(f, "\n## BASTION")
	for k, v := range vars {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}
	fmt.Fprintln(f, "## ENDBASTION")
	pterm.Success.Println("Variables saved to .env")
}

func saveToShellProfile(filename string, vars map[string]string) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, filename)
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	fmt.Fprintln(f, "\n## BASTION")
	for k, v := range vars {
		fmt.Fprintf(f, "export %s=\"%s\"\n", k, v)
	}
	fmt.Fprintln(f, "## ENDBASTION")
	pterm.Success.Printf("Variables saved to %s. Please run 'source ~/%s' to apply changes.\n", filename, filename)
}

func saveToPass(vars map[string]string) {
	storePath := os.Getenv("BASTION_STORE_DIR")
	if storePath == "" {
		storePath = filepath.Join(os.Getenv("HOME"), ".config", "bastion", "store")
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		initGit, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show("Initialize a new git-backed password store at " + storePath + "?")
		if initGit {
			gpgID, _ := pterm.DefaultInteractiveTextInput.Show("Enter GPG ID to initialize the store")
			if gpgID != "" {
				exec.Command("pass", "init", gpgID).Run()
				os.Setenv("PASSWORD_STORE_DIR", storePath)
				exec.Command("pass", "git", "init").Run()
				pterm.Info.Println("Password store initialized with Git. Remember to set up a private remote!")
			}
		}
	}

	for k, v := range vars {
		cmd := exec.Command("pass", "insert", "-m", "bastion/"+k)
		cmd.Stdin = strings.NewReader(v + "\n")
		cmd.Run()
	}
	pterm.Success.Println("Variables saved to Password Store (bastion/*)")
}

func displayEnv(vars map[string]string) {
	pterm.Info.Println("Please export the following variables manually:")
	for k, v := range vars {
		pterm.DefaultBox.Printf("export %s=%s\n", k, v)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
