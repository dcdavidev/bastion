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

var resetUser string

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset resources",
}

var resetCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Reset user credentials (password)",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("CREDENTIALS RESET")

		if resetUser == "" {
			var err error
			resetUser, err = pterm.DefaultInteractiveTextInput.Show("Enter Username to reset")
			if err != nil {
				return err
			}
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Searching for user '%s'...", resetUser))
		database, err := db.NewConnection()
		if err != nil {
			spinner.Fail("Connection failed: " + err.Error())
			return err
		}
		defer database.Close()

		user, hashHex, saltHexDB, err := database.GetUserByUsername(context.Background(), resetUser)
		if err != nil {
			spinner.Fail("User not found: " + err.Error())
			return err
		}
		spinner.Success(fmt.Sprintf("User '%s' found!", user.Username))

		oldPassword, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter OLD Password (leave empty to skip re-wrapping)")
		newPassword, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter NEW Password")
		confirmPassword, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Confirm NEW Password")

		if newPassword != confirmPassword {
			pterm.Error.Println("Passwords do not match!")
			return fmt.Errorf("passwords do not match")
		}

		// Re-wrapping logic for ADMIN
		reWrapped := false
		if user.Role == "ADMIN" && oldPassword != "" {
			spinner, _ = pterm.DefaultSpinner.Start("Authenticating and re-wrapping Master Key...")

			// Verify old password for LOGIN first
			saltDB, _ := hex.DecodeString(saltHexDB)
			loginKek := crypto.DeriveKey([]byte(oldPassword), saltDB)
			loginHash := hex.EncodeToString(loginKek)

			if loginHash != hashHex {
				spinner.Fail("Old password verification failed. Cannot re-wrap Master Key.")
			} else {
				// Re-wrap Master Key using the VAULT salt
				vault, err := database.GetVaultConfig(context.Background())
				if err != nil {
					spinner.Fail("Could not fetch vault configuration: " + err.Error())
				} else {
					vaultSalt, err := hex.DecodeString(vault.MasterKeySalt)
					if err != nil {
						spinner.Fail("Invalid salt format in vault: " + err.Error())
					} else {
						// Derive the KEK used for the vault
						oldVaultKek := crypto.DeriveKey([]byte(oldPassword), vaultSalt)

						wrappedMK, err := hex.DecodeString(vault.WrappedMasterKey)
						if err != nil {
							spinner.Fail("Invalid wrapped key format in vault: " + err.Error())
						} else {
							masterKey, err := crypto.UnwrapKey(oldVaultKek, wrappedMK)
							if err == nil {
								newVaultSalt, _ := crypto.GenerateSalt()
								newVaultKek := crypto.DeriveKey([]byte(newPassword), newVaultSalt)
								newWrappedMK, _ := crypto.WrapKey(newVaultKek, masterKey)

								err = database.UpdateVaultConfig(context.Background(), hex.EncodeToString(newWrappedMK), hex.EncodeToString(newVaultSalt))
								if err == nil {
									reWrapped = true
									spinner.Success("Master Key re-wrapped and vault configuration updated!")
								} else {
									spinner.Fail("Failed to update vault configuration: " + err.Error())
								}
							} else {
								spinner.Fail("Failed to decrypt Master Key with old password. Are you sure it's the right one? Error: " + err.Error())
							}
						}
					}
				}
				if !reWrapped && spinner.IsActive {
					spinner.Fail("Failed to re-wrap Master Key. Vault might be inaccessible with new password.")
				}
			}
		}

		spinner, _ = pterm.DefaultSpinner.Start("Updating credentials...")
		
		var finalSaltHex, finalHashHex string
		if reWrapped {
			// If we re-wrapped, the salt and hash were already used for vault, so we use them
			// Actually, let's keep it simple: generate new salt/hash for user login too.
			// Re-wrapping logic above already updated vault_config.
		}

		salt, _ := crypto.GenerateSalt()
		hash := crypto.DeriveKey([]byte(newPassword), salt)
		finalSaltHex = hex.EncodeToString(salt)
		finalHashHex = hex.EncodeToString(hash)

		err = database.UpdateUserPassword(context.Background(), user.ID, finalHashHex, finalSaltHex)
		if err != nil {
			spinner.Fail("Failed to update password: " + err.Error())
			return err
		}

		if user.Role == "ADMIN" && !reWrapped {
			pterm.Warning.Println("User is an ADMIN but Master Key was NOT re-wrapped.")
			pterm.Warning.Println("The vault will be inaccessible unless the old password is used to recover it.")
		}

		spinner.Success("Password updated successfully!")
		return nil
	},
}

func init() {
	resetCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return resetInteractive()
	}
	resetCredentialsCmd.Flags().StringVarP(&resetUser, "user", "u", "", "Username to reset")
	resetCmd.AddCommand(resetCredentialsCmd)
	rootCmd.AddCommand(resetCmd)
}
