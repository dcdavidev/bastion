package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var clientNameFlag string

var createClientCmd = &cobra.Command{
        Use:   "client",
        Short: "Create a new client in the vault",
        RunE: func(cmd *cobra.Command, args []string) error {
                clientName := clientNameFlag
                if clientName == "" {
                        var err error
                        clientName, err = pterm.DefaultInteractiveTextInput.Show("Client Name")
                        if err != nil {
                                return err
                        }
                }

                if clientName == "" {
                        return fmt.Errorf("client name is required")
                }

                url, _ := cmd.Flags().GetString("url")
                token, err := loadToken()
                if err != nil {
                        return fmt.Errorf("not logged in: %w", err)
                }

                spinner, _ := pterm.DefaultSpinner.Start("Creating client " + clientName + "...")

                payload := map[string]string{
                        "name": clientName,
                }

                jsonPayload, _ := json.Marshal(payload)
                req, _ := http.NewRequest("POST", url+"/api/v1/clients", bytes.NewBuffer(jsonPayload))
                req.Header.Set("Authorization", "Bearer "+token)
                req.Header.Set("Content-Type", "application/json")

                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                        spinner.Fail(err.Error())
                        return err
                }
                defer resp.Body.Close()

                if resp.StatusCode != http.StatusCreated {
                        spinner.Fail(fmt.Sprintf("Failed to create client: status %d", resp.StatusCode))
                        return fmt.Errorf("failed to create client: status %d", resp.StatusCode)
                }

                spinner.Success(fmt.Sprintf("Client %s created successfully!", clientName))
                return nil
        },
}

func init() {
        createClientCmd.Flags().StringVarP(&clientNameFlag, "name", "n", "", "Name of the client")
        createCmd.AddCommand(createClientCmd)
}

