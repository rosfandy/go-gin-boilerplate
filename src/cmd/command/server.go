package command

import (
	"owner-api-proxy/internal/config"
	"owner-api-proxy/pkg/cli/server"

	"github.com/spf13/cobra"
)

func ServerCommand() *cobra.Command {
	cmd := server.ServerCommand()
	originalRunE := cmd.RunE

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := initPostgres(cmd); err != nil {
			return err
		}

		if originalRunE != nil {
			return originalRunE(cmd, args)
		}

		return nil
	}

	return cmd
}

func initPostgres(cmd *cobra.Command) error {
	log := config.NewLogrusWithCategory("db")

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if _, err := config.LoadConfig(&configPath); err != nil {
		return err
	}

	if _, err := config.InitPostgresDB(); err != nil {
		return err
	}

	log.Info("postgres connection initialized")

	return nil
}
