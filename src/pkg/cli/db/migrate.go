package db

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type RunMigrateFn func(options MigrateOptions) error

type MigrateOptions struct {
	ConfigPath string
	SSLMode    string
	DBType     string
	Down       bool
}

func MigrateCommand(runMigrate RunMigrateFn) *cobra.Command {
	var down bool
	var configPath string
	var ssl string
	var conn string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			sslMode := strings.ToLower(strings.TrimSpace(ssl))
			if sslMode != "" && sslMode != "enable" && sslMode != "disable" {
				return fmt.Errorf("invalid --ssl value: %s (use enable or disable)", ssl)
			}

			actionDown := down
			if !cmd.Flags().Changed("down") {
				selectedDown, err := promptMigrationAction()
				if err != nil {
					return err
				}
				actionDown = selectedDown
			}

			dbType := strings.ToLower(strings.TrimSpace(conn))
			if dbType == "" {
				selectedDBType, err := promptMigrateDBType()
				if err != nil {
					return err
				}
				dbType = selectedDBType
			}
			if dbType != "postgres" && dbType != "mysql" {
				return fmt.Errorf("invalid --conn value: %s (use postgres or mysql)", conn)
			}

			return runMigrate(MigrateOptions{
				ConfigPath: configPath,
				SSLMode:    sslMode,
				DBType:     dbType,
				Down:       actionDown,
			})
		},
	}

	cmd.Flags().BoolVar(&down, "down", false, "Run down migration")
	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")
	cmd.Flags().StringVar(&ssl, "ssl", "", "Override postgres SSL (enable|disable)")
	cmd.Flags().StringVar(&conn, "conn", "", "Connection type (postgres|mysql) [required]")

	return cmd
}

func promptMigrationAction() (bool, error) {
	selected, err := SelectPrompt("Select migration action", []string{"up", "down"})
	if err != nil {
		return false, err
	}

	return selected == "down", nil
}

func promptMigrateDBType() (string, error) {
	return SelectPrompt("Select database", []string{"postgres", "mysql"})
}
