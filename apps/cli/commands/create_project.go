package commands

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dcdavidev/bastion/packages/crypto"
	"github.com/dcdavidev/bastion/packages/db"
	"github.com/dcdavidev/bastion/packages/models"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	projectClient string
	projectName   string
)

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new end-to-end encrypted project",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Ensure we have a profile or local access
		if activeProfile == nil && os.Getenv("BASTION_DATABASE_URL") == "" {
			return fmt.Errorf("no active profile. Please login first or set BASTION_DATABASE_URL")
		}

		// 2. Gather Inputs
		if projectClient == "" {
			var err error
			projectClient, err = pterm.DefaultInteractiveTextInput.Show("Enter Client ID or Name")
			if err != nil {
				return err
			}
		}

		if projectName == "" {
			var err error
			projectName, err = pterm.DefaultInteractiveTextInput.WithDefaultText("default").Show("Enter Project Name")
			if err != nil {
				return err
			}
		}

		password, err := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Admin Password to unwrap Master Key")
		if err != nil {
			return err
		}

		spinner, _ := pterm.DefaultSpinner.Start("Preparing encryption keys...")

		// 3. Mode Detection
		// If we have an active profile with a token, we use the API (Remote mode)
		// even if the URL is localhost. Otherwise we fallback to local DB access.
		isRemote := activeProfile != nil && activeProfile.Token != "" && activeProfile.URL != ""

		var vaultSalt []byte
		var wrappedMK []byte

		if isRemote {
			// Remote Mode
			spinner.UpdateText("Fetching vault configuration from server...")
			req, _ := http.NewRequest("GET", activeProfile.URL+"/api/v1/vault/config", nil)
			req.Header.Set("Authorization", "Bearer "+activeProfile.Token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				spinner.Fail("Failed to fetch vault configuration from server")
				return fmt.Errorf("vault config error")
			}
			defer resp.Body.Close()
			var vc struct {
				WrappedMasterKey string `json:"wrapped_master_key"`
				MasterKeySalt    string `json:"master_key_salt"`
			}
			json.NewDecoder(resp.Body).Decode(&vc)
			vaultSalt, _ = hex.DecodeString(vc.MasterKeySalt)
			wrappedMK, _ = hex.DecodeString(vc.WrappedMasterKey)
		} else {
			// Local Mode
			spinner.UpdateText("Connecting to local database...")
			database, err := db.NewConnection()
			if err != nil {
				spinner.Fail("Local database connection failed. Ensure your database is running and BASTION_DATABASE_URL is set.")
				return err
			}
			defer database.Close()
			vc, err := database.GetVaultConfig(context.Background())
			if err != nil {
				spinner.Fail("Vault not initialized")
				return err
			}
			vaultSalt, _ = hex.DecodeString(vc.MasterKeySalt)
			wrappedMK, _ = hex.DecodeString(vc.WrappedMasterKey)
		}

		// Derive KEK and unwrap Master Key
		spinner.UpdateText("Deriving keys and unwrapping Master Key...")
		kek := crypto.DeriveKey([]byte(password), vaultSalt)
		masterKey, err := crypto.UnwrapKey(kek, wrappedMK)
		if err != nil {
			spinner.Fail("Failed to unwrap Master Key. Invalid password?")
			return err
		}

		// Generate new Data Key for the project and wrap it with Master Key
		dataKey, _ := crypto.GenerateRandomKey()
		wrappedDK, _ := crypto.WrapKey(masterKey, dataKey)
		wrappedDKHex := hex.EncodeToString(wrappedDK)

		// 4. Resolve Client ID if name provided
		clientID := projectClient
		uid, err := uuid.Parse(projectClient)
		if err != nil {
			// It's a name, we need to find the UUID
			spinner.UpdateText("Resolving client name...")
			var clients []models.Client
			if isRemote {
				req, _ := http.NewRequest("GET", activeProfile.URL+"/api/v1/clients", nil)
				req.Header.Set("Authorization", "Bearer "+activeProfile.Token)
				resp, err := http.DefaultClient.Do(req)
				if err == nil && resp.StatusCode == http.StatusOK {
					json.NewDecoder(resp.Body).Decode(&clients)
					resp.Body.Close()
				}
			} else {
				database, _ := db.NewConnection()
				clients, _ = database.GetClients(context.Background())
				database.Close()
			}

			found := false
			for _, c := range clients {
				if c.Name == projectClient {
					clientID = c.ID.String()
					found = true
					break
				}
			}
			if !found {
				spinner.Fail(fmt.Sprintf("Client '%s' not found", projectClient))
				return fmt.Errorf("client not found")
			}
		} else {
			clientID = uid.String()
		}

		// 5. API Request
		if isRemote {
			spinner.UpdateText("Creating project on remote server...")
			// Remote API Call
			payload, _ := json.Marshal(map[string]interface{}{
				"client_id":        clientID,
				"name":             projectName,
				"wrapped_data_key": wrappedDKHex,
			})

			req, _ := http.NewRequest("POST", activeProfile.URL+"/api/v1/projects", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+activeProfile.Token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				spinner.Fail(fmt.Sprintf("Failed to create project on server: %s", string(body)))
				return fmt.Errorf("api error")
			}
		} else {
			spinner.UpdateText("Creating project in local database...")
			// Local DB Insert
			database, _ := db.NewConnection()
			defer database.Close()
			cid, _ := uuid.Parse(clientID)
			_, err = database.CreateProject(context.Background(), cid, projectName, wrappedDKHex)
			if err != nil {
				spinner.Fail("Failed to create project in database")
				return err
			}
		}

		spinner.Success(fmt.Sprintf("Project '%s' created and secured with E2EE!", projectName))
		return nil
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&projectClient, "client", "c", "", "Client ID or Name")
	createProjectCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project Name")
	createCmd.AddCommand(createProjectCmd)
}
