package db

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type RunPingFn func(options PingOptions) error

type PingOptions struct {
	ConfigPath string
	DBType     string
}

func PingCommand(runPing RunPingFn) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Test database connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			dbType, err := promptDBType()
			if err != nil {
				return err
			}

			return runPing(PingOptions{
				ConfigPath: configPath,
				DBType:     dbType,
			})
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}

func promptDBType() (string, error) {
	selected, err := SelectPrompt("Select database", []string{"postgres", "mysql"})
	if err != nil {
		return "", err
	}

	return selected, nil
}

func SelectPrompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  len(items),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}

	return result, nil
}
