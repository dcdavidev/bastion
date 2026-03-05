package commands

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/dcdavidev/bastion/packages/crypto"
	"github.com/dcdavidev/bastion/packages/db"
	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Bastion locally (database, migrations, and admin user)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("BASTION LOCAL INITIALIZATION")

		// 1. Load .env if exists
		_ = godotenv.Load()

		dbURL := os.Getenv("BASTION_DATABASE_URL")
		if dbURL == "" {
			dbURL = os.Getenv("DATABASE_URL")
		}

		if dbURL == "" {
			pterm.DefaultSection.Println("Database Configuration")
			dbHost, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("localhost").Show("Database Host")
			dbPort, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("5432").Show("Database Port")
			dbUser, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("bastion").Show("Database User")
			dbPass, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Database Password")
			dbName, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("bastion").Show("Database Name")

			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

			// Save to .env
			saveToEnv("BASTION_DATABASE_URL", dbURL)
			os.Setenv("BASTION_DATABASE_URL", dbURL)
			pterm.Success.Println("BASTION_DATABASE_URL configured and saved to .env")
		} else {
			pterm.Info.Printf("Using existing database connection: %s\n", dbURL)
		}

		// 2. JWT Secret
		jwtSecret := os.Getenv("BASTION_JWT_SECRET")
		if jwtSecret == "" || jwtSecret == "bastion_very_secret_key_change_me" {
			pterm.Info.Println("Generating a secure BASTION_JWT_SECRET...")
			secret := make([]byte, 64)
			if err := crypto.GenerateRandomKeyInto(secret); err != nil {
				return err
			}
			jwtSecret = hex.EncodeToString(secret)

			// Save to .env
			saveToEnv("BASTION_JWT_SECRET", jwtSecret)
			os.Setenv("BASTION_JWT_SECRET", jwtSecret)
			pterm.Success.Println("BASTION_JWT_SECRET generated and saved to .env")
		}

		// 3. Database Connection & Migrations
		spinner, _ := pterm.DefaultSpinner.Start("Connecting to database and running migrations...")
		database, err := db.NewConnection()
		if err != nil {
			spinner.Fail("Connection failed: " + err.Error())
			return err
		}
		defer database.Close()

		if err := database.RunMigrations(); err != nil {
			spinner.Fail("Migration failed: " + err.Error())
			return err
		}
		spinner.Success("Database is up to date!")

		// 4. Check for Admin
		hasAdmin, err := database.HasAdmin(context.Background())
		if err != nil {
			return err
		}

		if !hasAdmin {
			pterm.DefaultSection.Println("Create Admin User")
			email, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("admin@bastion.local").Show("Admin Email")
			username, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("admin").Show("Admin Username")
			password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Admin Password")

			spinner, _ = pterm.DefaultSpinner.Start("Creating admin and initializing vault...")

			salt, _ := crypto.GenerateSalt()
			hash := crypto.DeriveKey([]byte(password), salt)

			saltHex := hex.EncodeToString(salt)
			hashHex := hex.EncodeToString(hash)

			user, err := database.CreateUser(context.Background(), username, email, hashHex, saltHex, "ADMIN")
			if err != nil {
				spinner.Fail("Failed to create admin: " + err.Error())
				return err
			}

			// Initialize Vault
			masterKey, _ := crypto.GenerateRandomKey()
			kek := crypto.DeriveKey([]byte(password), salt)
			wrappedMK, _ := crypto.WrapKey(kek, masterKey)

			err = database.InitializeVault(context.Background(), hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))
			if err != nil {
				spinner.Fail("Failed to initialize vault: " + err.Error())
				return err
			}

			// Grant admin access to their own master key?
			// In our current architecture, the vault initialization already sets up the master key.
			// The admin will be able to decrypt it because they know the password and the salt is in the vault_config.

			spinner.Success(fmt.Sprintf("Admin user '%s' created successfully!", user.Username))
		} else {
			pterm.Info.Println("Admin user already exists.")
		}

		pterm.Success.Println("Bastion initialization complete!")
		return nil
	},
}

func saveToEnv(key, value string) {
	envPath := ".env"
	content, err := os.ReadFile(envPath)
	if err != nil {
		_ = os.WriteFile(envPath, []byte(fmt.Sprintf("%s=%s\n", key, value)), 0644)
		return
	}

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	_ = os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644)
}

func init() {
	rootCmd.AddCommand(initCmd)
}
