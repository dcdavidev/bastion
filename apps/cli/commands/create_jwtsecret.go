package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var printOnly bool

var createJWTSecretCmd = &cobra.Command{
	Use:   "jwtsecret",
	Short: "Generate a new random JWT secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate random secret: %w", err)
		}

		encodedSecret := hex.EncodeToString(secret)

		if printOnly {
			fmt.Print(encodedSecret)
			return nil
		}

		pterm.Success.Println("Generated new random JWT secret!")
		pterm.Info.Println("Please set it as an environment variable:")
		pterm.DefaultBox.WithTitle("Environment Variable").Println(fmt.Sprintf("export BASTION_JWT_SECRET=%s", encodedSecret))

		return nil
	},
}

func init() {
	createJWTSecretCmd.Flags().BoolVar(&printOnly, "print-only", false, "Only print the secret (useful for scripting)")
	createCmd.AddCommand(createJWTSecretCmd)
}
