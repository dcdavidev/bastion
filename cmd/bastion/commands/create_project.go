package commands

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	createProjectName     string
	createProjectClientID string
)

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new project with a dedicated encrypted data key",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		fmt.Print("Enter Admin Password to unlock vault: ")
		var password string
		fmt.Scanln(&password)

		// 1. Fetch Vault Config
		vaultConfig, err := fetchVaultConfig(serverURL, token)
		if err != nil {
			return err
		}

		// 2. Unwrap Master Key
		salt, _ := hex.DecodeString(vaultConfig.MasterKeySalt)
		wrappedMK, _ := hex.DecodeString(vaultConfig.WrappedMasterKey)
		adminKEK := crypto.DeriveKey([]byte(password), salt)
		
		masterKey, err := crypto.UnwrapKey(adminKEK, wrappedMK)
		if err != nil {
			return fmt.Errorf("failed to unlock vault: invalid password")
		}

		// 3. Generate and Wrap a new Data Key for the project
		dataKey, err := crypto.GenerateRandomKey()
		if err != nil {
			return err
		}

		wrappedDK, err := crypto.WrapKey(masterKey, dataKey)
		if err != nil {
			return err
		}

		// 4. Send to server
		payload, _ := json.Marshal(map[string]string{
			"client_id":        createProjectClientID,
			"name":             createProjectName,
			"wrapped_data_key": hex.EncodeToString(wrappedDK),
		})

		req, _ := http.NewRequest("POST", serverURL+"/api/v1/projects", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create project: %s", resp.Status)
		}

		fmt.Printf("Project '%s' created successfully with its own encrypted data key.
", createProjectName)
		return nil
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&createProjectName, "name", "n", "", "Project name")
	createProjectCmd.Flags().StringVarP(&createProjectClientID, "client", "c", "", "Client UUID")
	
	createProjectCmd.MarkFlagRequired("name")
	createProjectCmd.MarkFlagRequired("client")
	
	rootCmd.AddCommand(createProjectCmd)
}
