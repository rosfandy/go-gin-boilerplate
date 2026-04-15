package command

import (
	"owner-api-proxy/pkg/cli/db"

	"github.com/spf13/cobra"
)

func DbCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "db <command>",
		Short:   "Database command",
		Long:    "Command for database utilites",
		Example: "db pull",
	}

	cmd.AddCommand(
		db.PullCommand(),
		db.MigrateCommand(),
		db.ModelCommand(),
		db.PingCommand(),
	)

	return cmd
}
