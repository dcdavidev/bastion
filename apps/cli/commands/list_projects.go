package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listProjectsClientID string

var listProjectsCmd = &cobra.Command{
        Use:   "projects",
        Short: "List all projects for a specific client",
        RunE: func(cmd *cobra.Command, args []string) error {
                serverURL, _ := cmd.Flags().GetString("url")

                token, err := loadToken()
                if err != nil {
                        return fmt.Errorf("not authenticated: %w", err)
                }

                clientID := listProjectsClientID
                if clientID == "" {
                        var err error
                        clientID, err = pterm.DefaultInteractiveTextInput.Show("Client ID (UUID)")
                        if err != nil {
                                return err
                        }
                }

                if clientID == "" {
                        return fmt.Errorf("client ID is required")
                }

                spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Fetching projects for client %s...", clientID))
                // ... (rest of the function remains the same)
		req, _ := http.NewRequest("GET", serverURL+"/api/v1/projects?client_id="+clientID, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Connection error")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			spinner.Fail(fmt.Sprintf("Failed to fetch projects: %s", resp.Status))
			return fmt.Errorf("failed to fetch projects: %s", resp.Status)
		}

		var projects []models.Project
		json.NewDecoder(resp.Body).Decode(&projects)

		spinner.Success("Projects fetched successfully!")

		if len(projects) == 0 {
			pterm.Info.Println("No projects found for this client.")
			return nil
		}

		tableData := pterm.TableData{
			{"ID", "Name", "Created At"},
		}

		for _, p := range projects {
			tableData = append(tableData, []string{
				p.ID.String(),
				p.Name,
				p.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		return nil
	},
}

func init() {
        listProjectsCmd.Flags().StringVarP(&listProjectsClientID, "client", "c", "", "Client ID to list projects for")
        listCmd.AddCommand(listProjectsCmd)
}

