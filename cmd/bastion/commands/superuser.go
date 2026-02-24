package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/dcdavidev/bastion/internal/crypto"
	"github.com/spf13/cobra"
)

var createSuperuserCmd = &cobra.Command{
	Use:   "create-superuser",
	Short: "Generate credentials and master key for the initial admin setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter Admin Password: ")
		var password string
		fmt.Scanln(&password)

		if len(password) < 8 {
			return fmt.Errorf("password too short (min 8 chars)")
		}

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

		fmt.Println("
--- STEP 1: Update your .env file ---")
		fmt.Printf("ADMIN_PASSWORD_HASH=%s
", hex.EncodeToString(kek))
		fmt.Printf("ADMIN_PASSWORD_SALT=%s
", hex.EncodeToString(salt))

		fmt.Println("
--- STEP 2: Initialize your Database ---")
		fmt.Println("Run the following SQL command in your PostgreSQL database:")
		fmt.Printf("INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ('%s', '%s');
", 
			hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createSuperuserCmd)
}
