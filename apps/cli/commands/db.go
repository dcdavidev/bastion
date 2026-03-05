package commands

import (
	"fmt"

	"github.com/dcdavidev/bastion/packages/db"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Check and apply pending database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		pterm.DefaultHeader.WithFullWidth().Println("DATABASE MIGRATIONS")

		spinner, _ := pterm.DefaultSpinner.Start("Connecting to database...")
		database, err := db.NewConnection()
		if err != nil {
			spinner.Fail("Connection failed: " + err.Error())
			return err
		}
		defer database.Close()
		spinner.Success("Connected to database!")

		// 1. Check status
		spinner, _ = pterm.DefaultSpinner.Start("Checking migration status...")
		version, pending, err := database.GetMigrationStatus()
		if err != nil {
			spinner.Fail("Failed to get status: " + err.Error())
			return err
		}

		if !pending {
			spinner.Success(fmt.Sprintf("Database is up to date (Version: %d).", version))
			return nil
		}

		spinner.Warning(fmt.Sprintf("Pending migrations detected. Current version: %d.", version))

		// 2. Apply migrations
		confirm, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show("Do you want to apply pending migrations?")
		if !confirm {
			pterm.Info.Println("Operation cancelled.")
			return nil
		}

		spinner, _ = pterm.DefaultSpinner.Start("Applying migrations...")
		if err := database.RunMigrations(); err != nil {
			spinner.Fail("Migration failed: " + err.Error())
			return err
		}

		// Get new version
		newVersion, _, _ := database.GetMigrationStatus()
		spinner.Success(fmt.Sprintf("Migrations applied successfully! New version: %d.", newVersion))

		return nil
	},
}

func init() {
	dbCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return dbInteractive()
	}
	dbCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(dbCmd)
}
