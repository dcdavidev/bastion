package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createJWTSecretCmd = &cobra.Command{
	Use:   "jwtsecret",
	Short: "Generate a new random JWT secret and save it to config",
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate random secret: %w", err)
		}

		encodedSecret := hex.EncodeToString(secret)

		config, err := LoadConfig()
		if err != nil {
			return err
		}

		config.JWTSecret = encodedSecret
		if err := SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		pterm.Success.Println("Generated and saved new JWT secret to ~/.config/bastion.yml")
		pterm.Info.Printf("JWT_SECRET=%s\n", encodedSecret)

		return nil
	},
}

func init() {
	createCmd.AddCommand(createJWTSecretCmd)
}
