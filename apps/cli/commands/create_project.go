package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new project with a dedicated encrypted data key",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientID, _ := pterm.DefaultInteractiveTextInput.Show("Client ID (UUID)")
		projectName, _ := pterm.DefaultInteractiveTextInput.Show("Project Name")

		url, _ := cmd.Flags().GetString("url")
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}

		payload := map[string]string{
			"client_id":    clientID,
			"name":         projectName,
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url+"/api/v1/projects", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create project: status %d", resp.StatusCode)
		}

		pterm.Success.Printf("Project %s created successfully!\n", projectName)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createProjectCmd)
}
