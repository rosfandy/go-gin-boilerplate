package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type RunPullFn func(options PullOptions) error

type PullOptions struct {
	ConfigPath string
	Name       string
	DBType     string
}

func PullCommand(runPull RunPullFn) *cobra.Command {
	var name string
	var configPath string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			name = strings.TrimSpace(name)

			dbType, err := SelectPrompt("Select database", []string{"postgres", "mysql"})
			if err != nil {
				return err
			}

			fileName := strings.ToLower(nonAlphaNum.ReplaceAllString(name, "_"))
			fileName = strings.Trim(fileName, "_")
			if fileName == "" {
				return fmt.Errorf("invalid table name: %s", name)
			}

			modelPath := filepath.Join("internal", "database", "model", fileName+".go")
			if _, err := os.Stat(modelPath); err == nil {
				return fmt.Errorf("model file already exists: %s", modelPath)
			} else if !os.IsNotExist(err) {
				return err
			}

			return runPull(PullOptions{
				ConfigPath: configPath,
				Name:       name,
				DBType:     dbType,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Table name")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}

func PostgresTypeToGoType(dataType string, nullable bool) string {
	base := "string"

	switch strings.ToLower(dataType) {
	case "smallint", "integer", "bigint", "int2", "int4", "int8":
		base = "int64"
	case "real", "double precision", "numeric", "decimal", "float4", "float8":
		base = "float64"
	case "boolean", "bool":
		base = "bool"
	case "date", "timestamp without time zone", "timestamp with time zone", "time without time zone", "time with time zone", "timestamp", "timestamptz", "time", "timetz":
		base = "time.Time"
	case "bytea":
		base = "[]byte"
	}

	if nullable && base != "[]byte" {
		return "*" + base
	}

	return base
}

func postgresTypeToGoType(dataType string, nullable bool) string {
	return PostgresTypeToGoType(dataType, nullable)
}
