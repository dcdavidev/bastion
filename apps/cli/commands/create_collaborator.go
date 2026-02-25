package commands

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/packages/core/crypto"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	collabUsername string
	collabPassword string
	collabProject  string
)

var createCollaboratorCmd = &cobra.Command{
	Use:   "create-collaborator",
	Short: "Create a restricted collaborator with project-specific access",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		adminPwd, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Admin Password (to unlock vault)")

		pterm.DefaultSection.Println("Collaborator Credentials")
		if collabUsername == "" {
			collabUsername, _ = pterm.DefaultInteractiveTextInput.Show("Enter Collaborator Username")
		}
		if collabPassword == "" {
			collabPassword, _ = pterm.DefaultInteractiveTextInput.WithMask("*").Show("Enter Collaborator Password")
		}

		spinner, _ := pterm.DefaultSpinner.Start("Re-wrapping project keys for collaborator...")

		// 1. Fetch Vault Config
		vaultConfig, err := fetchVaultConfig(serverURL, token)
		if err != nil {
			spinner.Fail("Failed to fetch vault config")
			return err
		}

		// 2. Fetch Project (Admin's view)
		project, err := fetchProject(serverURL, token, collabProject)
		if err != nil {
			spinner.Fail("Project not found")
			return err
		}

		// 3. Unwrap Admin's Master Key
		adminSalt, _ := hex.DecodeString(vaultConfig.MasterKeySalt)
		wrappedMK, _ := hex.DecodeString(vaultConfig.WrappedMasterKey)
		adminKEK := crypto.DeriveKey([]byte(adminPwd), adminSalt)
		masterKey, err := crypto.UnwrapKey(adminKEK, wrappedMK)
		if err != nil {
			spinner.Fail("Failed to unlock admin vault: invalid password")
			return err
		}

		// 4. Unwrap Project Data Key
		wrappedDK, _ := hex.DecodeString(project.WrappedDataKey)
		dataKey, err := crypto.UnwrapKey(masterKey, wrappedDK)
		if err != nil {
			spinner.Fail("Failed to unwrap project data key")
			return err
		}

		// 5. Generate Collaborator material (Salt and KEK)
		collabSalt, err := crypto.GenerateSalt()
		if err != nil {
			spinner.Fail("Failed to generate salt for collaborator")
			return err
		}
		collabKEK := crypto.DeriveKey([]byte(collabPassword), collabSalt)

		// 6. Wrap Data Key for Collaborator (using their KEK)
		collabWrappedKey, err := crypto.WrapKey(collabKEK, dataKey)
		if err != nil {
			spinner.Fail("Failed to wrap key for collaborator")
			return err
		}

		spinner.UpdateText("Sending credentials to server...")

		// 7. Send to server
		payload, _ := json.Marshal(map[string]interface{}{
			"username":         collabUsername,
			"password_hash":    hex.EncodeToString(collabKEK), // Using KEK as hash for simple auth
			"salt":             hex.EncodeToString(collabSalt),
			"project_id":       collabProject,
			"wrapped_data_key": hex.EncodeToString(collabWrappedKey),
		})

		req, _ := http.NewRequest("POST", serverURL+"/api/v1/collaborators", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Connection error")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			spinner.Fail(fmt.Sprintf("Failed to create collaborator: %s", resp.Status))
			return fmt.Errorf("failed to create collaborator: %s", resp.Status)
		}

		spinner.Success(fmt.Sprintf("Collaborator '%s' created successfully!", collabUsername))
		pterm.Info.Printf("Project ID: %s\n", collabProject)
		pterm.Info.Println("The collaborator can now use 'bastion run' to access secrets.")
		return nil
	},
}

func init() {
	createCollaboratorCmd.Flags().StringVarP(&collabUsername, "username", "n", "", "Collaborator username")
	createCollaboratorCmd.Flags().StringVarP(&collabPassword, "password", "p", "", "Collaborator password")
	createCollaboratorCmd.Flags().StringVarP(&collabProject, "project", "j", "", "Project ID (UUID)")
	createCollaboratorCmd.MarkFlagRequired("project")
	rootCmd.AddCommand(createCollaboratorCmd)
}
