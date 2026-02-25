package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createCollaboratorCmd = &cobra.Command{
	Use:     "collaborator",
	Aliases: []string{"collab"},
	Short:   "Create a restricted collaborator with project-specific access",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := pterm.DefaultInteractiveTextInput.Show("Collaborator Username")
		password, _ := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Collaborator Password")
		role, _ := pterm.DefaultInteractiveTextInput.Show("Role (admin/collaborator)")

		url, _ := cmd.Flags().GetString("url")
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}

		payload := map[string]string{
			"username": username,
			"password": password,
			"role":     role,
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url+"/api/v1/collaborators", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create collaborator: status %d", resp.StatusCode)
		}

		pterm.Success.Printf("Collaborator %s created successfully!\n", username)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createCollaboratorCmd)
}
