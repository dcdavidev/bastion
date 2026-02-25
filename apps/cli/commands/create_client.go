package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create a new client in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientName, _ := pterm.DefaultInteractiveTextInput.Show("Client Name")

		url, _ := cmd.Flags().GetString("url")
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}

		payload := map[string]string{
			"name": clientName,
		}

		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url+"/api/v1/clients", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create client: status %d", resp.StatusCode)
		}

		pterm.Success.Printf("Client %s created successfully!\n", clientName)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createClientCmd)
}
