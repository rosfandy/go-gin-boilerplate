package command

import "github.com/spf13/cobra"

func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Application command",
	}

	cmd.AddCommand(DbCommands())
	cmd.AddCommand(ServerCommand())

	return cmd
}
