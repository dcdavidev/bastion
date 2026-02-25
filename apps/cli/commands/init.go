package commands

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/dcdavidev/bastion/packages/core/db"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Bastion (set DATABASE_URL, create superuser, apply migrations)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("BASTION INITIALIZATION")

		// 1. Set DATABASE_URL
		dbURL, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("postgres://bastion:bastion_password@localhost:5432/bastion?sslmode=disable").Show("Enter DATABASE_URL")

		config, err := LoadConfig()
		if err != nil {
			return err
		}
		config.DatabaseURL = dbURL
		if err := SaveConfig(config); err != nil {
			return err
		}

		// Set env for core/db to pick it up
		os.Setenv("DATABASE_URL", dbURL)

		// 2. Initialize DB & Run Migrations
		spinner, _ := pterm.DefaultSpinner.Start("Connecting to database and applying migrations...")
		database, err := db.NewConnection()
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer database.Close()

		if err := database.RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		spinner.Success("Database initialized and migrations applied!")

		// 3. Create Superuser
		pterm.DefaultHeader.WithFullWidth().Println("CREATE SUPERUSER")
		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Admin Password")

		if len(password) < 8 {
			return fmt.Errorf("password too short (min 8 chars)")
		}

		spinner, _ = pterm.DefaultSpinner.Start("Generating cryptographic material...")

		// Generate Salt
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return err
		}

		// Derive KEK
		kek := crypto.DeriveKey([]byte(password), salt)

		// Generate Master Key
		masterKey, err := crypto.GenerateRandomKey()
		if err != nil {
			return err
		}

		// Wrap Master Key with KEK
		wrappedMK, err := crypto.WrapKey(kek, masterKey)
		if err != nil {
			return err
		}
		spinner.Success("Cryptographic material generated!")

		// Save credentials to config
		config.AdminPasswordHash = hex.EncodeToString(kek)
		config.AdminPasswordSalt = hex.EncodeToString(salt)
		if err := SaveConfig(config); err != nil {
			return err
		}

		// Save to database
		spinner, _ = pterm.DefaultSpinner.Start("Saving configuration to database...")
		_, err = database.Pool.Exec(cmd.Context(), 
			"INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ($1, $2)",
			hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))
		if err != nil {
			return fmt.Errorf("failed to save vault config to database: %w", err)
		}
		spinner.Success("Superuser initialized successfully!")

		pterm.Success.Println("Bastion initialization complete! Configuration saved to ~/.config/bastion.yml")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
