package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createSuperuserCmd = &cobra.Command{
	Use:   "create-superuser",
	Short: "Generate credentials and master key for the initial admin setup",
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

		// 2. Derive KEK (and also the password hash for .env)
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

		pterm.DefaultHeader.WithFullWidth().Println("STEP 1: Update your .env file")
		pterm.Info.Printf("ADMIN_PASSWORD_HASH=%s\n", hex.EncodeToString(kek))
		pterm.Info.Printf("ADMIN_PASSWORD_SALT=%s\n", hex.EncodeToString(salt))

		pterm.DefaultHeader.WithFullWidth().Println("STEP 2: Initialize your Database")
		pterm.Println("Run the following SQL command in your PostgreSQL database:")
		
		pterm.DefaultBox.Println(fmt.Sprintf("INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ('%s', '%s');", 
			hex.EncodeToString(wrappedMK), hex.EncodeToString(salt)))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createSuperuserCmd)
}
