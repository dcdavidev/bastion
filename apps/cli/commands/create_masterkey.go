package commands

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/dcdavidev/bastion/packages/crypto"
	"github.com/dcdavidev/bastion/packages/db"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createMasterKeyCmd = &cobra.Command{
	Use:   "masterkey",
	Short: "Initialize the vault with a new Master Key",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("MASTER KEY INITIALIZATION")

		spinner, _ := pterm.DefaultSpinner.Start("Connecting to database...")
		database, err := db.NewConnection()
		if err != nil {
			spinner.Fail("Connection failed: " + err.Error())
			return err
		}
		defer database.Close()
		spinner.Success("Connected to database!")

		// 1. Check if vault is already initialized
		vault, err := database.GetVaultConfig(context.Background())
		if err == nil && vault.WrappedMasterKey != "" {
			pterm.Warning.Println("Vault is already initialized with a Master Key.")
			confirm, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show("Do you want to OVERWRITE the existing Master Key? (THIS WILL RENDER ALL EXISTING SECRETS UNREADABLE!)")
			if !confirm {
				pterm.Info.Println("Operation cancelled.")
				return nil
			}
		}

		// 2. Identify an admin to wrap the key
		pterm.Info.Println("A Master Key must be wrapped with an Admin password.")
		adminUser, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("admin").Show("Enter Admin Username to associate with this Master Key")

		user, _, _, err := database.GetUserByUsername(context.Background(), adminUser)
		if err != nil {
			return fmt.Errorf("user '%s' not found: %w", adminUser, err)
		}
		if user.Role != "ADMIN" {
			return fmt.Errorf("user '%s' is not an ADMIN", adminUser)
		}

		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show(fmt.Sprintf("Enter password for '%s' to wrap the Master Key", adminUser))

		// 3. Generate and wrap
		spinner, _ = pterm.DefaultSpinner.Start("Generating and wrapping new Master Key...")

		masterKey, _ := crypto.GenerateRandomKey()
		salt, _ := crypto.GenerateSalt()
		kek := crypto.DeriveKey([]byte(password), salt)
		wrappedMK, err := crypto.WrapKey(kek, masterKey)
		if err != nil {
			spinner.Fail("Failed to wrap key: " + err.Error())
			return err
		}

		// 4. Save
		if vault != nil && vault.WrappedMasterKey != "" {
			err = database.UpdateVaultConfig(context.Background(), hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))
		} else {
			err = database.InitializeVault(context.Background(), hex.EncodeToString(wrappedMK), hex.EncodeToString(salt))
		}

		if err != nil {
			spinner.Fail("Failed to save vault configuration: " + err.Error())
			return err
		}

		spinner.Success("Master Key successfully created and vault initialized!")
		pterm.Info.Println("Note: This key is now encrypted with the admin's password and stored in the database.")

		return nil
	},
}

func init() {
	createCmd.AddCommand(createMasterKeyCmd)
}
