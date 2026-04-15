package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"owner-api-proxy/internal/config"

	"github.com/spf13/cobra"
)

func PullCommand() *cobra.Command {
	var name string
	var configPath string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull database schema",
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
			defer sqlDB.Close()

			columns, err := db.Migrator().ColumnTypes(name)
			if err != nil {
				return err
			}

			if len(columns) == 0 {
				return fmt.Errorf("table not found or has no columns: %s", name)
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

			structName := toPascalCase(name)
			var builder strings.Builder
			builder.WriteString("package model\n\n")
			builder.WriteString("type ")
			builder.WriteString(structName)
			builder.WriteString(" struct {\n")

			for _, col := range columns {
				columnName := col.Name()
				nullable, ok := col.Nullable()
				if !ok {
					nullable = true
				}

				fieldName := toPascalCase(columnName)
				fieldType := postgresTypeToGoType(col.DatabaseTypeName(), nullable)
				isPrimaryKey, _ := col.PrimaryKey()

				builder.WriteString("\t")
				builder.WriteString(fieldName)
				builder.WriteString(" ")
				builder.WriteString(fieldType)
				builder.WriteString(" `gorm:\"column:")
				builder.WriteString(columnName)
				if isPrimaryKey || strings.EqualFold(columnName, "id") {
					builder.WriteString(";primaryKey")
				}
				builder.WriteString("\" json:\"")
				builder.WriteString(columnName)
				builder.WriteString("\"`\n")
			}

			builder.WriteString("}\n\n")
			builder.WriteString("func (")
			builder.WriteString(structName)
			builder.WriteString(") TableName() string {\n")
			builder.WriteString("\treturn \"")
			builder.WriteString(name)
			builder.WriteString("\"\n}\n")

			if err := os.WriteFile(modelPath, []byte(builder.String()), 0644); err != nil {
				return err
			}

			fmt.Println("Model", structName, "is generated at", modelPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Pull target name")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}

func postgresTypeToGoType(dataType string, nullable bool) string {
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
