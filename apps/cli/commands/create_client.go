package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var clientName string

var createClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create a new client",
	RunE: func(cmd *cobra.Command, args []string) error {
		if activeProfile == nil || activeProfile.URL == "" {
			return fmt.Errorf("no active profile. Please login first")
		}

		if clientName == "" {
			var err error
			clientName, err = pterm.DefaultInteractiveTextInput.Show("Enter Client Name")
			if err != nil {
				return err
			}
		}

		// Validation: [a-z0-9\-]
		reg := regexp.MustCompile(`^[a-z0-9\-]+$`)
		if !reg.MatchString(clientName) {
			return fmt.Errorf("invalid client name: must contain only lowercase letters, numbers, and hyphens")
		}

		spinner, _ := pterm.DefaultSpinner.Start("Creating client...")

		payload, _ := json.Marshal(map[string]string{
			"name": clientName,
		})

		req, err := http.NewRequest("POST", activeProfile.URL+"/api/v1/clients", bytes.NewBuffer(payload))
		if err != nil {
			spinner.Fail("Failed to create request")
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		if activeProfile.Token != "" {
			req.Header.Set("Authorization", "Bearer "+activeProfile.Token)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Failed to connect to server")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			spinner.Fail(fmt.Sprintf("Failed to create client: %s", resp.Status))
			return fmt.Errorf("server returned status: %s", resp.Status)
		}

		spinner.Success(fmt.Sprintf("Client '%s' created successfully!", clientName))
		return nil
	},
}

func init() {
	createClientCmd.Flags().StringVarP(&clientName, "name", "n", "", "Name of the client ([a-z0-9\\-])")
	createCmd.AddCommand(createClientCmd)
}
