package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createClientName string

var createClientCmd = &cobra.Command{
	Use:   "create-client",
	Short: "Create a new client in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Creating client '%s'...", createClientName))

		// Send to server
		payload, _ := json.Marshal(map[string]string{
			"name": createClientName,
		})

		req, _ := http.NewRequest("POST", serverURL+"/api/v1/clients", bytes.NewBuffer(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Connection error")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			spinner.Fail(fmt.Sprintf("Failed to create client: %s", resp.Status))
			return fmt.Errorf("failed to create client: %s", resp.Status)
		}

		var client struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		json.NewDecoder(resp.Body).Decode(&client)

		spinner.Success(fmt.Sprintf("Client '%s' created successfully!", client.Name))
		pterm.Info.Printf("ID: %s\n", client.ID)
		return nil
	},
}

func init() {
	createClientCmd.Flags().StringVarP(&createClientName, "name", "n", "", "Client name")
	createClientCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(createClientCmd)
}
