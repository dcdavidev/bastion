package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createSuperuserCmd = &cobra.Command{
	Use:     "superuser",
	Aliases: []string{"su"},
	Short:   "Generate credentials and master key for the initial admin setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Admin Password")

		if len(password) < 8 {
			return fmt.Errorf("password too short (min 8 chars)")
		}

		spinner, _ := pterm.DefaultSpinner.Start("Generating cryptographic material...")

		// 1. Generate Salt for Admin Password
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return err
		}

		// 2. Derive KEK (and also the password hash for config)
		kek := crypto.DeriveKey([]byte(password), salt)

		// 3. Generate Global Master Key
		masterKey, err := crypto.GenerateRandomKey()
		if err != nil {
			return err
		}

		// 4. Wrap Master Key with KEK
		wrappedMK, err := crypto.WrapKey(kek, masterKey)
		if err != nil {
			return err
		}

		spinner.Success("Cryptographic material generated!")

		// Save to Config
		config, err := LoadConfig()
		if err != nil {
			return err
		}

		config.AdminPasswordHash = hex.EncodeToString(kek)
		config.AdminPasswordSalt = hex.EncodeToString(salt)
		
		if err := SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		pterm.Success.Println("Saved credentials to ~/.config/bastion.yml")

		pterm.DefaultHeader.WithFullWidth().Println("DATABASE INITIALIZATION REQUIRED")
		pterm.Println("Run the following SQL command in your PostgreSQL database:")
		
		pterm.DefaultBox.Println(fmt.Sprintf("INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ('%s', '%s');", 
			hex.EncodeToString(wrappedMK), hex.EncodeToString(salt)))

		return nil
	},
}

func init() {
	createCmd.AddCommand(createSuperuserCmd)
}
