package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dcdavidev/bastion/internal/crypto"
	"github.com/dcdavidev/bastion/internal/models"
	"github.com/spf13/cobra"
)

var projectID string

var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Inject secrets into a command's environment",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated. Please run 'bastion login' first: %w", err)
		}

		fmt.Print("Enter Admin Password to unlock vault: ")
		var password string
		fmt.Scanln(&password)

		// 1. Fetch Vault Config (Wrapped Master Key and Salt)
		vaultConfig, err := fetchVaultConfig(serverURL, token)
		if err != nil {
			return err
		}

		// 2. Derive Admin KEK and Unwrap Master Key
		salt, _ := hex.DecodeString(vaultConfig.MasterKeySalt)
		wrappedMK, _ := hex.DecodeString(vaultConfig.WrappedMasterKey)
		adminKEK := crypto.DeriveKey([]byte(password), salt)
		
		masterKey, err := crypto.UnwrapKey(adminKEK, wrappedMK)
		if err != nil {
			return fmt.Errorf("failed to unlock vault: invalid password or corrupted master key")
		}

		// 3. Fetch Project and its Wrapped Data Key
		project, err := fetchProject(serverURL, token, projectID)
		if err != nil {
			return err
		}

		wrappedDK, _ := hex.DecodeString(project.WrappedDataKey)
		dataKey, err := crypto.UnwrapKey(masterKey, wrappedDK)
		if err != nil {
			return fmt.Errorf("failed to unwrap project data key: %w", err)
		}

		// 4. Fetch and Decrypt Secrets
		encryptedSecrets, err := fetchEncryptedSecrets(serverURL, token, projectID)
		if err != nil {
			return err
		}

		env := os.Environ()
		for _, s := range encryptedSecrets {
			ciphertext, _ := hex.DecodeString(s.Value)
			plaintext, err := crypto.Decrypt(dataKey, ciphertext)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to decrypt secret %s: %v\n", s.Key, err)
				continue
			}
			env = append(env, fmt.Sprintf("%s=%s", s.Key, string(plaintext)))
		}

		// 5. Execute command
		externalCmd := exec.Command(args[0], args[1:]...)
		externalCmd.Env = env
		externalCmd.Stdout = os.Stdout
		externalCmd.Stderr = os.Stderr
		externalCmd.Stdin = os.Stdin

		return externalCmd.Run()
	},
}

func fetchEncryptedSecrets(url, token, projectID string) ([]models.Secret, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/secrets?project_id="+projectID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch secrets: %s", resp.Status)
	}

	var secrets []models.Secret
	json.NewDecoder(resp.Body).Decode(&secrets)
	return secrets, nil
}

func init() {
	runCmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID to fetch secrets from")
	runCmd.MarkFlagRequired("project")
	rootCmd.AddCommand(runCmd)
}
