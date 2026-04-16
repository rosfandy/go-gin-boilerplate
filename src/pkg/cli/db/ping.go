package db

import (
	"fmt"

	"owner-api-proxy/internal/config"

	"github.com/manifoldco/promptui"
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

			dbType, err := promptDBType()
			if err != nil {
				return err
			}

			sqlClient := config.NewClientSql(dbType)
			if err := sqlClient.Open(); err != nil {
				return err
			}
			defer sqlClient.Close()

			db := sqlClient.Client()
			if db == nil {
				return fmt.Errorf("database client is not initialized")
			}

			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			if err := sqlDB.Ping(); err != nil {
				return err
			}

			fmt.Println("database connection is successful for", dbType)
			return nil
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
