package command

import (
	"owner-api-proxy/pkg/cli/server"

	"github.com/spf13/cobra"
)

func ServerCommand() *cobra.Command {
	return server.ServerCommand()
}
