package db

import (
	"fmt"
	"owner-api-proxy/internal/database/migration"
	"strings"

	"github.com/spf13/cobra"
)

func MigrateCommand() *cobra.Command {
	var down bool
	var configPath string
	var ssl string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			sslMode := strings.ToLower(strings.TrimSpace(ssl))
			if sslMode != "" && sslMode != "enable" && sslMode != "disable" {
				return fmt.Errorf("invalid --ssl value: %s (use enable or disable)", ssl)
			}

			if down {
				return migration.Down(&configPath, sslMode)
			}

			return migration.Up(&configPath, sslMode)
		},
	}

	cmd.Flags().BoolVar(&down, "down", false, "Run down migration")
	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")
	cmd.Flags().StringVar(&ssl, "ssl", "", "Override postgres SSL (enable|disable)")

	return cmd
}
