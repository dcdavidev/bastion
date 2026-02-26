package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
        runProjectID string
        runPassword  string
)

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

                password := runPassword
                if password == "" {
                        var err error
                        password, err = pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Password to unlock secrets")
                        if err != nil {
                                return err
                        }
                }

                spinner, _ := pterm.DefaultSpinner.Start("Fetching project credentials...")

                // 1. Fetch Vault Config (Salt)
                vaultConfig, err := fetchVaultConfig(serverURL, token)
                if err != nil {
                        spinner.Fail("Failed to fetch vault config")
                        return err
                }

                // 2. Fetch User-Specific Wrapped Data Key
                wrappedDK, err := fetchUserProjectKey(serverURL, token, runProjectID)
                if err != nil {
                        spinner.Fail("Access denied or project not found")
                        return err
                }
		// 3. Derive KEK and Unwrap Data Key
		salt, _ := hex.DecodeString(vaultConfig.MasterKeySalt)
		userKEK := crypto.DeriveKey([]byte(password), salt)
		
		// Note: For collaborators, the Data Key is wrapped by their KEK.
		// For Admins, we need to unwrap Master Key first.
		// Optimization: The server endpoint should ideally return the key wrapped for the current user.
		
		// Attempt direct unwrap (Collaborator flow)
		wrappedDKBytes, _ := hex.DecodeString(wrappedDK)
		dataKey, err := crypto.UnwrapKey(userKEK, wrappedDKBytes)
		
		if err != nil {
			// Try Admin flow: Unwrap Master Key first
			wrappedMK, _ := hex.DecodeString(vaultConfig.WrappedMasterKey)
			masterKey, errMK := crypto.UnwrapKey(userKEK, wrappedMK)
			if errMK != nil {
				spinner.Fail("Invalid password")
				return fmt.Errorf("invalid password")
			}
			dataKey, err = crypto.UnwrapKey(masterKey, wrappedDKBytes)
			if err != nil {
				spinner.Fail("Failed to unwrap project key")
				return err
			}
		}

		                // 4. Fetch and Decrypt Secrets
		                encryptedSecrets, err := fetchEncryptedSecrets(serverURL, token, runProjectID)
		                if err != nil {
		                        spinner.Fail("Failed to fetch secrets")
		                        return err
		                }
		
		                env := os.Environ()
		
		decryptedCount := 0
		for _, s := range encryptedSecrets {
			ciphertext, _ := hex.DecodeString(s.Value)
			plaintext, err := crypto.Decrypt(dataKey, ciphertext)
			if err != nil {
				pterm.Warning.Printf("Failed to decrypt secret %s: %v\n", s.Key, err)
				continue
			}
			env = append(env, fmt.Sprintf("%s=%s", s.Key, string(plaintext)))
			decryptedCount++
		}

		spinner.Success(fmt.Sprintf("Injected %d secrets into environment.", decryptedCount))

		// 5. Execute command
		pterm.Info.Printf("Executing: %v\n\n", args)
		externalCmd := exec.Command(args[0], args[1:]...)
		externalCmd.Env = env
		externalCmd.Stdout = os.Stdout
		externalCmd.Stderr = os.Stderr
		externalCmd.Stdin = os.Stdin

		return externalCmd.Run()
	},
}

func fetchUserProjectKey(url, token, id string) (string, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/projects/"+id+"/key", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch project key: %s", resp.Status)
	}

	var data struct {
		WrappedDataKey string `json:"wrapped_data_key"`
	}
	json.NewDecoder(resp.Body).Decode(&data)
	return data.WrappedDataKey, nil
}

func fetchEncryptedSecrets(url, token, projectID string) ([]models.Secret, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/secrets?project_id="+projectID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var secrets []models.Secret
	json.NewDecoder(resp.Body).Decode(&secrets)
	return secrets, nil
}

func init() {
        runCmd.Flags().StringVarP(&runProjectID, "project", "p", "", "Project ID to fetch secrets from")
        runCmd.Flags().StringVar(&runPassword, "password", "", "Password to unlock vault")
        runCmd.MarkFlagRequired("project")
        rootCmd.AddCommand(runCmd)
}

