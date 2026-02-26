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
	Short:   "Generate credentials and master key for manual initial setup",
	Long:    "This command generates the necessary cryptographic material for manual setup. Use 'bastion init' for a guided setup.",
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

		// 2. Derive KEK
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

		pterm.DefaultHeader.WithFullWidth().Println("MANUAL INITIALIZATION REQUIRED")
		
		pterm.Info.Println("1. Run the following SQL command in your PostgreSQL database:")
		pterm.DefaultBox.Println(fmt.Sprintf("INSERT INTO vault_config (wrapped_master_key, master_key_salt) VALUES ('%s', '%s');", 
			hex.EncodeToString(wrappedMK), hex.EncodeToString(salt)))

		pterm.Info.Println("2. Create your admin user in the 'users' table using these values:")
		pterm.DefaultBox.Println(fmt.Sprintf("Password Hash: %s\nSalt: %s", hex.EncodeToString(kek), hex.EncodeToString(salt)))

		pterm.Info.Println("3. Ensure your server environment variables are set correctly.")

		return nil
	},
}

func init() {
	createCmd.AddCommand(createSuperuserCmd)
}
