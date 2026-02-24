package commands

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/internal/crypto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	setKey   string
	setValue string
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Encrypt and store a secret in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Admin Password to unlock vault")

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Encrypting and storing secret '%s'...", setKey))

		// 1. Fetch Vault Config and Project Data Key
		vaultConfig, err := fetchVaultConfig(serverURL, token)
		if err != nil {
			spinner.Fail("Failed to fetch vault config")
			return err
		}

		project, err := fetchProject(serverURL, token, projectID)
		if err != nil {
			spinner.Fail("Project not found")
			return err
		}

		// 2. Unwrap Master Key
		salt, _ := hex.DecodeString(vaultConfig.MasterKeySalt)
		wrappedMK, _ := hex.DecodeString(vaultConfig.WrappedMasterKey)
		adminKEK := crypto.DeriveKey([]byte(password), salt)
		
		masterKey, err := crypto.UnwrapKey(adminKEK, wrappedMK)
		if err != nil {
			spinner.Fail("Failed to unlock vault: invalid password")
			return fmt.Errorf("failed to unlock vault: invalid password")
		}

		// 3. Unwrap Project Data Key
		wrappedDK, _ := hex.DecodeString(project.WrappedDataKey)
		dataKey, err := crypto.UnwrapKey(masterKey, wrappedDK)
		if err != nil {
			spinner.Fail("Failed to unwrap project data key")
			return fmt.Errorf("failed to unwrap project data key: %w", err)
		}

		// 4. Encrypt Secret Value
		ciphertext, err := crypto.Encrypt(dataKey, []byte(setValue))
		if err != nil {
			spinner.Fail("Encryption failed")
			return fmt.Errorf("encryption failed: %w", err)
		}

		spinner.UpdateText("Uploading encrypted secret to server...")

		// 5. Upload to Server
		payload, _ := json.Marshal(map[string]string{
			"project_id": projectID,
			"key":        setKey,
			"value":      hex.EncodeToString(ciphertext),
		})

		req, _ := http.NewRequest("POST", serverURL+"/api/v1/secrets", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Connection error")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			spinner.Fail(fmt.Sprintf("Failed to store secret: %s", resp.Status))
			return fmt.Errorf("failed to store secret: %s", resp.Status)
		}

		spinner.Success(fmt.Sprintf("Secret '%s' stored successfully in project %s", setKey, projectID))
		return nil
	},
}

func init() {
	setCmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID")
	setCmd.Flags().StringVarP(&setKey, "key", "k", "", "Secret key name")
	setCmd.Flags().StringVarP(&setValue, "value", "v", "", "Secret value")
	
	setCmd.MarkFlagRequired("project")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("value")
	
	rootCmd.AddCommand(setCmd)
}
