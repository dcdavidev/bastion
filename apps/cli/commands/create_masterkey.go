package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createMasterKeyCmd = &cobra.Command{
	Use:   "masterkey",
	Short: "Generate a new random master key",
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("failed to generate random key: %w", err)
		}

		encodedKey := hex.EncodeToString(key)

		pterm.Success.Println("Generated new random master key")
		pterm.Info.Printf("MASTER_KEY=%s\n", encodedKey)

		return nil
	},
}

func init() {
	createCmd.AddCommand(createMasterKeyCmd)
}
