package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listClientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List all clients in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated: %w", err)
		}

		spinner, _ := pterm.DefaultSpinner.Start("Fetching clients...")

		req, _ := http.NewRequest("GET", serverURL+"/api/v1/clients", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			spinner.Fail("Connection error")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			spinner.Fail(fmt.Sprintf("Failed to fetch clients: %s", resp.Status))
			return fmt.Errorf("failed to fetch clients: %s", resp.Status)
		}

		var clients []models.Client
		json.NewDecoder(resp.Body).Decode(&clients)

		spinner.Success("Clients fetched successfully!")

		if len(clients) == 0 {
			pterm.Info.Println("No clients found in the vault.")
			return nil
		}

		tableData := pterm.TableData{
			{"ID", "Name", "Created At"},
		}

		for _, c := range clients {
			tableData = append(tableData, []string{
				c.ID.String(),
				c.Name,
				c.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		return nil
	},
}

func init() {
	listCmd.AddCommand(listClientsCmd)
}
