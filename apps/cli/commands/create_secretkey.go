package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createSecretKeyCmd = &cobra.Command{
	Use:   "secretkey",
	Short: "Generate an ultra-secure random JWT secret key",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("SECRET KEY GENERATION")

		secret := make([]byte, 64) // 512-bit key
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate random secret: %w", err)
		}

		encodedSecret := hex.EncodeToString(secret)

		pterm.Success.Println("Generated new ultra-secure secret key!")

		// Attempt to save locally in .env
		envPath := ".env"
		content, err := os.ReadFile(envPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			updated := false
			for i, line := range lines {
				if strings.HasPrefix(line, "BASTION_JWT_SECRET=") {
					lines[i] = "BASTION_JWT_SECRET=" + encodedSecret
					updated = true
					break
				}
			}

			if !updated {
				lines = append(lines, "BASTION_JWT_SECRET="+encodedSecret)
			}

			err = os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644)
			if err == nil {
				pterm.Success.Println("Saved BASTION_JWT_SECRET to local .env file!")
			} else {
				pterm.Warning.Printf("Could not write to .env file: %v\n", err)
			}
		} else {
			pterm.Info.Println("No .env file found in current directory. Creating one...")
			err = os.WriteFile(envPath, []byte("BASTION_JWT_SECRET="+encodedSecret+"\n"), 0644)
			if err == nil {
				pterm.Success.Println("Created .env and saved BASTION_JWT_SECRET!")
			}
		}

		pterm.Info.Println("To apply remotely, please ensure this variable is updated on your Bastion server environment.")
		pterm.DefaultBox.WithTitle("New Secret Key").Println(encodedSecret)

		return nil
	},
}

func init() {
	createCmd.AddCommand(createSecretKeyCmd)
}
