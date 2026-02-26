package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createProjectClient string
var createProjectName string

var createProjectCmd = &cobra.Command{
        Use:   "project",
        Short: "Create a new project with a dedicated encrypted data key",
        RunE: func(cmd *cobra.Command, args []string) error {
                clientID := createProjectClient
                if clientID == "" {
                        var err error
                        clientID, err = pterm.DefaultInteractiveTextInput.Show("Client ID (UUID)")
                        if err != nil {
                                return err
                        }
                }

                projectName := createProjectName
                if projectName == "" {
                        var err error
                        projectName, err = pterm.DefaultInteractiveTextInput.Show("Project Name")
                        if err != nil {
                                return err
                        }
                }

                if clientID == "" || projectName == "" {
                        return fmt.Errorf("client ID and project name are required")
                }

                url, _ := cmd.Flags().GetString("url")
                token, err := loadToken()
                if err != nil {
                        return fmt.Errorf("not logged in: %w", err)
                }

                spinner, _ := pterm.DefaultSpinner.Start("Creating project " + projectName + "...")

                payload := map[string]string{
                        "client_id": clientID,
                        "name":      projectName,
                }

                jsonPayload, _ := json.Marshal(payload)
                req, _ := http.NewRequest("POST", url+"/api/v1/projects", bytes.NewBuffer(jsonPayload))
                req.Header.Set("Authorization", "Bearer "+token)
                req.Header.Set("Content-Type", "application/json")

                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                        spinner.Fail(err.Error())
                        return err
                }
                defer resp.Body.Close()

                if resp.StatusCode != http.StatusCreated {
                        spinner.Fail(fmt.Sprintf("Failed to create project: status %d", resp.StatusCode))
                        return fmt.Errorf("failed to create project: status %d", resp.StatusCode)
                }

                spinner.Success(fmt.Sprintf("Project %s created successfully!", projectName))
                return nil
        },
}

func init() {
        createProjectCmd.Flags().StringVarP(&createProjectClient, "client", "c", "", "Client ID (UUID)")
        createProjectCmd.Flags().StringVarP(&createProjectName, "name", "n", "", "Name of the project")
        createCmd.AddCommand(createProjectCmd)
}

