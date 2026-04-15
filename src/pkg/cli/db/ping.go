package db

import (
	"fmt"

	"owner-api-proxy/internal/config"

	"github.com/spf13/cobra"
)

func PingCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Test database connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := config.LoadConfig(&configPath); err != nil {
				return err
			}

			db, err := config.NewPostgresDB()
			if err != nil {
				return err
			}

			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			if err := sqlDB.Ping(); err != nil {
				return err
			}

			fmt.Println("database connection is successful")
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}
