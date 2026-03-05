package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/packages/db"
	"github.com/dcdavidev/bastion/packages/models"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	removeClientID  string
	removeProjectID string
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove resources",
}

var removeClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Remove a client and all its associated projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if removeClientID == "" {
			var err error
			removeClientID, err = pterm.DefaultInteractiveTextInput.Show("Enter Client ID or Name to remove")
			if err != nil {
				return err
			}
		}

		confirm, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show(fmt.Sprintf("Are you sure you want to remove client '%s'? THIS ACTION IS IRREVERSIBLE!", removeClientID))
		if !confirm {
			pterm.Info.Println("Operation cancelled.")
			return nil
		}

		spinner, _ := pterm.DefaultSpinner.Start("Removing client...")

		isRemote := activeProfile != nil && activeProfile.Token != "" && activeProfile.URL != ""

		var targetID uuid.UUID
		if isRemote {
			// Resolve name to ID if necessary
			uid, err := uuid.Parse(removeClientID)
			if err != nil {
				// Fetch clients to resolve name
				req, _ := http.NewRequest("GET", activeProfile.URL+"/api/v1/clients", nil)
				req.Header.Set("Authorization", "Bearer "+activeProfile.Token)
				resp, err := http.DefaultClient.Do(req)
				if err == nil && resp.StatusCode == http.StatusOK {
					var clients []models.Client
					json.NewDecoder(resp.Body).Decode(&clients)
					resp.Body.Close()
					for _, c := range clients {
						if c.Name == removeClientID {
							uid = c.ID
							break
						}
					}
				}
			}
			targetID = uid

			if targetID == uuid.Nil {
				spinner.Fail("Client not found")
				return fmt.Errorf("client not found")
			}

			// Remote Call
			req, _ := http.NewRequest("DELETE", activeProfile.URL+"/api/v1/clients/"+targetID.String(), nil)
			req.Header.Set("Authorization", "Bearer "+activeProfile.Token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusNoContent {
				spinner.Fail("Failed to remove client on server")
				return fmt.Errorf("api error")
			}
		} else {
			// Local Mode
			database, _ := db.NewConnection()
			defer database.Close()
			uid, err := uuid.Parse(removeClientID)
			if err != nil {
				clients, _ := database.GetClients(context.Background())
				for _, c := range clients {
					if c.Name == removeClientID {
						uid = c.ID
						break
					}
				}
			}
			targetID = uid
			if targetID == uuid.Nil {
				spinner.Fail("Client not found")
				return fmt.Errorf("client not found")
			}
			err = database.DeleteClient(context.Background(), targetID)
			if err != nil {
				spinner.Fail("Failed to remove client from database")
				return err
			}
		}

		spinner.Success(fmt.Sprintf("Client '%s' and all associated data removed.", removeClientID))
		return nil
	},
}

var removeProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Remove a project and all its secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if removeProjectID == "" {
			var err error
			removeProjectID, err = pterm.DefaultInteractiveTextInput.Show("Enter Project ID or Name to remove")
			if err != nil {
				return err
			}
		}

		confirm, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show(fmt.Sprintf("Are you sure you want to remove project '%s'? THIS ACTION IS IRREVERSIBLE!", removeProjectID))
		if !confirm {
			pterm.Info.Println("Operation cancelled.")
			return nil
		}

		spinner, _ := pterm.DefaultSpinner.Start("Removing project...")

		isRemote := activeProfile != nil && activeProfile.Token != "" && activeProfile.URL != ""

		var targetID uuid.UUID
		if isRemote {
			// Resolve name to ID would require fetching all projects for all clients,
			// for now let's assume ID or handle only if we have client context.
			// Simplified: assume UUID or fail for now in remote name resolution without client_id.
			uid, err := uuid.Parse(removeProjectID)
			if err != nil {
				spinner.Fail("Project removal by name requires UUID in remote mode currently")
				return fmt.Errorf("please use UUID")
			}
			targetID = uid

			req, _ := http.NewRequest("DELETE", activeProfile.URL+"/api/v1/projects/"+targetID.String(), nil)
			req.Header.Set("Authorization", "Bearer "+activeProfile.Token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusNoContent {
				spinner.Fail("Failed to remove project on server")
				return fmt.Errorf("api error")
			}
		} else {
			database, _ := db.NewConnection()
			defer database.Close()
			uid, err := uuid.Parse(removeProjectID)
			if err != nil {
				// Local name resolution is easier
				query := `SELECT id FROM projects WHERE name = $1 LIMIT 1`
				database.Pool.QueryRow(context.Background(), query, removeProjectID).Scan(&uid)
			}
			targetID = uid
			if targetID == uuid.Nil {
				spinner.Fail("Project not found")
				return fmt.Errorf("not found")
			}
			err = database.DeleteProject(context.Background(), targetID)
			if err != nil {
				spinner.Fail("Failed to remove project from database")
				return err
			}
		}

		spinner.Success(fmt.Sprintf("Project '%s' removed.", removeProjectID))
		return nil
	},
}

func init() {
	removeCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return removeInteractive()
	}
	removeClientCmd.Flags().StringVarP(&removeClientID, "id", "i", "", "Client ID or Name")
	removeProjectCmd.Flags().StringVarP(&removeProjectID, "id", "i", "", "Project ID or Name")

	removeCmd.AddCommand(removeClientCmd)
	removeCmd.AddCommand(removeProjectCmd)
	rootCmd.AddCommand(removeCmd)
}
