package server

import (
	"github.com/spf13/cobra"
)

type RunServerFn func(configPath string) error

func ServerCommand(runServer RunServerFn) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(configPath)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}
